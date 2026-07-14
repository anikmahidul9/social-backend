package store

import (
	"context"
	"database/sql"
	"encoding/json"
)

type Visibility string

const (
	Public  Visibility = "public"
	Private Visibility = "private"
)

type Post struct {
	ID          int64       `json:"id"`
	Content     string      `json:"content"`
	Title       string      `json:"title"`
	UserID      int64       `json:"user_ID"`
	Visibility  Visibility  `json:"visibility"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
	Comments    []Comment   `json:"comments"`
	User        User        `json:"user"`
	LikesCount  int         `json:"likes_count"`
	LatestLikes []User      `json:"latest_likes,omitempty"`
	Images      []PostImage `json:"images,omitempty"`
}

type PostWithMetadata struct {
	Post
	CommentCount int `json:"comments_count"`
}
type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `INSERT INTO posts (content, title, user_id, visibility)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, updated_at`

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		post.Visibility,
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostStore) Get(ctx context.Context, postID int64) (*Post, error) {
	var post Post
	post.User = User{}

	query := `
SELECT
    p.id,
    p.user_id,
    p.title,
    p.content,
    p.visibility,
    p.created_at,
    p.updated_at,

    u.id,
    u.first_name,
    u.last_name,
    u.username,
    u.email,
    u.created_at,
    u.updated_at,

    COALESCE(lc.likes_count, 0) AS likes_count,
    COALESCE(ll.latest_likes, '[]'::jsonb) AS latest_likes,
    COALESCE(img.images, '[]'::jsonb) AS images

FROM posts p

JOIN users u
    ON u.id = p.user_id

LEFT JOIN LATERAL (
    SELECT COUNT(*) AS likes_count
    FROM post_likes
    WHERE post_id = p.id
) lc ON TRUE

LEFT JOIN LATERAL (
    SELECT jsonb_agg(
        jsonb_build_object(
            'id', u2.id,
            'username', u2.username,
            'first_name', u2.first_name,
            'last_name', u2.last_name
        )
        ORDER BY pl.created_at DESC
    ) AS latest_likes
    FROM (
        SELECT *
        FROM post_likes
        WHERE post_id = p.id
        ORDER BY created_at DESC
        LIMIT 5
    ) pl
    JOIN users u2
        ON u2.id = pl.user_id
) ll ON TRUE

LEFT JOIN LATERAL (
    SELECT jsonb_agg(
        jsonb_build_object(
            'id', pi.id,
            'post_id', pi.post_id,
            'image_url', pi.image_url
        )
        ORDER BY pi.id
    ) AS images
    FROM post_images pi
    WHERE pi.post_id = p.id
) img ON TRUE

WHERE p.id = $1;
`

	var latestLikesJSON []byte
	var imagesJSON []byte

	err := s.db.QueryRowContext(ctx, query, postID).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.Visibility,
		&post.CreatedAt,
		&post.UpdatedAt,

		&post.User.ID,
		&post.User.FirstName,
		&post.User.LastName,
		&post.User.Username,
		&post.User.Email,
		&post.User.CreatedAt,
		&post.User.UpdatedAt,

		&post.LikesCount,
		&latestLikesJSON,
		&imagesJSON,
	)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(latestLikesJSON, &post.LatestLikes); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(imagesJSON, &post.Images); err != nil {
		return nil, err
	}

	return &post, nil
}
func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `UPDATE posts
SET
    title = $1,
    content = $2,
    visibility = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $4;`

	_, err := s.db.ExecContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.Visibility,
		post.ID,
	)
	if err != nil {
		return err
	}
	return nil

}

func (s *PostStore) Delete(ctx context.Context, postId int64) error {
	query := `DELETE FROM posts WHERE id = $1`

	res, err := s.db.ExecContext(ctx, query, postId)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *PostStore) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetadata, error) {
	query := `
SELECT
    p.id,
    p.user_id,
    p.title,
    p.content,
    p.visibility,
    p.created_at,

    u.username,

    COALESCE(cc.comment_count, 0) AS comment_count,
    COALESCE(lc.likes_count, 0) AS likes_count,

    COALESCE(ll.latest_likes, '[]'::jsonb) AS latest_likes,
    COALESCE(img.images, '[]'::jsonb) AS images

FROM posts p

JOIN users u
    ON u.id = p.user_id

LEFT JOIN LATERAL (
    SELECT COUNT(*) AS comment_count
    FROM comments
    WHERE post_id = p.id
) cc ON TRUE

LEFT JOIN LATERAL (
    SELECT COUNT(*) AS likes_count
    FROM post_likes
    WHERE post_id = p.id
) lc ON TRUE

LEFT JOIN LATERAL (
    SELECT jsonb_agg(
        jsonb_build_object(
            'id', u2.id,
            'username', u2.username,
            'first_name', u2.first_name,
            'last_name', u2.last_name
        )
        ORDER BY pl.created_at DESC
    ) AS latest_likes
    FROM (
        SELECT *
        FROM post_likes
        WHERE post_id = p.id
        ORDER BY created_at DESC
        LIMIT 5
    ) pl
    JOIN users u2
        ON u2.id = pl.user_id
) ll ON TRUE

LEFT JOIN LATERAL (
    SELECT jsonb_agg(
        jsonb_build_object(
            'id', pi.id,
            'post_id', pi.post_id,
            'image_url', pi.image_url
        )
        ORDER BY pi.id
    ) AS images
    FROM post_images pi
    WHERE pi.post_id = p.id
) img ON TRUE

WHERE
    p.user_id = $1
    OR p.visibility = 'public'

ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3;
`

	rows, err := s.db.QueryContext(
		ctx,
		query,
		userID,
		fq.Limit,
		fq.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var feed []PostWithMetadata
	for rows.Next() {
		var post PostWithMetadata
		var latestLikesJSON []byte
		var imagesJSON []byte

		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.Visibility,
			&post.CreatedAt,
			&post.User.Username,
			&post.CommentCount,
			&post.LikesCount,
			&latestLikesJSON,
			&imagesJSON,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(latestLikesJSON, &post.LatestLikes); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(imagesJSON, &post.Images); err != nil {
			return nil, err
		}
		feed = append(feed, post)
	}
	return feed, nil
}

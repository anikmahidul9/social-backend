package store

import (
	"context"
	"database/sql"
)

type Comment struct {
	ID              int64  `json:"id"`
	PostID          int64  `json:"post_id"`
	UserID          int64  `json:"user_id"`
	ParentCommentID *int64 `json:"parent_comment_id,omitempty"`
	Content         string `json:"content"`
	CreatedAt       string `json:"created_at"`

	User    User       `json:"user"`
	Replies []*Comment `json:"replies,omitempty"`
}

type CommentsStore struct {
	db *sql.DB
}

func (s *CommentsStore) GetByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	query := `SELECT
    c.id,
    c.post_id,
    c.user_id,
    c.parent_comment_id,
    c.content,
    c.created_at,
    u.id,
    u.username
FROM comments c
JOIN users u ON u.id = c.user_id
WHERE
    c.post_id = $1
    AND c.parent_comment_id IS NULL
ORDER BY c.created_at DESC;`

	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}

	for rows.Next() {
		var c Comment
		c.User = User{}

		err := rows.Scan(
			&c.ID,
			&c.PostID,
			&c.UserID,
			&c.ParentCommentID,
			&c.Content,
			&c.CreatedAt,
			&c.User.ID,
			&c.User.Username)
		if err != nil {
			return nil, err
		}

		replies, err := s.GetReplies(ctx, c.ID)
		if err != nil {
			return nil, err
		}

		c.Replies = replies
		comments = append(comments, c)
	}

	return comments, nil

}
func (s *CommentsStore) GetReplies(ctx context.Context, parentID int64) ([]*Comment, error) {

	query := `
	SELECT
		c.id,
		c.post_id,
		c.user_id,
		c.parent_comment_id,
		c.content,
		c.created_at,
		u.id,
		u.username
	FROM comments c
	JOIN users u ON u.id = c.user_id
	WHERE c.parent_comment_id = $1
	ORDER BY c.created_at ASC;
	`

	rows, err := s.db.QueryContext(ctx, query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []*Comment

	for rows.Next() {
		reply := &Comment{}

		err := rows.Scan(
			&reply.ID,
			&reply.PostID,
			&reply.UserID,
			&reply.ParentCommentID,
			&reply.Content,
			&reply.CreatedAt,
			&reply.User.ID,
			&reply.User.Username,
		)
		if err != nil {
			return nil, err
		}

		replies = append(replies, reply)
	}

	return replies, nil
}

func (s *CommentsStore) Create(ctx context.Context, comment *Comment) error {

	query := `
	INSERT INTO comments(
		post_id,
		user_id,
		parent_comment_id,
		content
	)
	VALUES($1,$2,$3,$4)
	RETURNING id,created_at;
	`

	return s.db.QueryRowContext(
		ctx,
		query,
		comment.PostID,
		comment.UserID,
		comment.ParentCommentID,
		comment.Content,
	).Scan(
		&comment.ID,
		&comment.CreatedAt,
	)
}

func (s *CommentsStore) CreateReplies(ctx context.Context, comment *Comment) error {

	query := `
	INSERT INTO comments(
		post_id,
		user_id,
		parent_comment_id,
		content
	)
	VALUES($1,$2,$3,$4)
	RETURNING id,created_at;
	`

	return s.db.QueryRowContext(
		ctx,
		query,
		comment.PostID,
		comment.UserID,
		comment.ParentCommentID,
		comment.Content,
	).Scan(
		&comment.ID,
		&comment.CreatedAt,
	)
}

func (s *CommentsStore) GetByID(ctx context.Context, id int64) (*Comment, error) {

	query := `
	SELECT
		id,
		post_id,
		user_id,
		parent_comment_id,
		content,
		created_at
	FROM comments
	WHERE id=$1;
	`

	comment := &Comment{}

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&comment.ID,
		&comment.PostID,
		&comment.UserID,
		&comment.ParentCommentID,
		&comment.Content,
		&comment.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return comment, nil
}

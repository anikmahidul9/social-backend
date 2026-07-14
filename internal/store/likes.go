package store

import (
	"context"
	"database/sql"
	"fmt"
)

const (
	PostLikeTable    = "post_likes"
	CommentLikeTable = "comment_likes"

	PostIDColumn    = "post_id"
	CommentIDColumn = "comment_id"
)

type LikeStore struct {
	db *sql.DB
}

func (s *LikeStore) Like(
	ctx context.Context,
	table string,
	column string,
	userID int64,
	targetID int64,
) error {

	query := fmt.Sprintf(`
		INSERT INTO %s(user_id,%s)
		VALUES($1,$2)
		ON CONFLICT DO NOTHING;
	`, table, column)

	_, err := s.db.ExecContext(ctx, query, userID, targetID)
	return err
}

func (s *LikeStore) Unlike(
	ctx context.Context,
	table string,
	column string,
	userID int64,
	targetID int64,
) error {

	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE user_id=$1
		AND %s=$2;
	`, table, column)

	_, err := s.db.ExecContext(ctx, query, userID, targetID)
	return err
}

func (s *LikeStore) GetLatestPostLikes(
	ctx context.Context,
	postID int64,
	limit int,
) ([]User, error) {

	query := `
	SELECT
		u.id,
		u.username,
		u.first_name,
		u.last_name
	FROM post_likes pl
	JOIN users u
		ON u.id = pl.user_id
	WHERE pl.post_id = $1
	ORDER BY pl.created_at DESC
	LIMIT $2;
	`

	rows, err := s.db.QueryContext(ctx, query, postID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}

	for rows.Next() {

		var user User

		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.FirstName,
			&user.LastName,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (s *LikeStore) CountPostLikes(
	ctx context.Context,
	postID int64,
) (int, error) {

	query := `
	SELECT COUNT(*)
	FROM post_likes
	WHERE post_id=$1;
	`

	var count int

	err := s.db.QueryRowContext(
		ctx,
		query,
		postID,
	).Scan(&count)

	return count, err
}
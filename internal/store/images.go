package store

import (
	"context"
	"database/sql"
)

type PostImage struct {
	ID        int64  `json:"id"`
	PostID    int64  `json:"post_id"`
	ImageURL  string `json:"image_url"`
	CreatedAt string `json:"created_at"`
}

type PostImageStore struct {
	db *sql.DB
}

func (s *PostImageStore) Create(
	ctx context.Context,
	postID int64,
	images []string,
) error {

	query := `
	INSERT INTO post_images(post_id,image_url)
	VALUES($1,$2);
	`

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, img := range images {

		_, err := stmt.ExecContext(
			ctx,
			postID,
			img,
		)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

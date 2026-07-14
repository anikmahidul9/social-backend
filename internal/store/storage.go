package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrDuplicateUsername = errors.New("duplicate username")
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		Get(context.Context, int64) (*Post, error)
		Update(context.Context, *Post) error
		Delete(context.Context, int64) error
		GetUserFeed(context.Context, int64, PaginatedFeedQuery) ([]PostWithMetadata, error)
	}

	Users interface {
		Create(context.Context, *User) error
		GetById(context.Context, int64) (*User, error)
		GetByEmail(ctx context.Context, email string) (*User, error)
	}

	Comments interface {
		GetByPostID(context.Context, int64) ([]Comment, error)
		Create(ctx context.Context, comment *Comment) error
		CreateReplies(ctx context.Context, comment *Comment) error
		GetByID(ctx context.Context, id int64) (*Comment, error)
		GetReplies(ctx context.Context, parentID int64) ([]*Comment, error)
	}

	Followers interface {
		Follow(ctx context.Context, followerId, userID int64) error
		UnFollow(ctx context.Context, followerId, userID int64) error
	}

	Reacts interface {
		Like(
			ctx context.Context,
			table string,
			column string,
			userID int64,
			targetID int64,
		) error
		Unlike(
			ctx context.Context,
			table string,
			column string,
			userID int64,
			targetID int64,
		) error
		GetLatestPostLikes(ctx context.Context, postID int64, limit int) ([]User, error)
		CountPostLikes(
			ctx context.Context,
			postID int64,
		) (int, error)
	}

	Images interface {
		Create(
			ctx context.Context,
			postID int64,
			images []string,
		) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStore{db},
		Users:     &UserStore{db},
		Comments:  &CommentsStore{db},
		Followers: &FollowerStore{db},
		Reacts:    &LikeStore{db},
		Images:    &PostImageStore{db},
	}
}

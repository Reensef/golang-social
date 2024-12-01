package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrDataNotFound = errors.New("data not found")
	QueryDBTimeout  = time.Second * 5
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (*Post, error)
		Update(context.Context, *Post) error
		DeleteByID(context.Context, int64) error
	}
	Users interface {
		Create(context.Context, *User) error
	}
	Comments interface {
		Create(context.Context, *Comment) error
		GetByID(context.Context, int64) (*Comment, error)
		GetByPostID(context.Context, int64) ([]Comment, error)
		Update(context.Context, *Comment) error
		DeleteByID(context.Context, int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostsStore{db},
		Users:    &UsersStore{db},
		Comments: &CommentsStore{db},
	}
}

package store

import (
	"context"
	"database/sql"
	"errors"
)

type Comment struct {
	ID        int64         `json:"id"`
	PostID    int64         `json:"post_id"`
	Sender    CommentSender `json:"sender"`
	Content   string        `json:"content"`
	CreatedAt string        `json:"created_at"`
	UpdatedAt string        `json:"updated_at"`
	Version   int64         `json:"-"`
}

type CommentSender struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type CommentsStore struct {
	db *sql.DB
}

func (s *CommentsStore) Create(ctx context.Context, comment *Comment) error {
	query := `
		INSERT INTO comments (content, post_id, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, (SELECT username FROM users WHERE id = $3), created_at, updated_at

	`
	ctx, cancel := context.WithTimeout(ctx, QueryDBTimeout)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		comment.Content,
		comment.PostID,
		comment.Sender.ID,
	).Scan(
		&comment.ID,
		&comment.Sender.Username,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *CommentsStore) GetByID(ctx context.Context, id int64) (*Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.user_id, u.username, c.content, c.created_at, c.updated_at, c.version
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryDBTimeout)
	defer cancel()

	comment := &Comment{}

	row := s.db.QueryRowContext(ctx, query, id)
	if row.Err() != nil {
		return nil, errors.Join(ErrDataNotFound, row.Err())
	}

	err := row.Scan(
		&comment.ID,
		&comment.PostID,
		&comment.Sender.ID,
		&comment.Sender.Username,
		&comment.Content,
		&comment.CreatedAt,
		&comment.UpdatedAt,
		&comment.Version,
	)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *CommentsStore) GetByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.content, users.id, users.username, c.created_at, c.updated_at
		FROM comments c
		JOIN users ON c.user_id = users.id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC
	`
	ctx, cancel := context.WithTimeout(ctx, QueryDBTimeout)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		comment := Comment{}
		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.Content,
			&comment.Sender.ID,
			&comment.Sender.Username,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (s *CommentsStore) Update(ctx context.Context, comment *Comment) error {
	query := `
		UPDATE comments
		SET content = $1, updated_at = now(), version = version + 1
		WHERE id = $2 AND version = $3
		RETURNING version
	`
	ctx, cancel := context.WithTimeout(ctx, QueryDBTimeout)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		comment.Content,
		comment.ID,
		comment.Version,
	).Scan(&comment.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrDataNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *CommentsStore) DeleteByID(ctx context.Context, id int64) error {
	query := `DELETE FROM comments WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryDBTimeout)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrDataNotFound
	}

	return nil
}

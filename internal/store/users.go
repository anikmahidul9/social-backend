package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64    `json:"id"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"_"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}
type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hash
	return nil
}

func (p *password) Matches(password string) error {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(password))
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users(first_name,last_name,username,email,password) VALUES ($1,$2,$3,$4,$5) RETURNING id,created_at`

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password.hash,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)

	if err != nil {
		var pqErr *pq.Error

		if errors.As(err, &pqErr) {
			switch pqErr.Constraint {
			case "users_username_key":
				return ErrDuplicateUsername

			case "users_email_key":
				return ErrDuplicateEmail
			}
		}
		return err

	}
	return nil
}

func (s *UserStore) GetById(ctx context.Context, userId int64) (*User, error) {

	query := `SELECT 
			id,username,email, created_at
			FROM users 
			WHERE id=$1`
	user := &User{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		userId,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil

}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {

	query := `
	SELECT
		id,
		first_name,
		last_name,
		username,
		email,
		password,
		created_at,
		updated_at
	FROM users
	WHERE email = $1
	`

	user := &User{}

	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

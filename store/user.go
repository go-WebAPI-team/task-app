package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-webapi-team/task-app/entity"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
)

func (r *Repository) CreateUser(ctx context.Context, db Execer, user *entity.User) error {
	const q = `INSERT INTO users(name, email, password, created_at, updated_at)
	VALUES (?,?,?,?)`
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	now := time.Now()
	_, err = db.ExecContext(ctx, q, user.Name, user.Email, user.Password, now, now)
	return err
}

func (r *Repository) Login(ctx context.Context, db Queryer, email, password string) (*entity.User, error) {
	const query = `SELECT id, name, email, password, created_at, updated_at
                   FROM users WHERE email = ?`
	var u entity.User

	err := db.QueryRowContext(ctx, query, email).Scan(
		&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// パスワードの検証
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, ErrInvalidPassword
	}

	return &u, nil
}

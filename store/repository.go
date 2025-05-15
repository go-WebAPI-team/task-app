package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // ドライバの side‑effect import
	"github.com/go-webapi-team/task-app/clock"
	"github.com/go-webapi-team/task-app/config"
)

const (
	pingTimeout = 2 * time.Second
)

// New は *sql.DB と cleanup 関数を返す。database/sql のみを使用。
func New(ctx context.Context, cfg *config.Config) (*sql.DB, func() error, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&loc=%s",
		cfg.DBUser, cfg.DBPass,
		cfg.DBHost, cfg.DBPort,
		cfg.DBName,
		cfg.DBLoc,
	)

	// ドライバを指定して Open（ここでは接続しない）
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	// Ping で実際に接続を確認
	ctxPing, cancel := context.WithTimeout(ctx, pingTimeout)

	defer cancel()
	if err := db.PingContext(ctxPing); err != nil {
		_ = db.Close()
		return nil, func() error { return nil }, err
	}

	return db, db.Close, nil
}

type Repository struct {
	Clocker clock.Clocker
}
type Queryer interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
type Execer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type ExecerQueryer interface {
	Queryer
	Execer
}

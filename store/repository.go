// store/repository.go
package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // ドライバの side‑effect import
	"github.com/your-module/config"    // 先ほど作成した package
)

const (
	defaultMaxOpenConns    = 10
	defaultMaxIdleConns    = 5
	defaultConnMaxLifetime = 2 * time.Hour
	pingTimeout            = 2 * time.Second
)

// New returns a ready‑to‑use *sql.DB, a cleanup function, and an error.
// ・database/sql だけを使用
// ・接続プールのチューニングもここで行う
func New(ctx context.Context, cfg *config.Config) (*sql.DB, func() error, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&loc=%s",
		cfg.DBUser, cfg.DBPass,
		cfg.DBHost, cfg.DBPort,
		cfg.DBName,	
		cfg.DBLoc,
	)

	// ❶ ドライバを指定して Open（ここでは接続しない）
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	// ❷ コネクションプール設定（任意で調整）
	db.SetMaxOpenConns(defaultMaxOpenConns)
	db.SetMaxIdleConns(defaultMaxIdleConns)
	db.SetConnMaxLifetime(defaultConnMaxLifetime)

	// ❸ Ping で実際に接続を確認
	ctxPing, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()
	if err := db.PingContext(ctxPing); err != nil {
		_ = db.Close()
		return nil, func() error { return nil }, err
	}

	// 呼び元で必ず Close してもらう
	cleanup := func() error { return db.Close() }
	return db, cleanup, nil
}

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/go-webapi-team/task-app/clock"
	"github.com/go-webapi-team/task-app/config"
	"github.com/go-webapi-team/task-app/store"
)

// run(ctx context.Context) errorに処理を委譲
func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, _, err := config.New()
	if err != nil {
		return err
	}

	// DB接続をここで作る
    db, cleanup, err := store.New(ctx, cfg)
    if err != nil { return err }
    defer cleanup()             // ← アプリ終了時に Close

    repo := &store.Repository{Clocker: clock.RealClocker{}}

	mux := NewMux(db, repo)


	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}
	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)

	s := NewServer(l, mux)
	
	return s.Run(ctx)
}
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/go-WebAPI-team/task-app/config"
)

// run(ctx context.Context) errorに処理を委譲
func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}
	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)

	mux := NewMux()
	s := NewServer(l, mux)
	
	return s.Run(ctx)
}
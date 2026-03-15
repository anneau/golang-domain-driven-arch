package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hkobori/golang-domain-driven-arch/internal/bootstrap"
)

func main() {
	app, err := bootstrap.NewWorkerApp()
	if err != nil {
		log.Fatalf("failed to bootstrap worker: %v", err)
	}
	defer app.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("worker starting...")
		if err := app.Start(ctx); err != nil {
			log.Printf("worker stopped: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down worker...")

	app.Stop()
	cancel()

	log.Println("worker exited")
}

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
	app, err := bootstrap.NewApp()
	if err != nil {
		log.Fatalf("failed to bootstrap app: %v", err)
	}
	defer app.Close()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("server starting on :%d", app.Port())
		if err := app.Start(); err != nil {
			log.Printf("server stopped: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down server...")

	if err := app.Shutdown(context.Background()); err != nil {
		log.Fatalf("failed to shutdown server: %v", err)
	}
	log.Println("server exited")
}

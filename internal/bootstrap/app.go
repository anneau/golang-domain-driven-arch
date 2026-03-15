package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	userusecase "github.com/hkobori/golang-domain-driven-arch/internal/app/usecase/user"
	adapterhttp "github.com/hkobori/golang-domain-driven-arch/internal/adapter/input/http"
	"github.com/hkobori/golang-domain-driven-arch/internal/adapter/input/http/handler"
	"github.com/hkobori/golang-domain-driven-arch/internal/adapter/output/database"
	outputqueue "github.com/hkobori/golang-domain-driven-arch/internal/adapter/output/queue"
	"github.com/hkobori/golang-domain-driven-arch/internal/domain/port"
	"github.com/hkobori/golang-domain-driven-arch/internal/domain/service"
)

type App struct {
	server *adapterhttp.Server
	db     *sql.DB
	port   int
}

func NewApp() (*App, error) {
	cfg := loadConfig()

	db, err := database.NewPostgres(database.DBConfig{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	userRepo := database.NewUserRepository(db)
	userDomainService := service.NewUserDomainService(userRepo)

	log.Printf("using SQS queue: %s (endpoint: %s)", cfg.SQS.QueueURL, cfg.SQS.EndpointURL)

	sqsPub, err := outputqueue.NewPublisher(context.Background(), cfg.SQS.QueueURL, cfg.SQS.EndpointURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create SQS publisher: %w", err)
	}

	var eventPublisher port.EventPublisher = sqsPub

	userUseCases := userusecase.NewUserUseCases(userRepo, userDomainService, eventPublisher)

	userHandler := handler.NewUserHandler(userUseCases.Create)

	srv := adapterhttp.NewServer(
		adapterhttp.ServerConfig{
			Port:         cfg.Server.Port,
			AllowOrigins: cfg.Server.AllowOrigins,
		},
		userHandler,
	)

	return &App{
		server: srv,
		db:     db,
		port:   cfg.Server.Port,
	}, nil
}

func (a *App) Start() error {
	return a.server.Start()
}

func (a *App) Port() int {
	return a.port
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}

func (a *App) Close() {
	if a.db != nil {
		_ = a.db.Close()
	}
}

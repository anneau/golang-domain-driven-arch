package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/hkobori/golang-domain-driven-arch/internal/adapter/input/http/handler"
	"github.com/hkobori/golang-domain-driven-arch/internal/adapter/input/http/middleware"
)

type ServerConfig struct {
	Port         int
	AllowOrigins []string
}

type Server struct {
	echo   *echo.Echo
	config ServerConfig
}

func NewServer(cfg ServerConfig, userHandler *handler.UserHandler) *Server {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true
	e.HTTPErrorHandler = middleware.ErrorHandler
	e.Validator = middleware.NewValidator()

	e.Use(middleware.RequestID())
	e.Use(middleware.RequestLogger())
	e.Use(echoMiddleware.Recover())
	e.Use(middleware.CORS(cfg.AllowOrigins))

	registerRoutes(e, userHandler)

	return &Server{echo: e, config: cfg}
}

func registerRoutes(e *echo.Echo, userHandler *handler.UserHandler) {
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	v1 := e.Group("/api/v1")
	users := v1.Group("/users")
	users.POST("", userHandler.Create)
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.Port)
	s.echo.Logger.Infof("starting server on %s", addr)
	if err := s.echo.Start(addr); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return s.echo.Shutdown(shutdownCtx)
}

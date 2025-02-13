package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/22Fariz22/wallet/config"
	"github.com/22Fariz22/wallet/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

// Server struct
type Server struct {
	echo        *echo.Echo
	cfg         *config.Config
	db          *sqlx.DB
	redisClient *redis.Client
	logger      logger.Logger
}

// CustomValidator wraps validator
type CustomValidator struct {
	Validator *validator.Validate
}

// Validate method for Echo
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

// NewServer New Server constructor
func NewServer(cfg *config.Config, db *sqlx.DB, redisClient *redis.Client, logger logger.Logger) *Server {
	e := echo.New()

	// Устанавливаем кастомный валидатор
	e.Validator = &CustomValidator{Validator: validator.New()}

	return &Server{echo: e, cfg: cfg, db: db, redisClient: redisClient, logger: logger}
}

func (s *Server) Run() error {
	// serverAddr := fmt.Sprintf("%s:%s", s.cfg.Server.BaseUrl, s.cfg.Server.Port)

	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", s.cfg.Server.Port),
		ReadTimeout:    time.Second * s.cfg.Server.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.Server.WriteTimeout,
		MaxHeaderBytes: s.cfg.Server.MaxHeaderBytes,
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				s.logger.Errorf("Recovered from panic: %v", r)
			}
		}()
		s.logger.Infof("Server is listening on PORT: %s", s.cfg.Server.Port)
		if err := s.echo.StartServer(server); err != nil {
			s.logger.Fatalf("Error starting Server: %v", err)
		}
	}()

	if err := s.MapHandlers(s.echo); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), s.cfg.Server.CtxTimeout)
	defer shutdown()

	s.logger.Info("Server Exited Properly")
	return s.echo.Server.Shutdown(ctx)
}

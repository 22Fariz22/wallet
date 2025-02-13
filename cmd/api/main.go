package main

import (
	"fmt"
	"log"

	"github.com/22Fariz22/wallet/config"
	"github.com/22Fariz22/wallet/internal/server"
	"github.com/22Fariz22/wallet/pkg/db/migrate"
	"github.com/22Fariz22/wallet/pkg/db/postgres"
	"github.com/22Fariz22/wallet/pkg/db/redis"
	"github.com/22Fariz22/wallet/pkg/logger"
)

func main() {
	log.Println("Starting api server")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	appLogger := logger.NewApiLogger(cfg)

	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

	appLogger.Debugf("Postgres config: host=%s port=%s user=%s dbname=%s sslmode=%s password=%s PgDriver=%s",
		cfg.Postgres.PostgresqlHost,
		cfg.Postgres.PostgresqlPort,
		cfg.Postgres.PostgresqlUser,
		cfg.Postgres.PostgresqlDbname,
		cfg.Postgres.PostgresqlSSLMode,
		cfg.Postgres.PostgresqlPassword,
		cfg.Postgres.PgDriver,
	)

	// Формирование строки подключения для GORM
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		cfg.Postgres.PostgresqlHost,
		cfg.Postgres.PostgresqlPort,
		cfg.Postgres.PostgresqlUser,
		cfg.Postgres.PostgresqlDbname,
		cfg.Postgres.PostgresqlPassword,
	)

	// Выполнение миграций
	if err := migrate.Migrate(appLogger, dsn); err != nil {
		appLogger.Errorf("Failed to run migrations: %v", err)
	}
	appLogger.Debug("Database migrated successfully")

	psqlDB, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	} else {
		appLogger.Infof("Postgres connected, Status: %#v", psqlDB.Stats())
	}
	defer psqlDB.Close()

	redisClient := redis.NewRedisClient(cfg)
	defer redisClient.Close()
	appLogger.Info("Redis connected")

	s := server.NewServer(cfg, psqlDB, redisClient, appLogger)
	if err = s.Run(); err != nil {
		appLogger.Fatalf("Error in main NewServer(): ", err)
	}
}

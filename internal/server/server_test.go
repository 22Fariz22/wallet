package server

import (
	"testing"

	"github.com/22Fariz22/wallet/config"
	"github.com/22Fariz22/wallet/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: "8080",
		},
	}
	db := &sqlx.DB{}
	redisClient := &redis.Client{}
	log := logger.NewApiLogger(cfg)

	server := NewServer(cfg, db, redisClient, log)

	assert.NotNil(t, server)
	assert.NotNil(t, server.echo)
	assert.Equal(t, cfg, server.cfg)
}

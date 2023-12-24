package servers

import (
	"context"
	"github.com/ndodanli/backend-api/configs"
	"github.com/ndodanli/backend-api/pkg/logger"
	"time"
)

const (
	certFile       = "ssl/server.crt"
	keyFile        = "ssl/server.pem"
	maxHeaderBytes = 1 << 20
	gzipLevel      = 5
	stackSize      = 4 << 10 // 4 KB
	bodyLimit      = "3M"
	timeout        = 120 * time.Second
)

// server
type server struct {
	cfg    *configs.Config
	ctx    *context.Context
	logger logger.ILogger
}

// NewServer constructor
func NewServer(cfg *configs.Config, ctx *context.Context, logger logger.ILogger) *server {
	return &server{cfg: cfg, ctx: ctx, logger: logger}
}

func (s *server) GetLoggerInstance() logger.ILogger {
	return s.logger
}

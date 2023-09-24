package servers

import (
	"context"
	"github.com/ndodanli/go-clean-architecture/configs"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
	"time"
)

const (
	certFile       = "ssl/server.crt"
	keyFile        = "ssl/server.pem"
	maxHeaderBytes = 1 << 20
	gzipLevel      = 5
	stackSize      = 4 << 10 // 4 KB
	bodyLimit      = "3M"
	timeout        = 3000 * time.Second
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

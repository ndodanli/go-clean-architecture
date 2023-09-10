package servers

import (
	"context"
	"github.com/ndodanli/go-clean-architecture/configs"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
)

const (
	certFile        = "ssl/server.crt"
	keyFile         = "ssl/server.pem"
	maxHeaderBytes  = 1 << 20
	gzipLevel       = 5
	stackSize       = 4 << 10 // 4 KB
	csrfTokenHeader = "X-CSRF-Token"
	bodyLimit       = "3M"
)

// server
type server struct {
	cfg    *configs.Config
	ctx    *context.Context
	logger logger.Logger
}

// NewServer constructor
func NewServer(cfg *configs.Config, ctx *context.Context, logger logger.Logger) *server {
	return &server{cfg: cfg, ctx: ctx, logger: logger}
}

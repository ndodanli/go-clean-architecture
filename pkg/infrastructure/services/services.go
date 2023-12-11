package services

import (
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/repo"
)

type AppServices struct {
	JWTService      IJWTService
	RedisService    IRedisService
	AppUserRepo     repo.IAppUserRepo
	SendgridService ISendgridService
}

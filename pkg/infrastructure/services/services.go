package services

import (
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/repo"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services/redissrv"
)

type AppServices struct {
	JWTService   IJWTService
	RedisService redissrv.IRedisService
	AppUserRepo  repo.IAppUserRepo
}

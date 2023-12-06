package queries

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/ndodanli/go-clean-architecture/pkg/constant"
	baseres "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql"
	uow "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql/unit_of_work"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services/redissrv"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
)

type TestQueryHandler struct {
	uow          uow.IUnitOfWork
	jwtService   services.IJWTService
	logger       logger.ILogger
	redisService redissrv.IRedisService
}

func NewTestQueryHandler(appServices *services.AppServices, uow uow.IUnitOfWork, logger logger.ILogger) *TestQueryHandler {
	return &TestQueryHandler{
		uow:          uow,
		jwtService:   appServices.JWTService,
		redisService: appServices.RedisService,
		logger:       logger,
	}
}

type TestQuery struct {
	TestID string
}

type TestQueryResponse struct {
	TestIDRes string `json:"testIDRes"`
}

func (h *TestQueryHandler) Handle(echoCtx echo.Context, query *TestQuery) *baseres.Result[TestQueryResponse, error, struct{}] {
	result := baseres.NewResult[TestQueryResponse, error, struct{}]()
	ctx := echoCtx.Request().Context()
	tm := echoCtx.Get(constant.General.TxSessionManagerKey).(*postgresql.TxSessionManager)
	appUserRepo := h.uow.AppUserRepo(ctx)

	redisErr := h.redisService.Ping(ctx)
	if redisErr != nil {
		h.logger.Error(redisErr.Error(), redisErr, "app")
	}

	redisResult, _ := redissrv.AcquireString(ctx, h.redisService.Client(), "test1", 30, func() (string, error) {
		return "test", nil
	})
	_ = redisResult

	redisHashResult, _ := redissrv.AcquireHash(ctx, h.redisService.Client(), "testHash", 600, []string{}, func() (struct {
		Test1 string
		Test2 string
		Test3 int64
	}, error) {
		return struct {
			Test1 string
			Test2 string
			Test3 int64
		}{
			Test1: "test1",
			Test2: "test2",
			Test3: 423423,
		}, nil
	})
	_ = redisHashResult
	fmt.Print(redisHashResult.Test2)

	updateProps := map[string]interface{}{
		"username": "testfdsfd",
	}

	_, err := appUserRepo.PatchAppUser(1, updateProps, tm)
	if err != nil {
		return result.Err(err)
	}

	return result.Ok()
}

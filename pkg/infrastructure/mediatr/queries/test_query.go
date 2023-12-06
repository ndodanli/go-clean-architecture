package queries

import (
	"github.com/labstack/echo/v4"
	"github.com/ndodanli/go-clean-architecture/pkg/constant"
	baseres "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql"
	uow "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql/unit_of_work"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services/redissrv"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
	"reflect"
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
	TestArray []int  `json:"testArray"`
}

func (h *TestQueryHandler) Handle(echoCtx echo.Context, query *TestQuery) *baseres.Result[TestQueryResponse, error, struct{}] {
	result := baseres.NewResult[TestQueryResponse, error, struct{}]()
	result.Data.TestIDRes = "Test 123"
	result.Data.TestArray = []int{}
	ctx := echoCtx.Request().Context()
	tm := echoCtx.Get(constant.General.TxSessionManagerKey).(*postgresql.TxSessionManager)
	appUserRepo := h.uow.AppUserRepo(ctx)

	redisErr := h.redisService.Ping(ctx)
	if redisErr != nil {
		h.logger.Error(redisErr.Error(), redisErr, "app")
	}

	redisSetStringResult := redissrv.SetString(ctx, h.redisService.Client(), "testKeySet", "testValueSet", 0)
	_ = redisSetStringResult

	//redisSetHashFieldResult := redissrv.SetHashField(ctx, h.redisService.Client(), "testMasterKey", "testHashField", result, 0)
	//_ = redisSetHashFieldResult

	redisSetHashResult := redissrv.SetHash(ctx, h.redisService.Client(), "testMasterKey", result, 0)
	_ = redisSetHashResult

	redisAcquireHashResult, err := redissrv.AcquireHash(ctx, h.redisService.Client(), "testMasterKey1", 0, []string{}, func() (*baseres.Result[TestQueryResponse, error, struct{}], error) {
		return result, nil
	})

	check := reflect.DeepEqual(result, redisAcquireHashResult)
	_ = check
	_ = redisAcquireHashResult
	updateProps := map[string]interface{}{
		"username": "testfdsfd",
	}

	_, err = appUserRepo.PatchAppUser(1, updateProps, tm)
	if err != nil {
		return result.Err(err)
	}

	return result.Ok()
}

package adminqueries

import (
	"github.com/labstack/echo/v4"
	baseres "github.com/ndodanli/backend-api/pkg/core/response"
	"github.com/ndodanli/backend-api/pkg/infrastructure/db/sqldb/pg"
	uow "github.com/ndodanli/backend-api/pkg/infrastructure/db/sqldb/pg/unit_of_work"
	"github.com/ndodanli/backend-api/pkg/infrastructure/services"
	"github.com/ndodanli/backend-api/pkg/logger"
	"reflect"
)

type TestQueryHandler struct {
	UOW         uow.IUnitOfWork
	Logger      logger.ILogger
	AppServices *services.AppServices
	TM          *pg.TxSessionManager
}

type TestQuery struct {
	TestID string
}

type TestQueryResponse struct {
	TestIDRes string `json:"testIDRes"`
	TestArray []int  `json:"testArray"`
}

func (h *TestQueryHandler) Handle(echoCtx echo.Context, query *TestQuery) *baseres.Result[*TestQueryResponse, error, struct{}] {
	result := baseres.NewResult[*TestQueryResponse, error, struct{}](nil)
	result.Data.TestIDRes = "Test 123"
	result.Data.TestArray = []int{}
	ctx := echoCtx.Request().Context()
	appUserRepo := h.UOW.AppUserRepo(ctx, h.TM)

	redisErr := h.AppServices.RedisService.Ping(ctx)
	if redisErr != nil {
		h.Logger.Error(redisErr.Error(), redisErr, "app")
	}

	redisSetStringResult := services.SetString(ctx, h.AppServices.RedisService.Client(), "testKeySet", "testValueSet", 0)
	_ = redisSetStringResult

	//redisSetHashFieldResult := redissrv.SetHashField(ctx, h.AppServices.RedisService.Client(), "testMasterKey", "testHashField", result, 0)
	//_ = redisSetHashFieldResult

	redisSetHashResult := services.SetHash(ctx, h.AppServices.RedisService.Client(), "testMasterKey", result, 0)
	_ = redisSetHashResult

	redisAcquireHashResult, err := services.AcquireHash(ctx, h.AppServices.RedisService.Client(), "testMasterKey1", 0, []string{}, func() (*baseres.Result[*TestQueryResponse, error, struct{}], error) {
		return result, nil
	})

	check := reflect.DeepEqual(result, redisAcquireHashResult)
	_ = check
	_ = redisAcquireHashResult
	updateProps := map[string]interface{}{
		"username": "testfdsfd",
	}

	_, err = appUserRepo.PatchUser(1, updateProps)
	if err != nil {
		return result.Err(err)
	}

	return result.Ok()
}

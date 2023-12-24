package testctrl

import (
	"github.com/labstack/echo/v4"
	_ "github.com/ndodanli/backend-api/pkg/core/response"
	baseres "github.com/ndodanli/backend-api/pkg/core/response"
	"github.com/ndodanli/backend-api/pkg/infrastructure/mediatr"
	"github.com/ndodanli/backend-api/pkg/infrastructure/mediatr/queries"
	mw "github.com/ndodanli/backend-api/pkg/infrastructure/middleware"
	"github.com/ndodanli/backend-api/pkg/logger"
	"net/http"
	"os"
)

type TestController struct {
	cGroup     *echo.Group
	httpClient *http.Client
	logger     logger.ILogger
}

func NewTestController(group *echo.Group, logger logger.ILogger) (*TestController, error) {
	err := RegisterMediatrHandlers()
	if err != nil {
		return nil, err
	}
	if os.Getenv("APP_ENV") == "test" {
		return nil, err
	}
	ac := &TestController{
		cGroup: group.Group("/test"),
		logger: logger,
	}

	ac.cGroup.GET("/test", ac.Test, mw.Auth)

	return ac, nil
}

func (ac *TestController) Test(c echo.Context) error {
	res := mediatr.Send[*queries.TestQuery, *baseres.Result[queries.TestQueryResponse, error, struct{}]](c, &queries.TestQuery{
		TestID: "test",
	})
	if res.IsErr() {
		return res.GetErr()
	}
	return c.JSON(http.StatusOK, res)
}

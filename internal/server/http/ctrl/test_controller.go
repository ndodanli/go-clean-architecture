package ctrl

import (
	"github.com/labstack/echo/v4"
	_ "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/mediatr"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/mediatr/queries"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
	"net/http"
)

type TestController struct {
	controllerGroup *echo.Group
	httpClient      *http.Client
	logger          logger.ILogger
}

func NewTestController(group *echo.Group, logger logger.ILogger) *TestController {
	ac := &TestController{
		controllerGroup: group.Group("/test"),
		logger:          logger,
	}

	ac.controllerGroup.GET("/test", ac.Test)

	return ac
}

func (ac *TestController) Test(c echo.Context) error {
	res := mediatr.Send[*queries.TestQuery, *queries.TestQueryResponse](c, &queries.TestQuery{
		TestID: "test",
	})

	return c.JSON(http.StatusOK, res)
}

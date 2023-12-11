package servers

import (
	"github.com/labstack/echo/v4"
	authctrl "github.com/ndodanli/go-clean-architecture/internal/server/http/ctrl/auth"
	testctrl "github.com/ndodanli/go-clean-architecture/internal/server/http/ctrl/test"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
)

type AppManager struct {
	*authctrl.AuthController
	*testctrl.TestController
	*echo.Echo
}

func RegisterControllers(e *echo.Group, logger logger.ILogger) error {
	_, err := authctrl.NewAuthController(e, logger)
	if err != nil {
		return err
	}
	_, err = testctrl.NewTestController(e, logger)
	if err != nil {
		return err
	}

	return nil
}

package servers

import (
	"github.com/labstack/echo/v4"
	adminauthctrl "github.com/ndodanli/backend-api/internal/server/http/admin/ctrl/admin_auth"
	adminuserctrl "github.com/ndodanli/backend-api/internal/server/http/admin/ctrl/admin_user"
	authctrl "github.com/ndodanli/backend-api/internal/server/http/ctrl/auth"
	testctrl "github.com/ndodanli/backend-api/internal/server/http/ctrl/test"
	"github.com/ndodanli/backend-api/pkg/logger"
)

type AppManager struct {
	*authctrl.AuthController
	*testctrl.TestController
	*echo.Echo
}

func RegisterControllers(e *echo.Group, logger logger.ILogger) error {
	var err error

	adminGroup := e.Group("/admin")
	// Admins
	_, err = adminauthctrl.NewAdminAuthController(adminGroup, logger)
	if err != nil {
		return err
	}
	_, err = adminuserctrl.NewAdminUserController(adminGroup, logger)
	if err != nil {
		return err
	}

	userGroup := e.Group("")
	// AppUsers
	_, err = authctrl.NewAuthController(userGroup, logger)
	if err != nil {
		return err
	}
	_, err = testctrl.NewTestController(userGroup, logger)
	if err != nil {
		return err
	}

	return nil
}

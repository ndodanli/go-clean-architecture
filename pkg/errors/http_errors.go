package httperr

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type ErrorData struct {
	Message   string
	ShouldLog bool
}

var (
	UserNotFoundError *echo.HTTPError
	UnauthorizedError *echo.HTTPError
)

func Init() {
	UnauthorizedError = echo.NewHTTPError(http.StatusUnauthorized, &ErrorData{
		Message:   "Unauthorized",
		ShouldLog: false,
	})
	UserNotFoundError = echo.NewHTTPError(http.StatusNotFound, &ErrorData{
		Message:   "User not found",
		ShouldLog: false,
	})
}

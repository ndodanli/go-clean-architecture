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
	InternalServerError *echo.HTTPError
)

var (
	UserNotFoundError           *echo.HTTPError
	UnauthorizedError           *echo.HTTPError
	UnAuthorizedAudienceError   *echo.HTTPError
	UsernameOrPasswordIncorrect *echo.HTTPError
)

func Init() {
	InternalServerError = echo.NewHTTPError(http.StatusInternalServerError, &ErrorData{
		Message:   "Internal server error",
		ShouldLog: true,
	})
	UnauthorizedError = echo.NewHTTPError(http.StatusUnauthorized, &ErrorData{
		Message:   "Unauthorized",
		ShouldLog: false,
	})
	UserNotFoundError = echo.NewHTTPError(http.StatusNotFound, &ErrorData{
		Message:   "User not found",
		ShouldLog: false,
	})
	UnAuthorizedAudienceError = echo.NewHTTPError(http.StatusUnauthorized, &ErrorData{
		Message:   "Unauthorized audience",
		ShouldLog: false,
	})
	UsernameOrPasswordIncorrect = echo.NewHTTPError(http.StatusUnauthorized, &ErrorData{
		Message:   "Username or password is incorrect",
		ShouldLog: false,
	})
}

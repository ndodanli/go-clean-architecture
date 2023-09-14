package httperr

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type ErrorData struct {
	Message   string
	ShouldLog bool
}

// Dynamic errors
var (
	BindingError = func(message string) *echo.HTTPError {
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorData{
			Message:   "Binding error: " + message,
			ShouldLog: false,
		})
	}
)

// Static errors
var (
	InternalServerError      *echo.HTTPError
	InvalidRefreshTokenError *echo.HTTPError
)

var (
	UserNotFoundError                *echo.HTTPError
	UnauthorizedError                *echo.HTTPError
	UnAuthorizedAudienceError        *echo.HTTPError
	UsernameOrPasswordIncorrectError *echo.HTTPError
	RefreshTokenNotFoundError        *echo.HTTPError
	RefreshTokenExpiredError         *echo.HTTPError
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
	UsernameOrPasswordIncorrectError = echo.NewHTTPError(http.StatusUnauthorized, &ErrorData{
		Message:   "Username or password is incorrect",
		ShouldLog: false,
	})
	RefreshTokenNotFoundError = echo.NewHTTPError(http.StatusUnauthorized, &ErrorData{
		Message:   "Refresh token not found",
		ShouldLog: true,
	})
	RefreshTokenExpiredError = echo.NewHTTPError(http.StatusUnauthorized, &ErrorData{
		Message:   "Refresh token expired",
		ShouldLog: false,
	})
	InvalidRefreshTokenError = echo.NewHTTPError(http.StatusBadRequest, &ErrorData{
		Message:   "Invalid refresh token",
		ShouldLog: true,
	})
}

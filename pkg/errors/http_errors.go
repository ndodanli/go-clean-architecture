package httperr

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type ErrorData struct {
	Message          string
	Metadata         interface{}
	ShouldLogAsError bool
	ShouldLogAsInfo  bool
}

// Dynamic errors
var (
	BindingError = func(message string) *echo.HTTPError {
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorData{
			Message:          "Binding error: " + message,
			ShouldLogAsError: false,
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
	InvalidAuthenticationError       *echo.HTTPError
	UsernameOrPasswordIncorrectError *echo.HTTPError
	RefreshTokenNotFoundError        *echo.HTTPError
	RefreshTokenExpiredError         *echo.HTTPError
)

func Init() {
	InternalServerError = echo.NewHTTPError(http.StatusInternalServerError, &ErrorData{
		Message:          "Internal server error",
		ShouldLogAsError: true,
	})
	UnauthorizedError = echo.NewHTTPError(http.StatusUnauthorized, &ErrorData{
		Message: "Unauthorized",
	})
	UserNotFoundError = echo.NewHTTPError(http.StatusNotFound, &ErrorData{
		Message: "User not found",
	})
	UnAuthorizedAudienceError = echo.NewHTTPError(http.StatusUnauthorized, &ErrorData{
		Message: "Unauthorized audience",
	})
	UsernameOrPasswordIncorrectError = echo.NewHTTPError(http.StatusUnauthorized, &ErrorData{
		Message: "Username or password is incorrect",
	})
	RefreshTokenNotFoundError = echo.NewHTTPError(http.StatusUnauthorized, &ErrorData{
		Message:         "Refresh token not found",
		ShouldLogAsInfo: true,
	})
	RefreshTokenExpiredError = echo.NewHTTPError(http.StatusUnauthorized, &ErrorData{
		Message: "Refresh token expired",
	})
	InvalidRefreshTokenError = echo.NewHTTPError(http.StatusBadRequest, &ErrorData{
		Message:         "Invalid refresh token",
		ShouldLogAsInfo: true,
	})
	InvalidAuthenticationError = echo.NewHTTPError(http.StatusUnauthorized, &ErrorData{
		Message: "Invalid authentication",
	})
}

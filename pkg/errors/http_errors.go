package httperr

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type ErrorData struct {
	Status           int
	Message          string
	Metadata         interface{}
	ShouldLogAsError bool
	ShouldLogAsInfo  bool
}

// Dynamic errors
var (
	BindingError = func(message string) *echo.HTTPError {
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorData{
			Status:           500,
			Message:          "Binding error: " + message,
			ShouldLogAsError: false,
		})
	}
	SendgridError = func(message string) *echo.HTTPError {
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorData{
			Status:           500,
			Message:          "Sendgrid error: " + message,
			ShouldLogAsError: true,
		})
	}
)

// Static errors
var (
	InternalServerError      *echo.HTTPError
	InvalidRefreshTokenError *echo.HTTPError
)

var (
	AppUserNotFoundError                       *echo.HTTPError
	UnauthorizedError                          *echo.HTTPError
	UnAuthorizedAudienceError                  *echo.HTTPError
	InvalidAuthenticationError                 *echo.HTTPError
	UsernameOrPasswordIncorrectError           *echo.HTTPError
	RefreshTokenNotFoundError                  *echo.HTTPError
	RefreshTokenExpiredError                   *echo.HTTPError
	AppUserAlreadyConfirmedError               *echo.HTTPError
	CodeRecentlySentError                      *echo.HTTPError
	ConfirmationCodeExpiredError               *echo.HTTPError
	InvalidConfirmationCodeError               *echo.HTTPError
	PasswordsDoNotMatchError                   *echo.HTTPError
	CannotChangePasswordEmailNotConfirmedError *echo.HTTPError
	PasswordCannotBeSameAsOldError             *echo.HTTPError
)

func Init() {
	InternalServerError = echo.NewHTTPError(http.StatusInternalServerError, &ErrorData{
		Status:           500,
		Message:          "Internal server error",
		ShouldLogAsError: true,
	})
	UnauthorizedError = echo.NewHTTPError(http.StatusUnauthorized, &ErrorData{
		Status:  401,
		Message: "Unauthorized",
	})
	AppUserNotFoundError = echo.NewHTTPError(http.StatusNotFound, &ErrorData{
		Status:  404,
		Message: "User not found",
	})
	UnAuthorizedAudienceError = echo.NewHTTPError(http.StatusUnauthorized, &ErrorData{
		Status:  401,
		Message: "Unauthorized audience",
	})
	UsernameOrPasswordIncorrectError = echo.NewHTTPError(http.StatusUnauthorized, &ErrorData{
		Status:  401,
		Message: "Username or password is incorrect",
	})
	RefreshTokenNotFoundError = echo.NewHTTPError(http.StatusUnauthorized, &ErrorData{
		Status:          401,
		Message:         "Refresh token not found",
		ShouldLogAsInfo: true,
	})
	RefreshTokenExpiredError = echo.NewHTTPError(http.StatusUnauthorized, &ErrorData{
		Status:  401,
		Message: "Refresh token expired",
	})
	InvalidRefreshTokenError = echo.NewHTTPError(http.StatusBadRequest, &ErrorData{
		Status:          400,
		Message:         "Invalid refresh token",
		ShouldLogAsInfo: true,
	})
	InvalidAuthenticationError = echo.NewHTTPError(http.StatusUnauthorized, &ErrorData{
		Status:  401,
		Message: "Invalid authentication",
	})
	AppUserAlreadyConfirmedError = echo.NewHTTPError(http.StatusBadRequest, &ErrorData{
		Status:          400,
		Message:         "User already confirmed",
		ShouldLogAsInfo: true,
	})
	CodeRecentlySentError = echo.NewHTTPError(http.StatusBadRequest, &ErrorData{
		Status:  400,
		Message: "Code recently sent",
	})
	ConfirmationCodeExpiredError = echo.NewHTTPError(http.StatusBadRequest, &ErrorData{
		Status:  400,
		Message: "Confirmation code expired",
	})
	InvalidConfirmationCodeError = echo.NewHTTPError(http.StatusBadRequest, &ErrorData{
		Status:  400,
		Message: "Invalid confirmation code",
	})
	PasswordsDoNotMatchError = echo.NewHTTPError(http.StatusBadRequest, &ErrorData{
		Status:  400,
		Message: "Passwords do not match",
	})
	CannotChangePasswordEmailNotConfirmedError = echo.NewHTTPError(http.StatusBadRequest, &ErrorData{
		Status:  400,
		Message: "Cannot change password because email not confirmed",
	})
	PasswordCannotBeSameAsOldError = echo.NewHTTPError(http.StatusBadRequest, &ErrorData{
		Status:  400,
		Message: "Password cannot be same as old",
	})
}

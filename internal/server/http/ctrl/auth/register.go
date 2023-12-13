package authctrl

import (
	baseres "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/mediatr"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/mediatr/queries"
)

func RegisterMediatrHandlers() error {
	var err error
	err = mediatr.RegisterRequestHandler[
		*queries.LoginQuery, *baseres.Result[*queries.LoginQueryResponse, error, struct{}],
	](&queries.LoginQueryHandler{})
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[
		*queries.RefreshTokenQuery, *baseres.Result[*queries.RefreshTokenQueryResponse, error, struct{}],
	](&queries.RefreshTokenQueryHandler{})
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[
		*queries.SendConfirmationEmailForgotPasswordQuery, *baseres.Result[*queries.SendConfirmationEmailForgotPasswordQueryResponse, error, struct{}],
	](&queries.SendConfirmationEmailForgotPasswordQueryHandler{})
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[
		*queries.ConfirmForgotPasswordCodeQuery, *baseres.Result[*queries.ConfirmForgotPasswordCodeQueryResponse, error, struct{}],
	](&queries.ConfirmForgotPasswordCodeQueryHandler{})
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[
		*queries.LoginWithGoogleQuery, *baseres.Result[*queries.LoginWithGoogleQueryResponse, error, struct{}],
	](&queries.LoginWithGoogleQueryHandler{})
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[
		*queries.EmailConfirmationQuery, *baseres.Result[*queries.EmailConfirmationQueryResponse, error, struct{}],
	](&queries.EmailConfirmationQueryHandler{})
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[
		*queries.RegisterQuery, *baseres.Result[*queries.RegisterQueryResponse, error, struct{}],
	](&queries.RegisterQueryHandler{})
	if err != nil {
		return err
	}

	return nil
}

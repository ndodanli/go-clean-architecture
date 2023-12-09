package authctrl

import (
	baseres "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/mediatr"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/mediatr/queries"
)

func RegisterMediatrHandlers() error {
	var err error
	err = mediatr.RegisterRequestHandler[
		*queries.LoginQuery, *baseres.Result[queries.LoginQueryResponse, error, struct{}],
	](&queries.LoginQueryHandler{})
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[
		*queries.RefreshTokenQuery, *baseres.Result[queries.RefreshTokenQueryResponse, error, struct{}],
	](&queries.RefreshTokenQueryHandler{})
	if err != nil {
		return err
	}

	return nil
}

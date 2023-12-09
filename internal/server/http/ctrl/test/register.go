package testctrl

import (
	baseres "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/mediatr"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/mediatr/queries"
)

func RegisterMediatrHandlers() error {
	var err error
	err = mediatr.RegisterRequestHandler[
		*queries.TestQuery, *baseres.Result[queries.TestQueryResponse, error, struct{}],
	](&queries.TestQueryHandler{})
	if err != nil {
		return err
	}

	return nil
}

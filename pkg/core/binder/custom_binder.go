package cstmbinder

import (
	"errors"
	"github.com/labstack/echo/v4"
	httperr "github.com/ndodanli/backend-api/pkg/errors"
)

type CustomBinder struct{}

func NewCustomBinder() *CustomBinder {
	return &CustomBinder{}
}

func (cb *CustomBinder) Bind(i interface{}, c echo.Context) error {
	db := new(echo.DefaultBinder)
	if err := db.Bind(i, c); err != echo.ErrUnsupportedMediaType {
		var he *echo.HTTPError
		if errors.As(err, &he) {
			return httperr.BindingError(he.Message.(string))
		}
	}

	return c.Bind(i)
}

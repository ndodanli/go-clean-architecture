package utils

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type ReqParams struct {
	Params  interface{}
	HttpReq *http.Request
}

func BindAndValidate(c echo.Context, params interface{}) error {
	if err := c.Bind(params); err != nil {
		return err
	}

	reqParams := &ReqParams{
		Params:  params,
		HttpReq: c.Request(),
	}

	if err := c.Validate(reqParams); err != nil {
		return err
	}

	return nil
}

package utils

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type ReqParams struct {
	Params  interface{}
	HttpReq *http.Request
}

func BindAndValidate(c echo.Context, params interface{}) error {
	handleQueryParams(c.QueryParams())

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

func handleQueryParams(queryParams map[string][]string) {
	for key, value := range queryParams {
		if len(value) > 0 {
			if value[0][0] == '[' {
				var arr []string
				if value[0] != "[]" {
					arr = strings.Split(value[0][1:len(value[0])-1], ",")
					queryParams[key] = arr
				} else {
					delete(queryParams, key)
				}
			}
		}
	}
}

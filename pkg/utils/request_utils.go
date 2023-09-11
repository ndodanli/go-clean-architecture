package utils

import "github.com/labstack/echo/v4"

func BindAndValidate(c echo.Context, reqParams interface{}) error {
	if err := c.Bind(reqParams); err != nil {
		return err
	}
	if err := c.Validate(reqParams); err != nil {
		return err
	}

	return nil
}

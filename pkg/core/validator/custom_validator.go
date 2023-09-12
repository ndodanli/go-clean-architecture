package cstmvalidator

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	res "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	"github.com/ndodanli/go-clean-architecture/pkg/utils"
	"net/http"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewCustomValidator(validator *validator.Validate) *CustomValidator {
	return &CustomValidator{
		validator: validator,
	}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			result := res.NewResult[any, error, string]()

			out := make([]res.ValidationError, len(ve))
			for j, fe := range ve {
				out[j] = res.ValidationError{Field: utils.LowercaseFirstLetter(fe.Field()), Error: parseValidationMessages(fe)}
			}
			result.SetValidationErrors(out)
			return echo.NewHTTPError(http.StatusBadRequest, result)
		}
	}
	return nil
}

func parseValidationMessages(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "min":
		return "Minimum length is " + fe.Param()
	case "max":
		return "Maximum length is " + fe.Param()
	}
	return fe.Error()
}

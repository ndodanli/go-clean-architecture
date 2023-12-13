package cstmvalidator

import (
	"errors"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/tr"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	res "github.com/ndodanli/go-clean-architecture/pkg/core/response"
	apperr "github.com/ndodanli/go-clean-architecture/pkg/errors/app_errors"
	"github.com/ndodanli/go-clean-architecture/pkg/utils"
	"net/http"
	"regexp"
)

type CustomValidator struct {
	validator   *validator.Validate
	translators map[string]ut.Translator
}

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
)

func NewCustomValidator(validator *validator.Validate) *CustomValidator {
	uni = ut.New(tr.New(), tr.New(), en.New())

	trTranslator, found := uni.GetTranslator("tr")
	if !found {
		panic("translator not found" + " tr")
	}

	enTranslator, found := uni.GetTranslator("en")
	if !found {
		panic("translator not found" + " en")
	}

	cv := &CustomValidator{
		validator: validator,
		translators: map[string]ut.Translator{
			"tr": trTranslator,
			"en": enTranslator,
		},
	}

	cv.RegisterTranslations()

	cv.RegisterCustomValidations()
	return cv
}

func (cv *CustomValidator) RegisterTranslations() {
	var err error
	err = cv.validator.RegisterTranslation("required", cv.GetDefaultTranslator(), func(ut ut.Translator) error {
		return ut.Add("required", "Bu alan boş bırakılamaz", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})
	if err != nil {
		panic(err)
	}
}

func (cv *CustomValidator) GetDefaultTranslator() ut.Translator {
	return cv.translators["tr"]
}

// #region Custom validations
func strongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	// Check for at least one uppercase letter and one digit using regular expressions
	match, _ := regexp.MatchString(`(?=.*[A-Z])(?=.*\d)`, password)
	return match
}

func (cv *CustomValidator) RegisterCustomValidations() {
	var err error
	err = cv.validator.RegisterTranslation("strongPassword", cv.GetDefaultTranslator(), func(ut ut.Translator) error {
		return ut.Add("strongPassword", "Parolanız en az 8 karakter uzunluğunda olmalı ve en az bir büyük harf ve bir rakam içermelidir", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("strongPassword", fe.Field())
		return t
	})
	if err != nil {
		panic(err)
	}
	err = cv.validator.RegisterValidation("strongPassword", strongPassword)
	if err != nil {
		panic(err)
	}
}

//#endregion

func (cv *CustomValidator) Validate(i interface{}) error {
	if reqParams, ok := i.(*utils.ReqParams); ok {
		var translator ut.Translator
		language := reqParams.HttpReq.Header.Get("X-Lang")
		if language == "" {
			language = reqParams.HttpReq.Header.Get("Accept-Language")
			if language == "" {
				translator = cv.GetDefaultTranslator()
			} else {
				var found bool
				translator, found = uni.GetTranslator(language)
				if !found {
					translator = cv.GetDefaultTranslator()
				}
			}
		}

		if err := cv.validator.Struct(reqParams.Params); err != nil {
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				result := res.NewResult[any, error, string](nil)

				out := make([]res.ValidationError, len(ve))
				for j, fe := range ve {
					out[j] = res.ValidationError{Field: utils.LowercaseFirstLetter(fe.Field()), Error: fe.Translate(translator)}
				}
				result.SetValidationErrors(out)
				return echo.NewHTTPError(http.StatusBadRequest, result)
			}
		}
	} else {
		return apperr.RequestParamsHasToBeReqParamsError
	}
	return nil
}

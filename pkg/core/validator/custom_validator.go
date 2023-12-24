package cstmvalidator

import (
	"errors"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/tr"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	res "github.com/ndodanli/backend-api/pkg/core/response"
	apperr "github.com/ndodanli/backend-api/pkg/errors/app_errors"
	"github.com/ndodanli/backend-api/pkg/utils"
	"net/http"
	"reflect"
	"regexp"
	"time"
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
func (cv *CustomValidator) strongPassword() error {
	var err error
	err = cv.validator.RegisterTranslation("strongPassword", cv.GetDefaultTranslator(), func(ut ut.Translator) error {
		return ut.Add("strongPassword", "Parolanız en az 8 karakter uzunluğunda olmalı ve en az bir büyük harf ve bir rakam içermelidir", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("strongPassword", fe.Field())
		return t
	})
	if err != nil {
		return err
	}
	err = cv.validator.RegisterValidation("strongPassword", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()

		if len(password) < 8 {
			return false
		}

		// Check for at least one uppercase letter and one digit using regular expressions
		match, _ := regexp.MatchString(`(?=.*[A-Z])(?=.*\d)`, password)
		return match
	})
	if err != nil {
		return err
	}

	return nil
}

func (cv *CustomValidator) gtNow() error {
	err := cv.validator.RegisterTranslation("gtNow", cv.GetDefaultTranslator(), func(ut ut.Translator) error {
		return ut.Add("gtNow", "{0} alanı şu anki tarihten daha ileri bir tarih olmalıdır", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("gtNow", fe.Field())
		return t
	})
	if err != nil {
		return err
	}
	err = cv.validator.RegisterValidation("gtNow", func(fl validator.FieldLevel) bool {
		fieldTime := fl.Field().Interface().(time.Time)
		if err != nil {
			return false
		}
		return fieldTime.After(time.Now())
	})
	if err != nil {
		return err
	}
	return nil
}

func (cv *CustomValidator) arrayInt64(typ reflect.Type) error {
	err := cv.validator.RegisterTranslation("arrayInt64", cv.GetDefaultTranslator(), func(ut ut.Translator) error {
		return ut.Add("arrayInt64", "{0} alanı int64 array olmalıdır", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("arrayInt64", fe.Field())
		return t
	})
	if err != nil {
		return err
	}
	err = cv.validator.RegisterValidation("arrayInt64", func(fl validator.FieldLevel) bool {
		return checkArray(fl, typ)
	})
	if err != nil {
		return err
	}
	return nil
}

func checkArray(fl validator.FieldLevel, typ reflect.Type) bool {
	field := fl.Field()

	// Check if the field is an array or slice
	if field.Kind() != reflect.Slice && field.Kind() != reflect.Array {
		return false
	}

	// Check if all elements are of type int64
	for i := 0; i < field.Len(); i++ {
		if field.Index(i).Kind() != typ.Elem().Kind() {
			return false
		}
	}

	return true
}

func (cv *CustomValidator) RegisterCustomValidations() {

	var err error
	err = cv.strongPassword()
	if err != nil {
		panic(err)
	}

	err = cv.gtNow()
	if err != nil {
		panic(err)
	}

	err = cv.arrayInt64(reflect.TypeOf([]int64{}))
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

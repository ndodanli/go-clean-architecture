package baseres

import (
	"errors"
	"github.com/labstack/echo/v4"
	httperr "github.com/ndodanli/go-clean-architecture/pkg/errors"
)

type resultType any
type ErrorType error
type metadataType any

type ValidationError struct {
	Field string `json:"f,omitempty" example:"age"`
	Error string `json:"e,omitempty" example:"age must be greater than 0"`
}

type Result[D resultType, E ErrorType, M metadataType] struct {
	Success          bool              `json:"s"`
	Message          string            `json:"m,omitempty"`
	Data             D                 `json:"d,omitempty"`
	ValidationErrors []ValidationError `json:"v,omitempty"`
	TestArray        []int             `json:"testArray,omitempty"`
	TestMap          map[string]int    `json:"testMap,omitempty"`
	error            ErrorType
	metadata         metadataType
}

func (r *Result[D, E, M]) IsSuccess() bool {
	return r.Success
}

func (r *Result[D, E, M]) IsErr() bool {
	return !r.Success
}

func (r *Result[D, E, M]) GetErr() ErrorType {
	var he *echo.HTTPError
	if ok := errors.As(r.error, &he); ok {
		_, ok = he.Message.(*httperr.ErrorData)
		if ok {
			he.Message.(*httperr.ErrorData).Metadata = r.metadata
		}
	}
	return r.error
}

func (r *Result[D, E, M]) GetErrorMessage() string {
	if r.error != nil {
		var he *echo.HTTPError
		if ok := errors.As(r.error, &he); ok {
			_, ok = he.Message.(*httperr.ErrorData)
			if ok {
				return he.Message.(*httperr.ErrorData).Message
			} else {
				if _, ok = he.Message.(string); ok {
					return he.Message.(string)
				}
			}
		}
		return r.error.Error()
	} else {
		return ""
	}
}

func (r *Result[D, E, M]) GetMessage() string {
	return r.Message
}

func (r *Result[D, E, M]) GetMetadata() any {
	return r.metadata
}

func (r *Result[D, E, M]) SetMessage(message string) {
	r.Message = message
}

func (r *Result[D, E, M]) SetErrorMessage(message string) {
	r.Success = false
	r.Message = message
}

func (r *Result[D, E, M]) SetMetadata(metadata interface{}) {
	r.metadata = metadata
}

func (r *Result[D, E, M]) SetValidationErrors(validationErrors []ValidationError) {
	r.Success = false
	r.ValidationErrors = validationErrors
}

func (r *Result[D, E, M]) Ok() *Result[D, E, M] {
	r.Success = true
	return r
}

func (r *Result[D, E, M]) OkWithData(data D) *Result[D, E, M] {
	r.Success = true
	r.Data = data
	return r
}

func (r *Result[D, E, M]) OkWithMetadata(metadata metadataType) *Result[D, E, M] {
	r.Success = true
	r.metadata = metadata
	return r
}

func (r *Result[D, E, M]) Err(error E) *Result[D, E, M] {
	r.Success = false
	r.error = error
	if r.Message == "" {
		r.Message = error.Error()
	}
	return r
}

func (r *Result[D, E, M]) ErrWithMetadata(error E, metadata metadataType) *Result[D, E, M] {
	r.Success = false
	r.error = error
	r.metadata = metadata
	return r
}

func NewResult[R resultType, E ErrorType, M metadataType]() *Result[R, E, M] {
	return &Result[R, E, M]{
		Success: true,
	}
}

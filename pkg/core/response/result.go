package res

type resultType any
type ErrorType error
type metadataType any

type ValidationError struct {
	Field string `json:"F,omitempty" example:"Age"`
	Error string `json:"E,omitempty" example:"Age must be greater than 0"`
}

type Result[D resultType, E ErrorType, M metadataType] struct {
	Success          bool              `json:"S"`
	Message          string            `json:"M,omitempty"`
	Data             D                 `json:"D,omitempty"`
	ValidationErrors []ValidationError `json:"V,omitempty"`
	error            ErrorType
	metadata         metadataType
}

func (r *Result[D, E, M]) IsSuccess() bool {
	return r.Success
}

func (r *Result[D, E, M]) IsError() bool {
	return !r.Success
}

func (r *Result[D, E, M]) GetError() ErrorType {
	return r.error
}

func (r *Result[D, E, M]) GetErrorMessage() string {
	return r.error.Error()
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

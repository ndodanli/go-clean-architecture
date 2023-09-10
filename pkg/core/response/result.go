package res

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

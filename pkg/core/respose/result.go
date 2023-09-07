package res

type resultType any
type errorType error
type metadataType any

type Result[D resultType, E errorType, M metadataType] struct {
	success  bool
	message  string
	Data     D `json:"data"`
	error    errorType
	metadata metadataType
}

func (r *Result[D, E, M]) IsSuccess() bool {
	return r.success
}

func (r *Result[D, E, M]) IsError() bool {
	return !r.success
}

func (r *Result[D, E, M]) GetError() errorType {
	return r.error
}

func (r *Result[D, E, M]) GetMessage() string {
	return r.message
}

func (r *Result[D, E, M]) GetMetadata() any {
	return r.metadata
}

func (r *Result[D, E, M]) SetMessage(message string) {
	r.message = message
}

func (r *Result[D, E, M]) SetErrorMessage(message string) {
	r.success = false
	r.message = message
}

func (r *Result[D, E, M]) SetMetadata(metadata interface{}) {
	r.metadata = metadata
}

func (r *Result[D, E, M]) Ok() *Result[D, E, M] {
	r.success = true
	return r
}

func (r *Result[D, E, M]) OkWithMetadata(metadata metadataType) *Result[D, E, M] {
	r.success = true
	r.metadata = metadata
	return r
}

func (r *Result[D, E, M]) Err(error E) *Result[D, E, M] {
	r.success = false
	r.error = error
	return r
}

func (r *Result[D, E, M]) ErrWithMetadata(error E, metadata metadataType) *Result[D, E, M] {
	r.success = false
	r.error = error
	r.metadata = metadata
	return r
}

func NewResult[R resultType, E errorType, M metadataType]() *Result[R, E, M] {
	return &Result[R, E, M]{
		success: true,
	}
}

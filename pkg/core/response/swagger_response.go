package res

type SwaggerSuccessRes[T any] struct {
	Success bool   `json:"s" example:"true" nullable:"false"`
	Data    T      `json:"d" nullable:"true"`
	Message string `json:"m,omitempty" example:"XXX Created/Updated/Deleted Successfully" nullable:"true"`
}

type SwaggerValidationErrRes struct {
	Success bool              `json:"s" example:"false" nullable:"false"`
	Errors  []ValidationError `json:"v" nullable:"false"`
}

type SwaggerInternalErrRes struct {
	Success bool   `json:"s" example:"false" nullable:"false"`
	Message string `json:"m,omitempty" example:"Internal Server Error" nullable:"true"`
}

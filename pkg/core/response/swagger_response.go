package res

type SwaggerSuccessRes[T any] struct {
	Success bool   `json:"S" example:"true" nullable:"false"`
	Data    T      `json:"D" nullable:"true"`
	Message string `json:"M,omitempty" example:"XXX Created/Updated/Deleted Successfully" nullable:"true"`
}

type SwaggerValidationErrRes struct {
	Success bool              `json:"S" example:"false" nullable:"false"`
	Errors  []ValidationError `json:"V" nullable:"false"`
}

type SwaggerInternalErrRes struct {
	Success bool   `json:"S" example:"false" nullable:"false"`
	Message string `json:"M,omitempty" example:"Internal Server Error" nullable:"true"`
}

type SwaggerUnauthorizedErrRes struct {
	Success bool   `json:"S" example:"false" nullable:"false"`
	Message string `json:"M,omitempty" example:"Unauthorized" nullable:"true"`
}

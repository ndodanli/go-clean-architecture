package apperr

import "errors"

var (
	ResultMustBeStruct                 = errors.New("Result must be struct")
	ValueIsSettableOrAddressable       = errors.New("Value is not settable or addressable")
	ReturnFuncValueNil                 = errors.New("Return func value is nil")
	RequestParamsHasToBeReqParamsError = errors.New("Request params has to be req params")
	FieldNotFoundWithColumnName        = errors.New("Field not found with column name")
)

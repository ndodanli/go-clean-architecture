package apperr

import "errors"

var (
	ResultMustBeStruct           = errors.New("result must be struct")
	ValueIsSettableOrAddressable = errors.New("Value is not settable or addressable")
	ReturnFuncValueNil           = errors.New("return func value is nil")
)

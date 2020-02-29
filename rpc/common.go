package rpc

import (
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/stretchr/objx"
	"gitlab.com/makeos/mosdef/types"
	"gitlab.com/makeos/mosdef/util"
)

func isStrFieldSet(val *objx.Value) bool {
	return val.Str() != ""
}

// GetStringFromObjxMap gets a field from the given objx.Map ensuring its type is a 'string'.
//
// o: The objx map.
// field: The target field to get.
// required: Whether to return error if field does not exist in the object map.
func GetStringFromObjxMap(o objx.Map, field string, required bool) (string, *Response) {

	vField := o.Get(field)
	if !isStrFieldSet(vField) && !required {
		return "", nil
	}

	if !isStrFieldSet(vField) {
		return "", Error(types.RPCErrCodeInvalidParamValue, field+" is required", nil)
	}

	if !vField.IsStr() {
		msg := util.WrongFieldValueMsg(field, "string", vField.Inter())
		return "", Error(types.RPCErrCodeInvalidParamValue, msg, nil)
	}

	return vField.Str(), nil
}

// GetStringToUint64FromObjxMap gets a field from the given objx.Map ensuring
// its type is a 'string' that is convertible to Uint64
//
// o: The objx map.
// field: The target field to get.
// required: Whether to return error if field does not exist in the object map.
func GetStringToUint64FromObjxMap(o objx.Map, field string, required bool) (uint64, *Response) {

	vField := o.Get(field)
	if !isStrFieldSet(vField) && !required {
		return 0, nil
	}

	if !isStrFieldSet(vField) {
		return 0, Error(types.RPCErrCodeInvalidParamValue, field+" is required", nil)
	}

	if !vField.IsStr() {
		msg := util.WrongFieldValueMsg(field, "string", vField.Inter())
		return 0, Error(types.RPCErrCodeInvalidParamValue, msg, nil)
	}

	if !govalidator.IsNumeric(vField.Str()) {
		msg := field + " requires a numeric value"
		return 0, Error(types.RPCErrCodeInvalidParamValue, msg, nil)
	}

	fieldValue, err := strconv.ParseUint(vField.Str(), 10, 64)
	if err != nil {
		return 0, Error(types.RPCErrCodeInvalidParamValue, "failed to parse to unsigned integer", nil)
	}

	return fieldValue, nil
}

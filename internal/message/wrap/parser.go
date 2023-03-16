package wrap

import (
	"github.com/whaoinfo/go-box/logger"
	"github.com/whaoinfo/go-box/mapping"
	"reflect"
)

func CallFieldUnmarshalBinary(ownRV reflect.Value, index int, args ...interface{}) error {
	if ownRV.Kind() == reflect.Ptr {
		ownRV = ownRV.Elem()
	}

	var fieldType string
	field := ownRV.Field(index)
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			field = mapping.NewFieldByRV(field)
			ownRV.Field(index).Set(field)
		}
		fieldType = field.Elem().Type().Name()
	} else {
		fieldType = field.Type().Name()
	}

	//logger.AllFmt("Call UnmarshalBinary function of %s type, field index: %d", fieldType, index)
	callArgs := args[0].([]reflect.Value)
	if err := mapping.CallFieldMethodByName(field, "UnmarshalBinary", callArgs...); err != nil {
		logger.WarnFmt("The UnmarshalBinary function of %s type has failed to call, field index: %d, err: %v",
			fieldType, index, err)
		return err
	}

	//logger.AllFmt("The UnmarshalBinary function of %s type has called, field index: %d", fieldType, index)
	return nil
}

func CallFieldMarshalBinary(ownRV reflect.Value, index int, args ...interface{}) error {
	callArgs := args[0]
	markFilter := args[1].(*bool)

	if ownRV.Kind() == reflect.Ptr {
		ownRV = ownRV.Elem()
	}
	field := ownRV.Field(index)
	var fieldType string
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			*markFilter = true
			return nil
		}
		fieldType = field.Elem().Type().Name()
	} else {
		fieldType = field.Type().Name()
	}

	*markFilter = false
	logger.AllFmt("Call MarshalBinary function, index: %d, type: %v", index, fieldType)
	if err := mapping.CallFieldMethodByName(field, "MarshalBinary", callArgs.([]reflect.Value)...); err != nil {
		logger.WarnFmt("the MarshalBinary function has failed to call, index: %d, type: %v, err: %v",
			index, fieldType, err)
		return err
	}

	logger.AllFmt("The MarshalBinary function has called, index: %d, type: %v", index, fieldType)
	return nil
}

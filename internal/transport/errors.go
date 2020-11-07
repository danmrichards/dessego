package transport

import (
	"reflect"
)

// InvalidDecodeTargetError is returned when a request is attempted to be
// decoded into an invalid type (type must be a non-nil pointer).
type InvalidDecodeTargetError struct {
	Type reflect.Type
}

func (i InvalidDecodeTargetError) Error() string {
	return "target must be a non-nil pointer, got:" + i.Type.String()
}

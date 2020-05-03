package gofig

import (
	"fmt"
	"reflect"
)

// ErrInvalidValue is returned when gogif is given an invalid type.
type ErrInvalidValue struct {
	Type reflect.Type
}

func (e ErrInvalidValue) Error() string {
	return fmt.Sprintf("destination must be a pointer to a struct got: %s", e.Type.String())
}

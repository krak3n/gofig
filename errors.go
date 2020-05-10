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

// ErrInvalidConversion is returned whe a type cannot be cast to a another.
type ErrInvalidConversion struct {
	From reflect.Kind
	To   reflect.Kind
}

func (e ErrInvalidConversion) Error() string {
	return fmt.Sprintf("invalid type conversion: %s > %s", e.From, e.To)
}

// ErrSetValue is returned when a field cannot be set to the given value.
type ErrSetValue struct {
	Field reflect.Value
	Value reflect.Value
}

func (e ErrSetValue) Error() string {
	return fmt.Sprintf("could not set value of %v(%s) on field %v(%s)",
		e.Value.Interface(),
		e.Value.Kind(),
		e.Field.Interface(),
		e.Field.Kind(),
	)
}

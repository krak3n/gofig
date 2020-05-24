package gofig

import (
	"fmt"
	"reflect"
	"strings"
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

// CloseError is returned by Close when one or more notifiers error on their Close.
type CloseError struct {
	errors []error
}

func (e *CloseError) Error() string {
	if e == nil || e.errors == nil {
		return "no errors"
	}

	points := make([]string, len(e.errors))
	for i, err := range e.errors {
		points[i] = fmt.Sprintf("* %s", err)
	}

	return fmt.Sprintf(
		"%d errors occurred:\n\t%s\n\n",
		len(e.errors), strings.Join(points, "\n\t"))
}

// Add adds one or more errors to the error group.
func (e *CloseError) Add(errs ...error) {
	if e == nil {
		*e = *new(CloseError)
	}

	if e.errors == nil {
		e.errors = make([]error, len(errs))
	}

	e.errors = append(e.errors, errs...)
}

// NilOrError returns an error if there are errors, else nil.
func (e *CloseError) NilOrError() error {
	if e == nil {
		return nil
	}

	if e.errors == nil || len(e.errors) == 0 {
		return nil
	}

	return e
}

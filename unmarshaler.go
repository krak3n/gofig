package gofig

import "reflect"

// Unmarshaler is an interface implemented by types that can unmarshal a values themselves.
type Unmarshaler interface {
	UnmarshalGoFig(value interface{}) error
}

// unmarshaler checks to see if the given field implements the Unmarshaler interface.
// If it does the Unmarshaler is returned, else nil is returned.
func unmarshaler(field reflect.Value) Unmarshaler {
	if field.Kind() != reflect.Ptr && field.Type().Name() != "" && field.CanAddr() {
		field = field.Addr()
	}

	for {
		if field.Kind() == reflect.Interface && !field.IsNil() {
			e := field.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() && e.Elem().Kind() == reflect.Ptr {
				field = e
				continue
			}
		}

		if field.Kind() != reflect.Ptr {
			break
		}

		if field.CanSet() {
			break
		}

		if field.Elem().Kind() == reflect.Interface && field.Elem().Elem() == field {
			field = field.Elem()
			break
		}

		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}

		if field.Type().NumMethod() > 0 && field.CanInterface() {
			if u, ok := field.Interface().(Unmarshaler); ok {
				return u
			}

			break
		}

		field = field.Elem()
	}

	return nil
}

package gofig

import "reflect"

// Unmarshaler is an interface implemented by types that can unmarshal a values themselves.
type Unmarshaler interface {
	UnmarshalGoFig(value interface{}) error
}

// unmarshaler checks to see if the given field implements the Unmarshaler interface.
// If it does the Unmarshaler is returned, else nil is returned.
// Lifted from go json stdlib
func unmarshaler(value reflect.Value) Unmarshaler {
	if value.Kind() != reflect.Ptr && value.Type().Name() != "" && value.CanAddr() {
		value = value.Addr()
	}

	for {
		if value.Kind() == reflect.Interface && !value.IsNil() {
			e := value.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() && e.Elem().Kind() == reflect.Ptr {
				value = e
				continue
			}
		}

		if value.Kind() != reflect.Ptr {
			break
		}

		if value.CanSet() {
			break
		}

		if value.Elem().Kind() == reflect.Interface && value.Elem().Elem() == value {
			value = value.Elem()
			break
		}

		if value.IsNil() {
			value.Set(reflect.New(value.Type().Elem()))
		}

		if value.Type().NumMethod() > 0 && value.CanInterface() {
			if u, ok := value.Interface().(Unmarshaler); ok {
				return u
			}

			break
		}

		value = value.Elem()
	}

	return nil
}

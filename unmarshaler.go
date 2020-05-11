package gofig

import "reflect"

// Unmarshaler is an interface implemented by types that can unmarshal a values themselves.
type Unmarshaler interface {
	UnmarshalGoFig(value interface{}) error
}

// unmarshaler checks to see if the given field implements the Unmarshaler interface.
// If it does the Unmarshaler is returned, else nil is returned.
// Lifted from go json stdlib
func unmarshaler(field Field) Unmarshaler {
	fv := field.Value

	if fv.Kind() != reflect.Ptr && fv.Type().Name() != "" && fv.CanAddr() {
		fv = fv.Addr()
	}

	for {
		if fv.Kind() == reflect.Interface && !fv.IsNil() {
			e := fv.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() && e.Elem().Kind() == reflect.Ptr {
				fv = e
				continue
			}
		}

		if fv.Kind() != reflect.Ptr {
			break
		}

		if fv.CanSet() {
			break
		}

		if fv.Elem().Kind() == reflect.Interface && fv.Elem().Elem() == fv {
			fv = fv.Elem()
			break
		}

		if fv.IsNil() {
			fv.Set(reflect.New(fv.Type().Elem()))
		}

		if fv.Type().NumMethod() > 0 && fv.CanInterface() {
			if u, ok := fv.Interface().(Unmarshaler); ok {
				return u
			}

			break
		}

		fv = fv.Elem()
	}

	return nil
}

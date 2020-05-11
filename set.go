package gofig

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// setValue sets the given fields value to that of the given interface value if it is possible to do
// so. If setttng the value fails an error is returned.
func setValue(field Field, value interface{}) error {
	fk := field.Value.Kind()
	vk := reflect.ValueOf(value).Kind()

	if u := unmarshaler(field); u != nil {
		return u.UnmarshalGoFig(value)
	}

	switch field.Value.Kind() {
	case reflect.String:
		return setString(field, value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return setInt64(field, value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return setUint64(field, value)
	case reflect.Float32, reflect.Float64:
		return setFloat64(field, value)
	case reflect.Slice, reflect.Array:
		return setSlice(field, value)
	}

	return ErrInvalidConversion{
		To:   fk,
		From: vk,
	}
}

// setString sets the fields value to a string.
func setString(field Field, value interface{}) error {
	v, ok := value.(string)
	if !ok {
		return ErrSetValue{
			Field: field.Value,
			Value: reflect.ValueOf(value),
		}
	}

	field.SetString(v)

	return nil
}

// setInt64 sets the fields value to an integer type.
func setInt64(field Field, value interface{}) error {
	var i int64

	switch t := value.(type) {
	case string:
		v, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return err
		}

		i = v
	case int:
		i = int64(t)
	case int8:
		i = int64(t)
	case int16:
		i = int64(t)
	case int32:
		i = int64(t)
	case int64:
		i = int64(t)
	case float32:
		i = int64(t)
	case float64:
		i = int64(t)
	default:
		return ErrSetValue{
			Field: field.Value,
			Value: reflect.ValueOf(value),
		}
	}

	return field.SetInt(i)
}

// setUint64 sets the fields value to an unsigned integer type.
func setUint64(field Field, value interface{}) error {
	var i uint64

	switch t := value.(type) {
	case string:
		v, err := strconv.ParseUint(t, 10, 64)
		if err != nil {
			return err
		}

		i = v
	case int:
		i = uint64(t)
	case int8:
		i = uint64(t)
	case int16:
		i = uint64(t)
	case int32:
		i = uint64(t)
	case int64:
		i = uint64(t)
	case float32:
		i = uint64(t)
	case float64:
		i = uint64(t)
	default:
		return ErrSetValue{
			Field: field.Value,
			Value: reflect.ValueOf(value),
		}
	}

	return field.SetUint(i)
}

// setFloat64 sets the fields value to an float type.
func setFloat64(field Field, value interface{}) error {
	var i float64

	switch t := value.(type) {
	case string:
		v, err := strconv.ParseFloat(t, 64)
		if err != nil {
			return err
		}

		i = v
	case float32:
		i = float64(t)
	case float64:
		i = float64(t)
	default:
		return ErrSetValue{
			Field: field.Value,
			Value: reflect.ValueOf(value),
		}
	}

	return field.SetFloat(i)
}

// setSlice sets the fields value to an slice.
func setSlice(field Field, value interface{}) error {
	ft := field.Value.Type()
	vv := reflect.ValueOf(value)
	if vv.Kind() != reflect.Array && vv.Kind() != reflect.Slice {
		return ErrInvalidValue{
			Type: ft,
		}
	}

	s := reflect.MakeSlice(reflect.SliceOf(ft.Elem()), vv.Len(), vv.Cap())

	for i := 0; i < vv.Len(); i++ {
		e := reflect.New(ft.Elem())
		f := Field{fmt.Sprintf("%s.%d", field.Key, i), e.Elem()}
		if err := setValue(f, vv.Index(i).Interface()); err != nil {
			return err
		}

		s.Index(i).Set(e.Elem())
	}

	field.Set(s)

	return nil
}

// setMap sets a field to a map, also handles nested maps.
func setMap(field Field, key string, value interface{}) error {
	fv := field.Value
	ft := field.Value.Type()

	if fv.IsNil() {
		if ft.Key().Kind() != reflect.String {
			return ErrInvalidValue{
				Type: ft.Key(),
			}
		}

		field.Set(reflect.MakeMap(reflect.MapOf(ft.Key(), ft.Elem())))
	}

	// Nested map
	if ft.Elem().Kind() == reflect.Map {
		elms := strings.Split(key, ".")
		key = strings.Join(elms[:len(elms)-1], ".")

		if key == "" {
			return nil
		}

		v := reflect.New(ft.Elem())
		f := Field{key, v.Elem()}
		if err := setMap(f, key, value); err != nil {
			return err
		}

		fv.SetMapIndex(reflect.ValueOf(elms[0]), v.Elem())

		return nil
	}

	if reflect.ValueOf(value).Kind() != ft.Elem().Kind() {
		return ErrInvalidConversion{
			From: reflect.ValueOf(value).Kind(),
			To:   ft.Elem().Kind(),
		}
	}

	v := reflect.New(ft.Elem())
	f := Field{key, v.Elem()}
	if err := setValue(f, value); err != nil {
		return err
	}

	fv.SetMapIndex(reflect.ValueOf(key), v.Elem())

	return nil
}

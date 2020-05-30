package gofig

import (
	"reflect"
	"strconv"
)

// A Field represents a type that we can set a value for based on its key.
type Field interface {
	// Set sets the fields value to the provided interface provided the types match.
	Set(interface{}) error
	// Key returns the full key path to the field
	Key() string
	// Value returns the destination reflect.Value values will be set into.
	Value() reflect.Value
}

// Fields holds a map of keys to fields.
type Fields map[string]Field

// Set sets a keys field value
func (f Fields) Set(k string, v Field) {
	f[k] = v
}

// A Field holds the fields struct path and reflected value.
type field struct {
	key   string // foo.bar.baz
	value reflect.Value
}

func newField(k string, v reflect.Value) *field {
	return &field{
		key:   k,
		value: v,
	}
}

func (f *field) Set(value interface{}) error {
	return set(f.value, value)
}

func (f *field) Key() string {
	return f.key
}

func (f *field) Value() reflect.Value {
	return f.value
}

// mapField embedded field wrapping map key values allowing setting map fields to be the same as
// setting struct fields.
type mapField struct {
	*field

	mk reflect.Value // Map key to set values in
	mp reflect.Value // Map to set values in
}

func newMapField(k, mk string, mp reflect.Value) *mapField {
	return &mapField{
		field: newField(k, reflect.New(mp.Type().Elem()).Elem()),

		mk: reflect.ValueOf(mk),
		mp: mp,
	}
}

func (f *mapField) Set(v interface{}) error {
	// TODO: ensure v == f.value
	if err := f.field.Set(v); err != nil {
		return err
	}

	// Set the map index to the given value
	f.mp.SetMapIndex(f.mk, f.value)

	return nil
}

func set(field reflect.Value, value interface{}) error {
	if u := unmarshaler(field); u != nil {
		return u.UnmarshalGoFig(value)
	}

	switch field.Kind() {
	case reflect.Ptr:
		return set(field.Elem(), value)
	case reflect.String:
		return setString(field, value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return setInt(field, value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return setUint(field, value)
	case reflect.Float32, reflect.Float64:
		return setFloat(field, value)
	case reflect.Slice, reflect.Array:
		return setSlice(field, value)
	}

	return ErrInvalidConversion{
		To:   field.Kind(),
		From: reflect.ValueOf(value).Kind(),
	}
}

// setString sets the fields value to a string.
func setString(field reflect.Value, value interface{}) error {
	v, ok := value.(string)
	if !ok {
		return ErrSetValue{
			Field: field,
			Value: reflect.ValueOf(value),
		}
	}

	field.SetString(v)

	return nil
}

// setInt sets the fields value to an integer.
func setInt(field reflect.Value, value interface{}) error {
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
			Field: field,
			Value: reflect.ValueOf(value),
		}
	}

	if field.OverflowInt(i) {
		return ErrSetValue{
			Field: field,
			Value: reflect.ValueOf(i),
		}
	}

	field.SetInt(i)

	return nil
}

// setUint64 sets the fields value to an unsigned integer type.
func setUint(field reflect.Value, value interface{}) error {
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
			Field: field,
			Value: reflect.ValueOf(value),
		}
	}

	if field.OverflowUint(i) {
		return ErrSetValue{
			Field: field,
			Value: reflect.ValueOf(i),
		}
	}

	field.SetUint(i)

	return nil
}

// setFloat sets the fields value to the a float.
func setFloat(field reflect.Value, value interface{}) error {
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
			Field: field,
			Value: reflect.ValueOf(value),
		}
	}

	if field.OverflowFloat(i) {
		return ErrSetValue{
			Field: field,
			Value: reflect.ValueOf(i),
		}
	}

	field.SetFloat(i)

	return nil
}

// setSlice sets the fields value to the given slice.
func setSlice(field reflect.Value, value interface{}) error {
	ft := field.Type()
	vv := reflect.ValueOf(value)

	if vv.Kind() != reflect.Array && vv.Kind() != reflect.Slice {
		return ErrInvalidValue{
			Type: ft,
		}
	}

	s := reflect.MakeSlice(reflect.SliceOf(ft.Elem()), vv.Len(), vv.Cap())

	for i := 0; i < vv.Len(); i++ {
		e := reflect.New(ft.Elem())
		if err := set(e, vv.Index(i).Interface()); err != nil {
			return err
		}

		s.Index(i).Set(e.Elem())
	}

	field.Set(s)

	return nil
}

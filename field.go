package gofig

import (
	"fmt"
	"reflect"
)

// Fields holds a map of keys to fields.
type Fields map[string]Field

// Set sets a field to the given val on the map.
func (f Fields) Set(key string, val Field) {
	f[key] = val
}

// Delete deletes a field from the map.
func (f Fields) Delete(key string) {
	delete(f, key)
}

// A Field holds the fields struct path and reflected value.
type Field struct {
	Key   string
	Value reflect.Value
}

func (f Field) String() string {
	return fmt.Sprintf("<%s: %s>", f.Key, f.Value.Type().Name())
}

// Set sets the fields value to the given value.
func (f Field) Set(x reflect.Value) {
	f.Value.Set(x)
}

// SetString sets the fields value to the given string.
func (f Field) SetString(x string) {
	f.Value.SetString(x)
}

// SetInt sets the fields value to the given int.
func (f Field) SetInt(x int64) error {
	if f.Value.OverflowInt(x) {
		return ErrSetValue{
			Field: f.Value,
			Value: reflect.ValueOf(x),
		}
	}

	f.Value.SetInt(x)

	return nil
}

// SetUint sets the fields value to the given uint.
func (f Field) SetUint(x uint64) error {
	if f.Value.OverflowUint(x) {
		return ErrSetValue{
			Field: f.Value,
			Value: reflect.ValueOf(x),
		}
	}

	f.Value.SetUint(x)

	return nil
}

// SetFloat sets the fields value to the given float.
func (f Field) SetFloat(x float64) error {
	if f.Value.OverflowFloat(x) {
		return ErrSetValue{
			Field: f.Value,
			Value: reflect.ValueOf(x),
		}
	}

	f.Value.SetFloat(x)

	return nil
}

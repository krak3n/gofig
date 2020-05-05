package gofig

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

// Gofig default configuration.
const (
	DefaultStructTag = "gofig"
)

// Unmarshaler is the interface implemented by types that can unmarshal a values themselves.
type Unmarshaler interface {
	UnmarshalGoFig(value interface{}) error
}

// A Logger can print log items
type Logger interface {
	Print(values ...interface{})
	Printf(format string, values ...interface{})
}

// A LoggerFunc is an adapter function allowing regular methods to act as Loggers.
type LoggerFunc func(v ...interface{})

// Print calls the wrapped fn.
func (fn LoggerFunc) Print(v ...interface{}) {
	fn(v...)
}

// Printf calls the wrapped fn.
func (fn LoggerFunc) Printf(format string, v ...interface{}) {
	fn(fmt.Sprintf(format, v...))
}

// DefaultLogger returns a standard lib logger.
func DefaultLogger() Logger {
	return LoggerFunc(func(v ...interface{}) {
		log.Println(v...)
	})
}

// NopLogger logs nothing.
func NopLogger() Logger {
	return LoggerFunc(func(...interface{}) {})
}

// An Option configures gofig.
type Option interface {
	apply(*Config)
}

// An OptionFunc is an adapter allowing regular methods to act as Option's.
type OptionFunc func(*Config)

func (fn OptionFunc) apply(c *Config) {
	fn(c)
}

// SetLogger sets gofig's logger.
func SetLogger(l Logger) Option {
	return OptionFunc(func(c *Config) {
		c.log = l
	})
}

// WithNopLogger sets gofig's logger to a no-op logger.
func WithNopLogger() Option {
	return OptionFunc(func(c *Config) {
		c.log = NopLogger()
	})
}

// Config parses configuration from one or more sources.
type Config struct {
	log    Logger
	fields map[string]reflect.Value
}

// New constructs a new Config
func New(dst interface{}, opts ...Option) (*Config, error) {
	fields, err := parse(dst)
	if err != nil {
		return nil, err
	}

	c := &Config{
		log:    DefaultLogger(),
		fields: fields,
	}

	for _, opt := range opts {
		opt.apply(c)
	}

	return c, nil
}

// Parse parses the given parsers in order. If any one parser fails an error will be returned.
func (c *Config) Parse(parsers ...Parser) error {
	for _, p := range parsers {
		if err := c.parse(p); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) parse(p Parser) error {
	ch, err := p.Values()
	if err != nil {
		return err
	}

	for {
		fn, ok := <-ch
		if !ok {
			return nil // Done
		}

		key, val := fn()
		field, ok := c.fields[strings.ToLower(key)]
		if !ok {
			c.log.Print("key:", key, "value:", val)
			continue
		}

		c.log.Printf("key: %s values: %s", key, val)

		if err := setValue(field, val); err != nil {
			return err
		}
	}
}

func setValue(field reflect.Value, value interface{}) error {
	fk := field.Kind()
	vk := reflect.ValueOf(value).Kind()

	if u := unmarshaler(field); u != nil {
		return u.UnmarshalGoFig(value)
	}

	switch field.Kind() {
	case reflect.String:
		return setString(field, value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return setInt64(field, value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return setUint64(field, value)
	case reflect.Float32, reflect.Float64:
		return setFloat64(field, value)
	case reflect.Slice, reflect.Array:
		// TODO: slice / array
	case reflect.Map:
		// TODO: map
	case reflect.Struct:
		// TODO: struct
	}

	return ErrInvalidConversion{
		To:   fk,
		From: vk,
	}
}

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

func setString(field reflect.Value, value interface{}) error {
	v, ok := value.(string)
	if !ok {
		return nil // TODO: error
	}

	field.SetString(v)

	return nil
}

func setInt64(field reflect.Value, value interface{}) error {
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
	default:
		return nil // TODO: error type
	}

	if field.OverflowInt(i) {
		return nil // TODO: error type
	}

	field.SetInt(i)

	return nil
}

func setUint64(field reflect.Value, value interface{}) error {
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
	default:
		return nil // TODO: error type
	}

	if field.OverflowUint(i) {
		return nil // TODO: error type
	}

	field.SetUint(i)

	return nil
}

func setFloat64(field reflect.Value, value interface{}) error {
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
		return nil // TODO: error type
	}

	if field.OverflowFloat(i) {
		return nil // TODO: error type
	}

	field.SetFloat(i)

	return nil
}

func parse(v interface{}) (map[string]reflect.Value, error) {
	rt := reflect.TypeOf(v)
	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return nil, ErrInvalidValue{reflect.TypeOf(v)}
	}

	if rv.Elem().Kind() != reflect.Struct {
		return nil, ErrInvalidValue{reflect.TypeOf(v)}
	}

	rv = rv.Elem()
	rt = rt.Elem()

	fields := make(map[string]reflect.Value)

	flatten(rv, rt, "", fields)

	return fields, nil
}

func flatten(rv reflect.Value, rt reflect.Type, key string, fields map[string]reflect.Value) {
	for i := 0; i < rv.NumField(); i++ {
		fv := rv.Field(i)
		ft := rt.Field(i)

		if fv.CanSet() {
			var t tag

			if v, ok := ft.Tag.Lookup(DefaultStructTag); ok {
				t = parseTag(v)
			} else {
				t = tag{
					name: ft.Name,
				}
			}

			name := strings.Trim(strings.Join(append(strings.Split(key, "."), t.name), "."), ".")

			if fv.Kind() == reflect.Struct {
				flatten(fv, ft.Type, name, fields)
			} else {
				fields[name] = fv
			}
		}
	}
}

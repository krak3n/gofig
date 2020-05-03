package gofig

import (
	"log"
	"reflect"
	"strings"
)

// Gofig default configuration.
const (
	DefaultStructTag = "gofig"
)

// A Logger can print log items
type Logger interface {
	Print(v ...interface{})
}

// A LoggerFunc is an adapter function allowing regular methods to act as Loggers.
type LoggerFunc func(v ...interface{})

// Print calls the wrapped fn.
func (fn LoggerFunc) Print(v ...interface{}) {
	fn(v...)
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

		c.log.Print("key:", key, "value:", val)

		if field.CanSet() {
			// TODO: more types
			// TODO: nested structured
			// TODO: Unmarshal interface for custom types
			switch field.Kind() {
			case reflect.String:
				if v, ok := val.(string); ok {
					field.SetString(v)
				}

			default:
				// TODO: return error
			}

			continue
		}

		// TODO: return error
	}
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

		if fv.CanInterface() {
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

package gofig

import (
	"reflect"
	"strings"
)

// Gofig default configuration.
const (
	DefaultStructTag = "gofig"
)

// Config parses configuration from one or more sources.
type Config struct {
	logger Logger
	debug  bool
	fields map[string]reflect.Value
}

// New constructs a new Config
func New(dst interface{}, opts ...Option) (*Config, error) {
	fields, err := load(dst)
	if err != nil {
		return nil, err
	}

	c := &Config{
		logger: DefaultLogger(),
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

// log returns a logger if debug is true
func (c *Config) log() Logger {
	if c.debug {
		return c.logger
	}

	return NopLogger()
}

// parse parses an single parser.
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

		if field, ok := c.fields[strings.ToLower(key)]; ok {
			c.log().Printf("set key %s to %v", key, val)

			if err := setValue(field, val); err != nil {
				return err
			}

			continue
		}

		// If the field was not found this could be a map element value.
		// This will always be the leaf node.
		// First we find the root map at the top of the stack.
		if path, m, ok := c.mapRoot(key); ok {
			c.log().Printf("set key %s to %v", key, val)

			key = strings.Trim(strings.Replace(key, path, "", -1), ".")
			if err := setMap(m, key, val); err != nil {
				return err
			}

			continue
		}

		c.log().Printf("ignoring %s of value %v", key, val)
	}
}

func (c *Config) mapRoot(key string) (string, reflect.Value, bool) {
	var v reflect.Value

	elms := strings.Split(key, ".")

	key = strings.Join(elms[:len(elms)-1], ".")
	if key == "" {
		return key, v, false
	}

	field, ok := c.fields[key]
	if !ok {
		return c.mapRoot(key)
	}

	if field.Kind() != reflect.Map {
		return key, v, false
	}

	return key, field, true
}

// load validates the given value ensuring it is a pointer to a struct. Once validated the struct
// fields will be flattened into a single map where the key is path to a field and the value
// a reflect.Value.
func load(v interface{}) (map[string]reflect.Value, error) {
	rt := reflect.TypeOf(v)
	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return nil, ErrInvalidValue{reflect.TypeOf(v)}
	}

	if rv.Elem().Kind() != reflect.Struct {
		return nil, ErrInvalidValue{reflect.TypeOf(v)}
	}

	fields := make(map[string]reflect.Value)

	flatten(rv.Elem(), rt.Elem(), "", fields)

	return fields, nil
}

// flatten recursively flattens a struct.
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

			switch fv.Kind() {
			case reflect.Struct:
				flatten(fv, ft.Type, name, fields)
			default:
				fields[strings.ToLower(name)] = fv
			}
		}
	}
}

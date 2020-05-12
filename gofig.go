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
	fields Fields
	kfmtr  Formatter
}

// New constructs a new Config
func New(dst interface{}, opts ...Option) (*Config, error) {
	t := reflect.TypeOf(dst)
	v := reflect.ValueOf(dst)

	if t.Kind() != reflect.Ptr || v.IsNil() {
		return nil, ErrInvalidValue{reflect.TypeOf(v)}
	}

	if v.Elem().Kind() != reflect.Struct {
		return nil, ErrInvalidValue{reflect.TypeOf(v)}
	}

	c := &Config{
		logger: DefaultLogger(),
		fields: make(Fields),
		kfmtr:  KeyFormatter(LowerCaseFormatter()),
	}

	for _, opt := range opts {
		opt.apply(c)
	}

	c.flatten(v.Elem(), t.Elem(), "")

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

		// Call the function passed on the channel returnin key value pair
		key, val := fn()
		key = c.kfmtr.Format(key)

		// Lookup the key
		field, ok := c.fields[key]
		if !ok {
			// If the field was not found this could be a map element value.
			// This will always be the leaf node.
			if field, ok = c.mapRoot(key); ok {
				key = strings.Trim(strings.Replace(key, field.Key, "", -1), ".")
				if err := setMap(field, key, val); err != nil {
					return err
				}
			} else {
				c.log().Printf("%s key not found", key)
			}

			continue
		}

		// Attempt to set the fields value
		if err := setValue(field, val); err != nil {
			return err
		}
	}
}

// flatten recursively flattens a struct.
func (c *Config) flatten(rv reflect.Value, rt reflect.Type, key string) {
	for i := 0; i < rv.NumField(); i++ {
		fv := rv.Field(i)
		ft := rt.Field(i)

		if fv.CanSet() {
			tag := TagFromStructField(ft, DefaultStructTag)

			path := c.kfmtr.Format(strings.Trim(strings.Join(append(strings.Split(key, "."), tag.Name), "."), "."))

			c.log().Printf("<Field %s kind:%s path:%s tag:%s>", ft.Name, fv.Kind(), path, tag)

			switch fv.Kind() {
			case reflect.Struct:
				c.flatten(fv, ft.Type, path)
			default:
				c.fields.Set(path, Field{path, fv})
			}
		}
	}
}

// mapRoot recursively looks for a root map for the given key.
func (c *Config) mapRoot(key string) (Field, bool) {
	var f Field

	elms := strings.Split(key, ".")

	key = strings.Join(elms[:len(elms)-1], ".")
	if key == "" {
		return f, false
	}

	field, ok := c.fields[key]
	if !ok {
		return c.mapRoot(key)
	}

	if field.Value.Kind() != reflect.Map {
		return f, false
	}

	return field, true
}

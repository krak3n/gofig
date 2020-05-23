package gofig

import (
	"reflect"
	"strings"
)

// Gofig default configuration.
const (
	DefaultStructTag = "gofig"
)

// Loader parses configuration from one or more sources.
type Loader struct {
	// flattened map of field keys to struct reflect values
	fields Fields

	// Configurable options
	logger       Logger
	debug        bool
	keyFormatter Formatter
	structTag    string
}

// New constructs a new Loader
func New(dst interface{}, opts ...Option) (*Loader, error) {
	t := reflect.TypeOf(dst)
	v := reflect.ValueOf(dst)

	if t.Kind() != reflect.Ptr || v.IsNil() {
		return nil, ErrInvalidValue{reflect.TypeOf(v)}
	}

	if v.Elem().Kind() != reflect.Struct {
		return nil, ErrInvalidValue{reflect.TypeOf(v)}
	}

	l := &Loader{
		logger:       DefaultLogger(),
		fields:       make(Fields),
		keyFormatter: KeyFormatter(LowerCaseFormatter()),
		structTag:    DefaultStructTag,
	}

	for _, opt := range opts {
		opt.apply(l)
	}

	l.flatten(v.Elem(), t.Elem(), "")

	return l, nil
}

// Parse parses the given parsers in order. If any one parser fails an error will be returned.
func (l *Loader) Parse(parsers ...Parser) error {
	for _, p := range parsers {
		if err := l.parse(p); err != nil {
			return err
		}
	}

	return nil
}

// log returns a logger if debug is true
func (l *Loader) log() Logger {
	if l.debug {
		return l.logger
	}

	return NopLogger()
}

// parse parses an single parser.
func (l *Loader) parse(p Parser) error {
	ch, err := p.Values()
	if err != nil {
		return err
	}

	for {
		fn, ok := <-ch
		if !ok {
			return nil // Done
		}

		// Call the function passed on the channel returning key value pair
		key, val := fn()
		key = l.keyFormatter.Format(key)

		// Lookup the key
		field, ok := l.fields[key]
		if !ok {
			// If the field was not found this could be a map element value.
			// This will always be the leaf node.
			if field, ok = l.mapRoot(key); ok {
				key = strings.Trim(strings.Replace(key, field.Key, "", -1), ".")
				if err := setMap(field, key, val); err != nil {
					return err
				}
			} else {
				l.log().Printf("%s key not found", key)
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
func (l *Loader) flatten(rv reflect.Value, rt reflect.Type, key string) {
	for i := 0; i < rv.NumField(); i++ {
		fv := rv.Field(i)
		ft := rt.Field(i)

		if fv.CanSet() {
			tag := TagFromStructField(ft, l.structTag)

			path := l.keyFormatter.Format(strings.Trim(strings.Join(append(strings.Split(key, "."), tag.Name), "."), "."))

			l.log().Printf("<Field %s kind:%s path:%s tag:%s>", ft.Name, fv.Kind(), path, tag)

			switch fv.Kind() {
			case reflect.Struct:
				l.flatten(fv, ft.Type, path)
			default:
				l.fields.Set(path, Field{path, fv})
			}
		}
	}
}

// mapRoot recursively looks for a root map for the given key.
func (l *Loader) mapRoot(key string) (Field, bool) {
	var f Field

	elms := strings.Split(key, ".")

	key = strings.Join(elms[:len(elms)-1], ".")
	if key == "" {
		return f, false
	}

	field, ok := l.fields[key]
	if !ok {
		return l.mapRoot(key)
	}

	if field.Value.Kind() != reflect.Map {
		return f, false
	}

	return field, true
}

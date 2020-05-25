package gofig

import (
	"reflect"
	"strings"
	"sync"
)

// Gofig default configuration.
const (
	DefaultStructTag = "gofig"
)

// Loader parses configuration from one or more sources.
type Loader struct {
	// flattened map of field keys to struct reflect values
	fields Fields

	// notifiers we are currently watching
	notifiers []NotifyParser
	wg        sync.WaitGroup

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
		fields:    make(Fields),
		notifiers: make([]NotifyParser, 0),

		// Defaults
		logger: DefaultLogger(),
		keyFormatter: KeyFormatter(FormatterFunc(func(key string) string {
			return key
		})),
		structTag: DefaultStructTag,
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
	// Send the keys
	errCh := make(chan error, 1)
	keyCh := make(chan string, len(l.fields))

	go func() {
		close(errCh)

		if err := p.Keys(keyCh); err != nil {
			errCh <- err
		}
	}()

	for _, f := range l.fields {
		keyCh <- f.Key
	}

	close(keyCh)

	if err := <-errCh; err != nil {
		return err
	}

	// Get the 	values
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

		// Lookup the field
		field, ok := l.find(key)
		if ok {
			if field.Value.Kind() == reflect.Map {
				if err := setMap(field, key, val); err != nil {
					return err
				}

				continue
			}

			if err := setValue(field, val); err != nil {
				return err
			}
		} else {
			l.log().Printf("%s key not found", key)
		}
	}
}

// flatten recursively flattens a struct.
func (l *Loader) flatten(rv reflect.Value, rt reflect.Type, key string) {
	for i := 0; i < rv.NumField(); i++ {
		fv := rv.Field(i)
		ft := rt.Field(i)

		// TODO: embedded support
		if fv.CanSet() {
			tag := TagFromStructField(ft, l.structTag)

			k := strings.Trim(strings.Join(append(strings.Split(key, "."), tag.Name), "."), ".")

			l.log().Printf("<Field %s kind:%s key:%s tag:%s>", ft.Name, fv.Kind(), k, tag)

			switch fv.Kind() {
			case reflect.Struct:
				l.flatten(fv, ft.Type, k)
			default:
				l.fields.Set(k, Field{k, fv})
			}
		}
	}
}

func (l *Loader) find(key string) (Field, bool) {
	f, ok := l.fields[key]
	if ok {
		return f, ok
	}

	elms := strings.Split(key, ".")

	key = strings.Join(elms[:len(elms)-1], ".")
	if key == "" {
		return f, false
	}

	return l.find(key)
}

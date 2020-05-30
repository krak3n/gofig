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
	delimiter    string
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
		delimiter: ".",
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
	// Set the delimiter
	p.SetDelimeter(l.delimiter)

	// Send keys to the parser
	if err := l.sendKeys(p); err != nil {
		return nil
	}

	// Get the 	values
	ch, err := p.Values()
	if err != nil {
		return err
	}

	// Range over the channel until it's closed processing the returned key / values
	for fn := range ch {
		// Call the function passed on the channel returning key value pair
		key, val := fn()
		key = l.keyFormatter.Format(key)

		// Lookup the field
		field, ok := l.lookup(key)
		if !ok {
			l.log().Printf("%s key not found", key)
			continue
		}

		// Set the value on the field
		if err := field.Set(val); err != nil {
			return err
		}
	}

	return nil
}

// sends keys to the parser.
func (l *Loader) sendKeys(p Parser) error {
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
		keyCh <- f.Key()
	}

	close(keyCh)

	return <-errCh
}

// lookup finds a field by it's key. If it finds a field that is a map that map is initialised
// returning a field values can be set on.
func (l *Loader) lookup(key string) (Field, bool) {
	// Look up the field
	field, ok := l.find(key)
	if !ok {
		return nil, false
	}

	// Return the field if it is not a map
	if field.Value().Kind() == reflect.Map {

		// The field is a map, this could be a leaf node, init the map
		// Generate the map key path by removing the root key from the field key
		// e.g foo.bar.baz becomes baz where bar is a map ahd baz the map key
		mk := strings.Trim(strings.Replace(key, field.Key(), "", -1), l.delimiter)

		// Returns the leaf map that the value should be set into
		mv, err := l.initMap(field.Value(), mk)
		if err != nil {
			return nil, false
		}

		// Make a field we can set map index values on
		kp := strings.Split(mk, l.delimiter)
		field = newMapField(key, kp[len(kp)-1], mv)

		// Insert the field into the field map so we don't have to initMap again for this value
		l.fields.Set(key, field)
	}

	return field, true
}

// find recursively finds a field based on the key path until a field is found or the key is empty.
func (l *Loader) find(key string) (Field, bool) {
	// Look up the field
	field, ok := l.fields[key]
	if ok {
		return field, true
	}

	elms := strings.Split(key, l.delimiter)

	if key := strings.Join(elms[:len(elms)-1], l.delimiter); key != "" {
		return l.find(key)
	}

	return nil, false
}

// initMap initialises maps with zero values for the given keep. Also handles deeply nested maps.
// Returns the map to set values into.
func (l *Loader) initMap(elem reflect.Value, key string) (reflect.Value, error) {
	// If the elem value is nil, initialise a new map of the correct types and set it as the fields value
	if elem.IsNil() {
		if elem.Type().Key().Kind() != reflect.String {
			return reflect.Value{}, ErrInvalidValue{
				Type: elem.Type().Key(),
			}
		}

		elem.Set(reflect.MakeMap(reflect.MapOf(
			elem.Type().Key(),
			elem.Type().Elem())))
	}

	// If the maps value is of another map this is a nested map
	// We need too initialise a new map of the correct types and set the keys index to be that new
	// nested map.
	if elem.Type().Elem().Kind() == reflect.Map {
		// Split the key at the delimiter extracting the parent and children key elements
		elms := strings.Split(key, l.delimiter)
		parent, children := elms[0], elms[1:]

		// Remove the parent from the key, e.g foo.bar.baz becomes bar.baz
		key = strings.Join(children, l.delimiter)
		if key == "" {
			return reflect.Value{}, nil // error
		}

		// Check if we have this map index in the map, if not create a new value for the index of
		// the correct type - this will be another map.
		m := elem.MapIndex(reflect.ValueOf(parent))
		if !m.IsValid() {
			m = reflect.New(elem.Type().Elem()).Elem()
		}

		// As this is a map we now init that map
		field, err := l.initMap(m, key)
		if err != nil {
			return reflect.Value{}, err
		}

		// Set the map index of key parent to the value of the new map.
		elem.SetMapIndex(reflect.ValueOf(elms[0]), m)

		return field, nil
	}

	// Return the map for index value setting.
	return elem, nil
}

// flatten recursively flattens a struct.
func (l *Loader) flatten(rv reflect.Value, rt reflect.Type, key string) {
	for i := 0; i < rv.NumField(); i++ {
		fv := rv.Field(i)
		ft := rt.Field(i)

		// TODO: embedded support
		if fv.CanSet() {
			tag := TagFromStructField(ft, l.structTag)

			fk := l.keyFormatter.Format(
				strings.Trim(
					strings.Join(
						append(strings.Split(key, l.delimiter), tag.Name), l.delimiter),
					l.delimiter))

			l.log().Printf("<Field %s kind:%s key:%s tag:%s>", ft.Name, fv.Kind(), fk, tag)

			switch fv.Kind() {
			case reflect.Struct:
				l.flatten(fv, ft.Type, fk)
			default:
				l.fields.Set(fk, newField(fk, fv))
			}
		}
	}
}

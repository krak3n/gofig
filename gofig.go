package gofig

import (
	"log"
	"reflect"
)

// Gofig default configuration.
const (
	DefaultStructTag = "gofig"
)

// Config parses configuration from one or more sources.
type Config struct {
	fields map[string]*field
}

// New constructs a new Config
func New(dst interface{}) (*Config, error) {
	return &Config{}, nil
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
		log.Println("key:", key, "value:", val)
	}
}

type field struct {
	kind reflect.Kind
	ptr  reflect.Value
}

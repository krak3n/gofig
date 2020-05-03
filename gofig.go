package gofig

import (
	"log"
	"reflect"
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

// Config parses configuration from one or more sources.
type Config struct {
	log    Logger
	fields map[string]*field
}

// New constructs a new Config
func New(dst interface{}, opts ...Option) (*Config, error) {
	c := &Config{
		log: DefaultLogger(),
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
		c.log.Print("key:", key, "value:", val)
	}
}

type field struct {
	kind reflect.Kind
	ptr  reflect.Value
}

package env

import (
	"os"
	"strings"
)

// Filter filters environment varibales.
type Filter func(key, value string) bool

// A Splitter splits environment variable keys into a slice.
type Splitter func(key string) []string

// An Option configures the Parser.
type Option interface {
	apply(*Parser)
}

// An OptionFunc is an adapter allowing regular methods to act as Option's.
type OptionFunc func(p *Parser)

func (fn OptionFunc) apply(p *Parser) {
	fn(p)
}

// WithFilter sets the Filter function to use.
func WithFilter(f Filter) Option {
	return OptionFunc(func(p *Parser) {
		p.filter = f
	})
}

// WithSplitter sets the Splitter function to use.
func WithSplitter(s Splitter) Option {
	return OptionFunc(func(p *Parser) {
		p.splitter = s
	})
}

// DefaultFilter takes the environment variable key and values and returns true if either of those
// values are empty, filtering them out from the values passed back to GoFig.
func DefaultFilter(key, val string) bool {
	if key == "" || val == "" {
		return true
	}

	return false
}

// DefaultSplitter splits the environment variable key at underscore.
func DefaultSplitter(key string) []string {
	return strings.Split(key, "_")
}

// Parser parsers OS environment variables.
type Parser struct {
	filter   Filter
	splitter Splitter
}

// New constructs a new OS environment variables parser.
func New(opts ...Option) *Parser {
	p := &Parser{
		filter:   DefaultFilter,
		splitter: DefaultSplitter,
	}

	for _, opt := range opts {
		opt.apply(p)
	}

	return p
}

// Values returns a channel of funcs that return each environment variable key values.
func (p *Parser) Values() (<-chan func() (string, interface{}), error) {
	ch := make(chan func() (string, interface{}))

	go func() {
		defer close(ch)

		for _, env := range os.Environ() {
			key, val := split(env)

			if p.filter(key, val) {
				continue
			}

			key = strings.Join(p.splitter(key), ".")

			ch <- (func(key string, val interface{}) func() (string, interface{}) {
				return func() (string, interface{}) {
					return key, val
				}
			}(key, val))
		}
	}()

	return ch, nil
}

// split splits an environment string at the = separator returning the key value pair.
func split(env string) (string, string) {
	var (
		key string
		val string
	)

	for i := 0; i < len(env); i++ {
		if env[i] == '=' {
			key = env[:i]
			val = env[i+1:]
		}
	}

	return key, val
}

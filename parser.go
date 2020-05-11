package gofig

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// A Parser parses configuration.
type Parser interface {
	// Values returns a channel of functions that returns an individual key value pair.
	// Each key returned should be an absolute flattened path and the value being the keys value.
	// The Parser should close the channel once parsing the configuration is complete and no more
	// values are to be returned.
	// If there are any errors in parsing the config Values should return a nil channel and an
	// error.
	// This interface allows Parsers and GoFig to remain loosely coupled.
	//
	// Given parsing the following yaml:
	//
	//   foo:
	//     bar:
	//       baz: fizz
	//
	// The values returned by the function would be:
	//
	// * key: foo.bar.baz.
	// * value: fizz as a string.
	Values() (<-chan func() (key string, value interface{}), error)
}

// A ParserFunc is an adapter allowing regular methods to act as Parser's.
type ParserFunc func() (<-chan func() (key string, value interface{}), error)

// Values calls the wrapped fn returning it's values.
func (fn ParserFunc) Values() (<-chan func() (string, interface{}), error) {
	return fn()
}

// A ReaderParser parses configuration from an io.Reader.
type ReaderParser interface {
	Values(src io.ReadCloser) (<-chan func() (key string, value interface{}), error)
}

// An InMemoryParser holds key value pairs in memory implementing the Parser interface.
type InMemoryParser struct {
	values map[string]interface{}
}

// Add adds a value to the in memory values.
func (p *InMemoryParser) Add(k string, v interface{}) {
	p.values[k] = v
}

// Delete deletes a value.
func (p *InMemoryParser) Delete(k string) {
	delete(p.values, k)
}

// Values iterates over the in memory values returning then on the returned channel.
func (p *InMemoryParser) Values() (<-chan func() (string, interface{}), error) {
	ch := make(chan func() (string, interface{}))

	go func() {
		for k, v := range p.values {
			ch <- (func(key string, val interface{}) func() (string, interface{}) {
				return func() (string, interface{}) {
					return key, val
				}
			}(k, v))
		}

		close(ch)
	}()

	return ch, nil
}

// FromString parsers configuration from a string.
func FromString(parser ReaderParser, v string) Parser {
	return ParserFunc(func() (<-chan func() (string, interface{}), error) {
		return parser.Values(ioutil.NopCloser(strings.NewReader(v)))
	})
}

// FromBytes parsers configuration from a byte slice.
func FromBytes(parser ReaderParser, b []byte) Parser {
	return ParserFunc(func() (<-chan func() (string, interface{}), error) {
		return parser.Values(ioutil.NopCloser(bytes.NewReader(b)))
	})
}

// FromFile reads a file.
func FromFile(parser ReaderParser, path string) Parser {
	return ParserFunc(func() (<-chan func() (string, interface{}), error) {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		return parser.Values(f)
	})
}

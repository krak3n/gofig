package gofig

import (
	"bytes"
	"io"
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
	Values(src io.Reader) (<-chan func() (key string, value interface{}), error)
}

// ParseReadCloser parses a ReadCloser.
type ParseReadCloser struct {
	parser ReaderParser
	rc     io.ReadCloser
}

// Values returns a channel of a func that returns a key value pair.
func (p *ParseReadCloser) Values() (<-chan func() (key string, value interface{}), error) {
	return p.parser.Values(p.rc)
}

// Close closes the ReadCloser.
func (p *ParseReadCloser) Close() error {
	if p.rc != nil {
		return p.rc.Close()
	}

	return nil
}

// FromString parsers configuration from a string.
func FromString(parser ReaderParser, v string) Parser {
	return ParserFunc(func() (<-chan func() (string, interface{}), error) {
		return parser.Values(strings.NewReader(v))
	})
}

// FromBytes parsers configuration from a byte slice.
func FromBytes(parser ReaderParser, b []byte) Parser {
	return ParserFunc(func() (<-chan func() (string, interface{}), error) {
		return parser.Values(bytes.NewReader(b))
	})
}

// FromFile reads a file and returns a ParseReadCloser.
// Rember to defer Close().
func FromFile(parser ReaderParser, path string) (*ParseReadCloser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return &ParseReadCloser{
		parser: parser,
		rc:     f,
	}, nil
}

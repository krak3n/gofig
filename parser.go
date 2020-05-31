package gofig

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// Parsers stores a map of pointers to PrioritisedParser's.
type Parsers map[Parser]PrioritisedParser

// Add adds a Parser to the Parsers map returning a PrioritisedParser Parser. The priority of the
// parser is automatically set based on the size the parsers map.
func (p Parsers) Add(parser Parser) PrioritisedParser {
	prioritised, ok := p[parser]
	if ok {
		return prioritised
	}

	prioritised = PrioritiseParser(parser)
	prioritised.SetPriority(uint8(p.Len() + 1))

	p[parser] = prioritised

	return prioritised
}

// Get returns a PrioritisedParser from the Parsers map.
func (p Parsers) Get(parser Parser) PrioritisedParser {
	return p.Add(parser)
}

// Len returns the size of the parsers map.
func (p Parsers) Len() int {
	return len(p)
}

// A DelimeterSetter sets the key delimeter used for flattened keys, default is .
type DelimeterSetter interface {
	SetDelimeter(string)
}

// A Prioritiser prioritises a Parser.
type Prioritiser interface {
	SetPriority(uint8)
	Priority() uint8
}

// PrioritiseParser wraps a Parser so it can be prioritised.
func PrioritiseParser(p Parser) PrioritisedParser {
	return &prioritised{
		Parser: p,
	}
}

type prioritised struct {
	Parser

	priority uint8
}

func (p *prioritised) SetPriority(v uint8) {
	p.priority = v
}

func (p *prioritised) Priority() uint8 {
	return p.priority
}

// A Parser parses configuration.
type Parser interface {
	DelimeterSetter

	// Keys sends flattened keys (e.g foo.bar.fizz_buzz) to the parser. The Parser then can then decide, if
	// it wishes to format the key and store internal mapping or not.
	// This is useful for parsers like environment variables where keys such as foo.bar.fizz_buzz would need to be
	// converted too FOO_BAR_FIZZ_BUZZ with a mapping to the original key.
	// This allows us to  maintain case sensitivity in key lookups within the loader.
	// Most parsers such as YAML, TOML and JSON will not process these keys.
	Keys(keys <-chan string) error

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

// A PrioritisedParser is a Parser that has been prioritised.
type PrioritisedParser interface {
	Parser
	Prioritiser
}

// A ParseReadCloser parses configuration from an io.ReadCloser.
type ParseReadCloser interface {
	DelimeterSetter

	Values(src io.ReadCloser) (<-chan func() (key string, value interface{}), error)
}

// An InMemoryParser holds key value pairs in memory implementing the Parser interface.
type InMemoryParser struct {
	values   map[string]interface{}
	notify   bool
	notifyCh chan error
}

// NewInMemoryParser constructs a new InMemoryParser.
func NewInMemoryParser() *InMemoryParser {
	return &InMemoryParser{
		values:   make(map[string]interface{}),
		notifyCh: make(chan error),
	}
}

// Add adds a value to the in memory values.
func (p *InMemoryParser) Add(k string, v interface{}) {
	p.values[k] = v

	if p.notify {
		p.notifyCh <- nil
	}
}

// Delete deletes a value.
func (p *InMemoryParser) Delete(k string) {
	delete(p.values, k)
}

// SetDelimeter is a no-op.
func (p *InMemoryParser) SetDelimeter(string) {}

// Keys consumes the keys but does nothing with them.
func (p *InMemoryParser) Keys(c <-chan string) error {
	for {
		_, ok := <-c
		if !ok {
			return nil
		}
	}
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

// Notify notifies when in memory values have changed.
func (p *InMemoryParser) Notify() <-chan error {
	if p.notifyCh == nil {
		p.notifyCh = make(chan error)
	}

	p.notify = true

	return p.notifyCh
}

// Close does closes the notify channel.
func (p *InMemoryParser) Close() error {
	p.notify = false
	close(p.notifyCh)
	p.notifyCh = nil

	return nil
}

// ReadCloseParser parses config from io.ReadCloser's.
type ReadCloseParser struct {
	parser   ParseReadCloser
	src      io.ReadCloser
	priority uint8
}

// NewReadCloseParser constructs a new ReadCloseParser.
func NewReadCloseParser(parser ParseReadCloser, src io.ReadCloser) *ReadCloseParser {
	return &ReadCloseParser{
		parser: parser,
		src:    src,
	}
}

// SetDelimeter sets the parsers delimeter.
func (p *ReadCloseParser) SetDelimeter(d string) {
	p.parser.SetDelimeter(d)
}

// SetPriority sets the parsers priority.
func (p *ReadCloseParser) SetPriority(v uint8) {
	p.priority = v
}

// Priority returns parsers priority.
func (p *ReadCloseParser) Priority() uint8 {
	return p.priority
}

// Keys is a no-op key consumer.
func (p *ReadCloseParser) Keys(c <-chan string) error {
	for {
		_, ok := <-c
		if !ok {
			return nil
		}
	}
}

// Values returns values from the parser back to gofig.
func (p *ReadCloseParser) Values() (<-chan func() (string, interface{}), error) {
	return p.parser.Values(p.src)
}

// FromString parsers configuration from a string.
func FromString(parser ParseReadCloser, v string) Parser {
	return NewReadCloseParser(parser, ioutil.NopCloser(strings.NewReader(v)))
}

// FromBytes parsers configuration from a byte slice.
func FromBytes(parser ParseReadCloser, b []byte) Parser {
	return NewReadCloseParser(parser, ioutil.NopCloser(bytes.NewReader(b)))
}

// FileParser parsers configuration from a file.
type FileParser struct {
	parser   ParseReadCloser
	path     string
	priority uint8
}

// NewFileParser constructs a new FileParser.
func NewFileParser(parser ParseReadCloser, path string) *FileParser {
	return &FileParser{
		parser: parser,
		path:   path,
	}
}

// SetDelimeter sets the parsers delimeter.
func (p *FileParser) SetDelimeter(d string) {
	p.parser.SetDelimeter(d)
}

// SetPriority sets the parsers priority.
func (p *FileParser) SetPriority(v uint8) {
	p.priority = v
}

// Priority returns parsers priority.
func (p *FileParser) Priority() uint8 {
	return p.priority
}

// Keys is a no-op key consumer.
func (p *FileParser) Keys(c <-chan string) error {
	for {
		_, ok := <-c
		if !ok {
			return nil
		}
	}
}

// Values opens the file for reading and passed it to the parser to return values back to gofig.
func (p *FileParser) Values() (<-chan func() (string, interface{}), error) {
	f, err := os.Open(p.path)
	if err != nil {
		return nil, err
	}

	return p.parser.Values(f)
}

// FromFile reads a file.
func FromFile(parser ParseReadCloser, path string) Parser {
	return NewFileParser(parser, path)
}

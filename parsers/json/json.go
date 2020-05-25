package json

import (
	"encoding/json"
	"io"
	"reflect"
	"strings"
)

// Parser parses YAML documents.
type Parser struct {
	delimiter string
}

// New constructs a new Parser.
func New() *Parser {
	return &Parser{
		delimiter: ".",
	}
}

// SetDelimeter sets the key delimiter.
func (p *Parser) SetDelimeter(v string) {
	p.delimiter = v
}

// Values parses yaml configuration, iterating over each key value pair and returning them until
// parsing has been completed.
func (p *Parser) Values(src io.ReadCloser) (<-chan func() (string, interface{}), error) {
	var dst map[string]interface{}

	d := json.NewDecoder(src)
	if err := d.Decode(&dst); err != nil {
		return nil, err
	}

	ch := make(chan func() (string, interface{}))

	go func() {
		defer close(ch)
		p.recurse("", dst, ch)
	}()

	return ch, src.Close()
}

func (p *Parser) recurse(key string, m map[string]interface{}, ch chan func() (string, interface{})) {
	for k, v := range m {
		name := strings.Trim(strings.Join(append(strings.Split(key, p.delimiter), k), p.delimiter), p.delimiter)

		if reflect.ValueOf(v).Kind() == reflect.Map {
			p.recurse(name, v.(map[string]interface{}), ch)

			continue
		}

		ch <- (func(key string, val interface{}) func() (string, interface{}) {
			return func() (string, interface{}) {
				return key, val
			}
		}(name, v))
	}
}

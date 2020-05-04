package env

import (
	"os"
	"strings"
)

// Parser parsers OS environment variables.
type Parser struct {
	filter    Filterer
	formatter Formatter
}

// Values returns a channel of funcs that return each environment variable key values.
func (p *Parser) Values() (<-chan func() (string, interface{}), error) {
	ch := make(chan func() (string, interface{}))

	go func() {
		defer close(ch)

		for _, env := range os.Environ() {
			key, val := split(env)

			if p.filter.Filter(key, val) {
				continue
			}

			key = p.formatter.Format(key)

			key = strings.Join(strings.Split(key, "_"), ".")

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

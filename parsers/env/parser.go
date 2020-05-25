package env

import (
	"os"
	"strings"
)

// Parser parsers OS environment variables.
type Parser struct {
	prefix string
	suffix string

	delimiter string
	keys      map[string]string
}

// SetDelimeter sets the key delimiter.
func (p *Parser) SetDelimeter(v string) {
	p.delimiter = v
}

// Keys consumes the keys from the channel.
func (p *Parser) Keys(c <-chan string) error {
	// Range over the keys we need to look for and convert to env variables formats.
	for key := range c {
		// Break the key at the . delimiter
		elms := strings.Split(key, p.delimiter)

		// Add prefix / suffix
		elms = append([]string{p.prefix}, elms...)
		elms = append(elms, p.suffix)

		// Join the elements elms together at _
		env := strings.Trim(strings.ToUpper(strings.Join(elms, "_")), "_")

		// Store the env var to key mapping
		p.keys[env] = key
	}

	return nil
}

// Values returns a channel of funcs that return each environment variable key values.
func (p *Parser) Values() (<-chan func() (string, interface{}), error) {
	ch := make(chan func() (string, interface{}))

	go func() {
		defer close(ch)

		for _, env := range os.Environ() {
			// Split the environment variable at =
			name, val := split(env)

			// Lookup the key, if found, send the key and the value
			key, ok := p.keys[name]
			if ok {
				ch <- (func(key string, val interface{}) func() (string, interface{}) {
					return func() (string, interface{}) {
						return key, val
					}
				}(key, val))
			}
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

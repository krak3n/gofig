package env

import "strings"

// An Option configures the Parser.
type Option interface {
	apply(*Parser)
}

// An OptionFunc is an adapter allowing regular methods to act as Option's.
type OptionFunc func(p *Parser)

func (fn OptionFunc) apply(p *Parser) {
	fn(p)
}

// Options holds muliple Option. This also implements the Option interface.
type Options []Option

func (opts Options) apply(p *Parser) {
	for _, opt := range opts {
		opt.apply(p)
	}
}

// HasAndTrimPrefix will filter out environment variables where the key does not have the given
// prefix. It will also then strip the given prefix from the key.
// NOTE: This chains with pre-existing filters and formatters.
func HasAndTrimPrefix(prefix string) Option {
	return Options{
		HasPrefix(prefix),
		TrimPrefix(prefix),
	}
}

// HasPrefix filters out environment variables that do not have the given prefix.
// NOTE: This chains with pre-existing filters.
func HasPrefix(prefix string) Option {
	return OptionFunc(func(p *Parser) {
		p.filter = func(next Filterer) Filterer {
			return FilterFunc(func(key, value string) bool {
				if next.Filter(key, value) {
					return true
				}

				return !strings.HasPrefix(key, prefix)
			})
		}(p.filter)
	})
}

// TrimPrefix removes the given prefix from the environment variable key.
// NOTE: This chains with pre-existing formatters.
func TrimPrefix(prefix string) Option {
	return OptionFunc(func(p *Parser) {
		p.formatter = func(next Formatter) Formatter {
			return FormatterFunc(func(key string) string {
				return strings.TrimLeft(strings.Replace(key, prefix, "", -1), "_")
			})
		}(p.formatter)
	})
}

// HasAndTrimSuffix will filter out environment variables where the key does not have the given
// suffix. It will also then strip the given suffix from the key.
// NOTE: This chains with pre-existing filters and formatters.
func HasAndTrimSuffix(prefix string) Option {
	return Options{
		HasSuffix(prefix),
		TrimSuffix(prefix),
	}
}

// HasSuffix filters out environment variables that do not have the given suffix.
// NOTE: This chains with pre-existing filters.
func HasSuffix(suffix string) Option {
	return OptionFunc(func(p *Parser) {
		p.filter = func(next Filterer) Filterer {
			return FilterFunc(func(key, value string) bool {
				if next.Filter(key, value) {
					return true
				}

				return !strings.HasSuffix(key, suffix)
			})
		}(p.filter)
	})
}

// TrimSuffix removes the given suffix from the environment variable key.
// NOTE: This chains with pre-existing formatters.
func TrimSuffix(suffix string) Option {
	return OptionFunc(func(p *Parser) {
		p.formatter = func(next Formatter) Formatter {
			return FormatterFunc(func(key string) string {
				return strings.TrimRight(strings.Replace(key, suffix, "", -1), "_")
			})
		}(p.formatter)
	})
}

package env

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

// WithPrefix sets a prefix on the environment variable name.
func WithPrefix(prefix string) Option {
	return OptionFunc(func(p *Parser) {
		p.prefix = prefix
	})
}

// WithSuffix sets a suffix on the environment variable name.
func WithSuffix(suffix string) Option {
	return OptionFunc(func(p *Parser) {
		p.suffix = suffix
	})
}

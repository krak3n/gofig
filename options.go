package gofig

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
		c.logger = l
	})
}

// WithDebug enables debugging. Use SetLogger to customise the logging output.
func WithDebug() Option {
	return OptionFunc(func(c *Config) {
		c.debug = true
	})
}

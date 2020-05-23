package gofig

// An Option configures gofig.
type Option interface {
	apply(*Loader)
}

// An OptionFunc is an adapter allowing regular methods to act as Option's.
type OptionFunc func(*Loader)

func (fn OptionFunc) apply(c *Loader) {
	fn(c)
}

// SetLogger sets gofig's logger.
func SetLogger(l Logger) Option {
	return OptionFunc(func(c *Loader) {
		c.logger = l
	})
}

// SetKeyFormatter sets the formatter to be used for recieved keys from parsers.
func SetKeyFormatter(fmtr Formatter) Option {
	return OptionFunc(func(c *Loader) {
		c.kfmtr = KeyFormatter(fmtr)
	})
}

// WithDebug enables debugging. Use SetLogger to customise the logging output.
func WithDebug() Option {
	return OptionFunc(func(c *Loader) {
		c.debug = true
	})
}

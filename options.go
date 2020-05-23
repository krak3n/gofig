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
func SetLogger(v Logger) Option {
	return OptionFunc(func(l *Loader) {
		l.logger = v
	})
}

// SetKeyFormatter sets the formatter to be used for recieved keys from parsers.
func SetKeyFormatter(fmtr Formatter) Option {
	return OptionFunc(func(l *Loader) {
		l.keyFormatter = KeyFormatter(fmtr)
	})
}

// SetStructTag changes the struct tag gofig looks for on struct fields to the value provided.
func SetStructTag(t string) Option {
	return OptionFunc(func(l *Loader) {
		l.structTag = t
	})
}

// WithDebug enables debugging. Use SetLogger to customise the logging output.
func WithDebug() Option {
	return OptionFunc(func(l *Loader) {
		l.debug = true
	})
}

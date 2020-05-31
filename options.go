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

// SetLogger sets logger.
func SetLogger(v Logger) Option {
	return OptionFunc(func(l *Loader) {
		l.logger = v
	})
}

// SetKeyFormatter sets the formatter to be used for received keys from parsers.
func SetKeyFormatter(formatter Formatter) Option {
	return OptionFunc(func(l *Loader) {
		l.keyFormatter = formatter
	})
}

// SetStructTag changes the struct tag gofig looks for on struct fields to the value provided.
func SetStructTag(t string) Option {
	return OptionFunc(func(l *Loader) {
		l.structTag = t
	})
}

// SetEnforcePriority enable or disable parser priority enforcement.
func SetEnforcePriority(v bool) Option {
	return OptionFunc(func(l *Loader) {
		l.enforcePriority = v
	})
}

// WithDebug enables debugging. Use SetLogger to customise the logging output.
func WithDebug() Option {
	return OptionFunc(func(l *Loader) {
		l.debug = true
	})
}

package env

// A Formatter formats an environment variable key. This allows for the modification of an
// environment variable key, for example, replacing a prefix or suffix.
type Formatter interface {
	Format(key string) string
}

// A FormatterFunc is an adapter function allowing regular methods to act a Formatter.
type FormatterFunc func(key string) string

// Format calls the wrapped fn.
func (fn FormatterFunc) Format(key string) string {
	return fn(key)
}

// FormatrChain allows for the chaining for formatters.
type FormatrChain func(next Formatter) Formatter

// NopFormatter does nothing.
func NopFormatter() Formatter {
	return FormatterFunc(func(key string) string {
		return key
	})
}

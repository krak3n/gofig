package env

// A Filterer filters out unwanted environment variables from being parsed by gofig.
// Filters can be performed on both the key and value. The Filter function should return true if the
// environment variable should be ignored.
type Filterer interface {
	Filter(key, value string) bool
}

// A FilterFunc is an adapter allowing regular methods to act as a Filterer.
type FilterFunc func(key, value string) bool

// Filter calls the wrapped fn.
func (fn FilterFunc) Filter(key, value string) bool {
	return fn(key, value)
}

// EmptyFilter filters out environment variables where the key or value are empty strings.
// This filter will always be applied first.
func EmptyFilter() Filterer {
	return FilterFunc(func(key, value string) bool {
		return (key == "" || value == "")
	})
}

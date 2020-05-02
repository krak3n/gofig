package gofig

// Gofig default configuration.
const (
	DefaultStructTag = "gofig"
)

// A Parser parses configuration.
type Parser interface {
	// Values iterates over the configuration returning key value pairs
	// where the key is a absolute flattened key path and valye is the keys value. The Parser
	// should return an io.EOF error when parsing has completed, any other error value will cause
	// parsing to error.
	//
	// Given parsing the following yaml:
	//
	//   foo:
	//     bar:
	//       baz: fizz
	//
	// The values returned by the parser would be:
	//
	// * key would be foo.bar.baz.
	// * value would be fizz as a string.
	// * err would be an io.EOF as only one key/value pair should be returned by the Parser.
	Values() (key string, value interface{}, err error)
}

// A Notifier notifies via a channel if changes to configuration have occurred.
// Remember to check the error on the channel.
type Notifier interface {
	Notify() <-chan error
}

// A ParseNotifier can parse config and notify on changes to configuration.
type ParseNotifier interface {
	Parser
	Notifier
}

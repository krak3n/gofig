package env

// New constructs a new OS environment variables parser.
// Use Option methods to configure the parsers behaviour.
func New(opts ...Option) *Parser {
	p := &Parser{
		keys:      map[string]string{},
		delimiter: ".",
	}

	for _, opt := range opts {
		opt.apply(p)
	}

	return p
}

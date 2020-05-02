package gofig

// Must calls the given function. On any error a panic will be thrown.
// fn must be a function and must return one value which must be an error.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

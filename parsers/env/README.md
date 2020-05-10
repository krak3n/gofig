# Environment Variable Parser

[![Workflow Status][workflow-image]][workflow-url]
[![Go Version][goversion-image]][goversion-url]
[![Example][playground-image]][playground-url]
[![Documentation][doc-image]][doc-url]

This parser loads configuration from OS Environment Variables.

## Example

Click the Playground badge above to see the example running in the Go Playground.

``` go
package main

import (
	"fmt"
	"os"

	"go.krak3n.codes/gofig"
	"go.krak3n.codes/gofig/parsers/env"
)

// Config is our configuration structure.
type Config struct {
	Foo  string `gofig:"foo"`
	Bar  string `gofig:"bar"`
	Fizz struct {
		Buzz string `gofig:"buzz"`
	} `gofig:"fizz"`
}

func main() {
	// Initialise Config
	var cfg Config

	// Set environment variables
	os.Setenv("GOFIG_FOO", "foo")
	os.Setenv("GOFIG_BAR", "bar")
	os.Setenv("GOFIG_FIZZ_BUZZ", "buzz")

	// Initialise gofig with the struct values will be parsed into
	gfg, err := gofig.New(&cfg)
	gofig.Must(err)

	// Parse the environment variables
	// This will filter out environment variables that do not have the given prefix and also trim
	// the prefix from the environment variable key.
	gofig.Must(gfg.Parse(env.New(
		env.HasAndTrimPrefix("GOFIG"),
	)))

	// Use the config
	fmt.Println("Foo:", cfg.Foo)             // foo
	fmt.Println("Bar:", cfg.Bar)             // bar
	fmt.Println("Fizz.Buzz:", cfg.Fizz.Buzz) // buzz
}
```

[workflow-image]: https://github.com/krak3n/gofig/workflows/Environment%20Variable%20Parser/badge.svg
[workflow-url]: https://github.com/krak3n/gofig/actions?query=workflow%3A%22Environment+Variable+Parser%22
[goversion-image]: https://img.shields.io/badge/Go-1.13+-00ADD8.svg
[goversion-url]: https://golang.org/
[playground-image]: https://img.shields.io/badge/Example-play.golang.org-00ADD8.svg
[playground-url]: https://play.golang.org/p/atkM_FbS0fq
[doc-image]: https://img.shields.io/badge/Documentation-pkg.go.dev-00ADD8.svg
[doc-url]: https://pkg.go.dev/go.krak3n.codes/gofig/parsers/env

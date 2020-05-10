# TOML Parser

[![Workflow Status][workflow-image]][workflow-url]
[![Go Version][goversion-image]][goversion-url]
[![Example][playground-image]][playground-url]
[![Documentation][doc-image]][doc-url]

This parser loads `toml` formatted configuration from an `io.ReadCloser`.

## Example

Click the Playground badge above to see the example running in the Go Playground.

``` go
package main

import (
	"fmt"

	"go.krak3n.codes/gofig"
	"go.krak3n.codes/gofig/parsers/toml"
)

type Config struct {
	Foo  string `gofig:"foo"`
	Bar  int    `gofig:"bar"`
	Fizz struct {
		Buzz string `gofig:"buzz"`
	} `gofig:"fizz"`
}

const blob = `
foo = "bar"
bar = 12

[fizz]
buzz = "fizz"`

func main() {
	var cfg Config

	// Initialise gofig with the struct config values will be placed into
	gfg, err := gofig.New(&cfg)
	gofig.Must(err)

	// Parse
	gofig.Must(gfg.Parse(gofig.FromString(toml.New(), blob)))

	fmt.Println(fmt.Sprintf("%+v", cfg))
}
```

[workflow-image]: https://github.com/krak3n/gofig/workflows/TOML%20Parser/badge.svg
[workflow-url]: https://github.com/krak3n/gofig/actions?query=workflow%3A%22TOML+Parser%22
[goversion-image]: https://img.shields.io/badge/Go-1.13+-00ADD8.svg
[goversion-url]: https://golang.org/
[playground-image]: https://img.shields.io/badge/Example-play.golang.org-00ADD8.svg
[playground-url]: https://play.golang.org/p/NUGuKxfNuot
[doc-image]: https://img.shields.io/badge/Documentation-pkg.go.dev-00ADD8.svg
[doc-url]: https://pkg.go.dev/go.krak3n.codes/gofig/parsers/toml

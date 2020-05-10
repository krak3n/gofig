# 🗃️ GoFig

[![Workflow Status][workflow-image]][workflow-url]
[![Go Version][goversion-image]][goversion-url]
[![Documentation][doc-image]][doc-url]

GoFig is a configuration loading library for Go. It aims to provide a simple and intuitive API that is
unopinionated for all your configuration loading needs.

* **Status**: PoC (Proof of Concept)

## Example

``` go
package main

import (
	"fmt"

	"go.krak3n.codes/gofig"
	"go.krak3n.codes/gofig/parsers/env"
)

type Config  struct {
	Foo struct {
		Bar string `gofig:"bar"`
	} `gofig:"foo"`
	Fizz struct {
		Buzz string `gofig:"buzz"`
	} `gofig:"fizz"`
}

func main() {
	var cfg Config

	// Initialise gofig with the destination struct
	gfg, err := gofig.New(&cfg)
	gofig.Must(err)

	// Parse the yaml file and then the envs
	gofig.Must(gfg.Parse(
		gofig.FromFile(yaml.New(), "./config.yaml"),
		env.New(env.HasAndTrimPrefix("GOFIG")),
	))

	fmt.Println(fmt.Sprintf("%+v", cfg))
}
```

## Parsers

GoFig implements it's parsers as sub modules. Currently it supports:

* [Environment Variables][env-url]
* [YAML][yaml-url]

### Planned

* TOML
* JSON

[goversion-image]: https://img.shields.io/badge/Go-1.13+-00ADD8.svg
[goversion-url]: https://golang.org/
[workflow-image]: https://github.com/krak3n/gofig/workflows/GoFig/badge.svg
[workflow-url]: https://github.com/krak3n/gofig/actions?query=workflow%3AGoFig
[doc-image]: https://img.shields.io/badge/Documentation-pkg.go.dev-00ADD8.svg
[doc-url]: https://pkg.go.dev/go.krak3n.codes/gofig
[env-url]: ./parsers/env
[yaml-url]: ./parsers/yaml

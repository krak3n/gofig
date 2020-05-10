# 🗃️ GoFig

[![Test Workflow][actions-image]][actions-url]
[![Go 1.12+][goversion-image]][goversion-url]
[![Documentation][gofig-doc-image]][gofig-doc-url]

GoFig is a configuration loading library for Go. It aims to provide a simple and intuitive API that is
unopinionated for your configuration loading needs.

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

## Examples

* [Environment Variables](examples/001_env)
* [YAML](examples/002_yaml)
* [Multi Source](examples/003_multisource)

[goversion-image]: https://img.shields.io/badge/Go-1.13+-00ADD8.svg
[goversion-url]: https://golang.org/
[actions-image]: https://github.com/krak3n/gofig/workflows/Test%20Library/badge.svg
[actions-url]: https://github.com/krak3n/gofig/actions?query=workflow%3ATest%20Library
[gofig-doc-image]: https://img.shields.io/badge/Documentation-pkg.go.dev-00ADD8.svg
[gofig-doc-url]: https://pkg.go.dev/go.krak3n.codes/gofig
[env-url]: ./parsers/env
[yaml-url]: ./parsers/yaml

# 🗃️ GoFig

[![Go 1.12+][goversion-image]][goversion-url]
[![Documentation][gofig-godoc-image]][gofig-godoc-url]

GoFig is a configuration loading library for Go. It aims to provide a simple and intuitive API that is
unopinionated for your configuration loading needs.

* **Status**: PoC (Proof of Concept)

## Example

``` go
package main

import (
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

	fig, err := gofig.New(&cfg)
	gofig.Must(err)

	parsers := []gofig.Parser{
		env.New(),
		gofig.FromFile(yaml.New(), "/path/to/config.yaml")
	}

	gofig.Must(fig.Parse(parsers...))

	fmt.Println(fmt.Sprintf("%+v", cfg))
}
```

## Parsers

GoFig implements it's parsers as sub modules. Currently it supports:

* [Environment Variables][env-godoc-url]
* [YAML][yaml-godoc-url]

### Planned

* TOML
* JSON

## Examples

* [Environment Variables](examples/001_env)
* [YAML](examples/002_yaml)
* [Multi Source](examples/003_multisource)

[goversion-image]: https://img.shields.io/badge/Go-1.13+-00ADD8.svg
[goversion-url]: https://golang.org/
[gofig-godoc-image]: https://img.shields.io/badge/godoc-reference-00ADD8.svg
[gofig-godoc-url]: https://godoc.org/go.krak3n.codes/gofig
[env-godoc-url]: https://godoc.org/go.krak3n.codes/gofig/parsers/env
[yaml-godoc-url]: https://godoc.org/go.krak3n.codes/gofig/parsers/yaml

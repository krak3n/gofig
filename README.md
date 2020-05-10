# 💡GoFig

[![Workflow Status][workflow-image]][workflow-url]
[![Go Version][goversion-image]][goversion-url]
[![Documentation][doc-image]][doc-url]

GoFig is a configuration loading library for Go. It aims to provide a simple, flexible and
decoupled API for all your configuration loading needs.

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
* [JSON][json-url]
* [TOML][toml-url]
* [YAML][yaml-url]

[workflow-image]: https://img.shields.io/github/workflow/status/krak3n/gofig/GoFig?style=flat&logo=github&logoColor=white&label=Workflow
[workflow-url]: https://github.com/krak3n/gofig/actions?query=workflow%3AGoFig
[goversion-image]: https://img.shields.io/badge/Go-1.13+-00ADD8.svg?style=flat&logo=go&logoColor=white
[goversion-url]: https://golang.org/
[doc-image]: https://img.shields.io/badge/Documentation-pkg.go.dev-00ADD8.svg?style=flat&logo=go&logoColor=white
[doc-url]: https://pkg.go.dev/go.krak3n.codes/gofig
[env-url]: ./parsers/env
[json-url]: ./parsers/json
[toml-url]: ./parsers/toml
[yaml-url]: ./parsers/yaml

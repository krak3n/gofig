# 💡GoFig

[![Go Version][goversion-image]][goversion-url]
[![Documentation][doc-image]][doc-url]
[![Workflow Status][workflow-image]][workflow-url]
[![Coverage][coverage-image]][coverage-url]
[![Go Report Card][report-image]][report-url]

GoFig is a configuration loading library for Go. It aims to provide a simple, flexible and
decoupled API for all your configuration loading needs.

* **Status**: PoC (Proof of Concept)

## Example

``` go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.krak3n.codes/gofig"
	"go.krak3n.codes/gofig/notifiers/fsnotify"
	"go.krak3n.codes/gofig/parsers/env"
	"go.krak3n.codes/gofig/parsers/yaml"
)

// Config is our configuration structure.
type Config struct {
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
	gfg, err := gofig.New(&cfg, gofig.WithDebug())
	gofig.Must(err)

	// Setsup a yaml parser with file notification support
	yml := gofig.FromFileAndNotify(yaml.New(), fsnotify.New("./config.yaml"))

	// Parse the yaml file and then environment variables
	gofig.Must(gfg.Parse(yml, env.New(env.WithPrefix("GOFIG"))))

	// Setup gofig notification channel to send notification of configuration updates
	notifyCh := make(chan error, 1)
	gfg.Notify(notifyCh, yml)

	// Setup OS signal notification
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// Wait for OS signal or configuration changes to reload your application
	for {
		log.Printf("configuration: %+v\n", cfg)

		select {
		case err := <-notifyCh:
			if err != nil {
				log.Println(err)
				return
			}
		case sig := <-signalCh:
			log.Println(sig)
			return
		}
	}

	// Check close errors.
	// You only need to call Close if using gofig.Notify or gofig.NotifyWithContext.
	gofig.Must(gfg.Close())
}
```

## Parsers

GoFig implements it's parsers as sub modules. Currently it supports:

* [Environment Variables][env-url]
* [JSON][json-url]
* [TOML][toml-url]
* [YAML][yaml-url]

## Priority

> Note priority enforcement can be disabled by using the `SetEnforcePriority()` option function.

By default parsers are parsed in the order in which they are given to `Parse`. This also applies to
`Notify`. Each field tracks which parser set it's values priority, this is then checked when each
parser attempts to set a value and is rejected if it's priority is not equal to or higher than the
priority set on the field.

For example:

``` go
cfg := struct{
	A string `gofig:"a"`
	B string `gofig:"b"`
}{}

gofig.New(&cfg)

p1 := gofig.NewInMemoryParser()
p1.Add("a", "Foo")
p1.Add("b", "Bar")

p2 := gofig.NewInMemoryParser()
p2.Add("a", "Fizz")

p3 := gofig.NewInMemoryParser()
p3.Add("b", "Buzz")

gofig.Must(gofig.Parse(p1, p2, p3))
```

When parsed the `cfg` struct will hold these values:

 * `A`: `Fizz` - Set by `p2` overriding `p1`
 * `B`: `Fizz` - Set by `p3` overriding `p1`

Since the `InMemoryParser` implements the `Notifier` interface we can also dynamically update values
via `Notify`. We will add `p2` to `Notify` and change the `b` value.

``` go
ch := make(chan error)
gofig.Notify(ch, p3, p2)

p2.Add("b", "Fizz")

gofig.Must(<-ch)
```

Here the values will be unchanged, even though `p2` updated `b`. This is because `p3` has higher
priority than `p2`.

New Parsers add will always have a higher priority than previously added parsers.

# Roadmap

* [x] (PoC) Support notification of config changes via `Notifier` interface
* [x] (PoC) Implement File notifier on changes to files via `fsnotify`
* [x] (Poc) Parser Order Priority on Notify events, e.g file changes should not override env var config
* [ ] Test Suite / Code Coverage reporting
* [ ] Helpful errors
* [ ] Support pointer values
* [ ] Default Values via a struct tag, e.g: `gofig:"foo,default=bar"`
* [ ] Support `omitempty` for pointer values which should not be initialised to their zero value.
* [ ] Add support for:
  * [ ] ETCD Parser / Notifier
  * [ ] Consul Parser / Notifier

[workflow-image]: https://img.shields.io/github/workflow/status/krak3n/gofig/GoFig?style=flat&logo=github&logoColor=white&label=Workflow
[workflow-url]: https://github.com/krak3n/gofig/actions?query=workflow%3AGoFig
[goversion-image]: https://img.shields.io/badge/Go-1.13+-00ADD8.svg?style=flat&logo=go&logoColor=white
[goversion-url]: https://golang.org/
[doc-image]: https://img.shields.io/badge/Documentation-pkg.go.dev-00ADD8.svg?style=flat&logo=go&logoColor=white
[doc-url]: https://pkg.go.dev/go.krak3n.codes/gofig
[report-image]: https://goreportcard.com/badge/github.com/krak3n/gofig?style=flat-square
[report-url]: https://goreportcard.com/report/github.com/krak3n/gofig
[coverage-image]: https://img.shields.io/codecov/c/gh/krak3n/gofig?label=Coverage&logo=codecov&logoColor=white
[coverage-url]: https://codecov.io/gh/krak3n/gofig
[env-url]: ./parsers/env
[json-url]: ./parsers/json
[toml-url]: ./parsers/toml
[yaml-url]: ./parsers/yaml

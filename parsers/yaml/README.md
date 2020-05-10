# YAML Parser

[![Workflow Status][workflow-image]][workflow-url]
[![Go Version][goversion-image]][goversion-url]
[![Example][playground-image]][playground-url]
[![Documentation][doc-image]][doc-url]

This parser loads `yaml` formatted configuration from an `io.ReadCloser`.

## Example

Click the Playground badge above to see the example running in the Go Playground.

``` go
package main

import (
	"fmt"
	"io/ioutil"

	"go.krak3n.codes/gofig"
	"go.krak3n.codes/gofig/parsers/yaml"
)

// Config is a struct to unpack configuration into.
type Config struct {
	Foo struct {
		Bar struct {
			Baz string `gofig:"baz"`
		} `gofig:"bar"`
	} `gofig:"foo"`
	Fizz struct {
		Buzz map[string]string `gofig:"buzz"`
	} `gofig:"fizz"`
	A struct {
		B map[string][]int `gofig:"b"`
	} `gofig:"a"`
	C struct {
		D map[string]map[string][]int `gofig:"d"`
	} `gofig:"c"`
}

const blob string = `
foo:
  bar:
    baz: bar
fizz:
  buzz:
    hello: world
    bill: ben`

func main() {
	var cfg Config

	// Initialise gofig with the struct config values will be placed into
	gfg, err := gofig.New(&cfg)
	gofig.Must(err)

	// Create a parser
	parser := yaml.New()

	// write some data to a config file
	path, err := create()
	gofig.Must(err)

	// Parse in order
	gofig.Must(gfg.Parse(
		gofig.FromFile(parser, path),
		gofig.FromString(parser, blob)))

	fmt.Println(fmt.Sprintf("%+v", cfg))
}

const contents = `
a:
  b:
    c: [1,2,3]
c:
  d:
    e:
      f: [1,2,3]`

func create() (string, error) {
	f, err := ioutil.TempFile("", "yaml")
	if err != nil {
		return "", err
	}

	if _, err := f.Write([]byte(contents)); err != nil {
		return "", err
	}

	return f.Name(), nil
}
```

[workflow-image]: https://img.shields.io/github/workflow/status/krak3n/gofig/YAML%20Parser?style=flat&logo=github&logoColor=white&label=Workflow
[workflow-url]: https://github.com/krak3n/gofig/actions?query=workflow%3A%22YAML+Parser%22
[goversion-image]: https://img.shields.io/badge/Go-1.13+-00ADD8.svg?style=flat&logo=go&logoColor=white
[goversion-url]: https://golang.org/
[playground-image]: https://img.shields.io/badge/Example-play.golang.org-00ADD8.svg?style=flat&logo=go&logoColor=white
[playground-url]: https://play.golang.org/p/hJLRH9pdhON
[doc-image]: https://img.shields.io/badge/Documentation-pkg.go.dev-00ADD8.svg?style=flat&logo=go&logoColor=white
[doc-url]: https://pkg.go.dev/go.krak3n.codes/gofig/parsers/yaml

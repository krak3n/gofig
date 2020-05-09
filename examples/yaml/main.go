package main

import (
	"fmt"

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

	// Create a parser
	yml := yaml.New()

	// Initialise gofig with the struct config values will be placed into
	gfg, err := gofig.New(&cfg)
	gofig.Must(err)

	file, err := gofig.FromFile(yml, "./config.yaml")
	gofig.Must(err)

	// Parse in order
	gofig.Must(gfg.Parse(file, gofig.FromString(yml, blob)))

	fmt.Println(fmt.Sprintf("%+v", cfg))
}

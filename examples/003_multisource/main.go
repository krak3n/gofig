package main

import (
	"fmt"
	"os"

	"go.krak3n.codes/gofig"
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
	// Initialise Config
	var cfg Config

	// Set environment variables
	os.Setenv("GOFIG_FOO_BAR", "bar")

	// Initialise gofig with the struct values will be parsed into
	gfg, err := gofig.New(&cfg)
	gofig.Must(err)

	// Parse the yaml file and then the envs
	gofig.Must(gfg.Parse(
		gofig.FromFile(yaml.New(), "./config.yaml"),
		env.New(env.HasAndTrimPrefix("GOFIG")),
	))

	fmt.Println("Foo.Bar:", cfg.Foo.Bar)     // bar
	fmt.Println("Fizz.Buzz:", cfg.Fizz.Buzz) // buzz
}

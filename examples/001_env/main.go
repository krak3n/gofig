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

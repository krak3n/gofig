package main

import (
	"fmt"
	"os"

	"go.krak3n.codes/gofig"
	"go.krak3n.codes/gofig/parsers/env"
)

type Config struct {
	Foo  string `gofig:"foo"`
	Bar  string `gofig:"bar"`
	Fizz struct {
		Buzz string `gofig:"buzz"`
	} `gofig:"fizz"`
}

func main() {
	var cfg Config

	os.Setenv("GOFIG_FOO", "foo")
	os.Setenv("GOFIG_BAR", "bar")
	os.Setenv("GOFIG_FIZZ_BUZZ", "buzz")

	// Initialise gofig with the struct config values will be placed into
	gfg, err := gofig.New(&cfg, gofig.WithNopLogger())
	gofig.Must(err)

	// Parse so environment variables
	gofig.Must(gfg.Parse(env.New(
		env.WithKeyPrefix("GOFIG_"),
		env.WithStrip("GOFIG_"),
	)))

	// Use the config
	fmt.Println("Foo:", cfg.Foo)             // foo
	fmt.Println("Bar:", cfg.Bar)             // bar
	fmt.Println("Fizz.Buzz:", cfg.Fizz.Buzz) // buzz
}

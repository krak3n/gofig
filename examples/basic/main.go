package main

import (
	"fmt"

	"go.krak3n.codes/gofig"
	"go.krak3n.codes/gofig/parsers/env"
)

type Config struct {
	Foo string `gofig:"foo"`
	Bar string `gofig:"bar"`
}

func main() {
	var cfg Config

	// Initialise gofig with the struct config values will be placed into
	gfg, err := gofig.New(&cfg)
	gofig.Must(err)

	// Parse so environment variables
	gofig.Must(gfg.Parse(env.New()))

	// Use the config
	fmt.Println("Foo:", cfg.Foo, "Bar:", cfg.Bar)
}

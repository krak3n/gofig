// Package gofig is a configuration loading library for Go. It aims to provide a simple and intuitive API that is
// unopinionated for your configuration loading needs.
//
// Example.
//
//   package main
//
//   import (
//   	"fmt"
//
//   	"go.krak3n.codes/gofig"
//   	"go.krak3n.codes/gofig/parsers/env"
//   )
//
//   type Config  struct {
//   	Foo struct {
//   		Bar string `gofig:"bar"`
//   	} `gofig:"foo"`
//   	Fizz struct {
//   		Buzz string `gofig:"buzz"`
//   	} `gofig:"fizz"`
//   }
//
//   func main() {
//   	var cfg Config
//
//   	// Initialise gofig with the destination struct
//   	gfg, err := gofig.New(&cfg)
//   	gofig.Must(err)
//
//   	// Parse the yaml file and then the envs
//   	gofig.Must(gfg.Parse(
//   		gofig.FromFile(yaml.New(), "./config.yaml"),
//   		env.New(env.HasAndTrimPrefix("GOFIG")),
//   	))
//
//   	fmt.Println(fmt.Sprintf("%+v", cfg))
//   }
package gofig

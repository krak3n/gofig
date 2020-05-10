module go.krak3n.codes/gofig/examples/003_multisource

go 1.13

replace go.krak3n.codes/gofig => ../..

replace go.krak3n.codes/gofig/parsers/env => ../../parsers/env

replace go.krak3n.codes/gofig/parsers/yaml => ../../parsers/yaml

require (
	go.krak3n.codes/gofig v0.0.0-00010101000000-000000000000
	go.krak3n.codes/gofig/parsers/env v0.0.0-20200510103519-0c8fcfa65524
	go.krak3n.codes/gofig/parsers/yaml v0.0.0-20200510102756-55cbeb2fd4a5
)

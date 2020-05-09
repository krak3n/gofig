package gofig

import (
	"strings"
)

const omitempty = "omitempty"

type tag struct {
	name      string
	omitempty bool
}

func parseTag(v string) tag {
	var t tag

	for i, v := range strings.Split(v, ",") {
		if i == 0 {
			t.name = v
			continue
		}

		if v == omitempty {
			t.omitempty = true
			continue
		}
	}

	return t
}

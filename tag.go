package gofig

import (
	"fmt"
	"strings"
)

const omitempty = "omitempty"

type tag struct {
	name      string
	omitempty bool
}

func (t tag) Stirng() string {
	return fmt.Sprintf("name=%s, omitempty=%t", t.name, t.omitempty)
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

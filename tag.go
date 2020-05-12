package gofig

import (
	"reflect"
	"strings"
)

const omitempty = "omitempty"

// Tag is a gofig struct tag.
type Tag struct {
	Name      string
	OmitEmpty bool
	RawTag    string
}

func (t Tag) String() string {
	return t.RawTag
}

// TagFromStructField returns a Tag from the struct fields tag.
func TagFromStructField(field reflect.StructField, tag string) Tag {
	t := Tag{
		Name: field.Name,
	}

	if v, ok := field.Tag.Lookup(DefaultStructTag); ok {
		t.RawTag = v

		for i, v := range strings.Split(v, ",") {
			if i == 0 {
				t.Name = v

				continue
			}

			if v == omitempty {
				t.OmitEmpty = true

				continue
			}
		}
	}

	return t
}

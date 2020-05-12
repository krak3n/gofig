package gofig

import (
	"strings"
)

// A Formatter formats a str.
type Formatter interface {
	Format(string) string
}

// FormatterFunc is an adapter function allowing regular methods to act as Formatter's.
type FormatterFunc func(string) string

// Format formats the str calling the wrapping fn.
func (fn FormatterFunc) Format(str string) string {
	return fn(str)
}

// LowerCaseFormatter returns a Formatter that formats str to lowercase.
func LowerCaseFormatter() Formatter {
	return FormatterFunc(strings.ToLower)
}

// KeyFormatter returns a Formatter that breaks a . delimited key path into it's constituent
// parts, formats each part with the given formatter returning the reconstituted . delimited key.
func KeyFormatter(fmtr Formatter) Formatter {
	return FormatterFunc(func(key string) string {
		elm := strings.Split(key, ".")

		for i, k := range elm {
			elm[i] = fmtr.Format(k)
		}

		return strings.Join(elm, ".")
	})
}

package gofig

import (
	"strings"
)

// A Formatter formats a str.
type Formatter interface {
	Format(key string, delimiter string) string
}

// FormatterFunc is an adapter function allowing regular methods to act as Formatter's.
type FormatterFunc func(string, string) string

// Format formats the str calling the wrapping fn.
func (fn FormatterFunc) Format(str string, delimiter string) string {
	return fn(str, delimiter)
}

// CaseSensitiveKeys returns a Formatter that maintains case sensitivity.
func CaseSensitiveKeys() Formatter {
	return FormatterFunc(func(key string, _ string) string {
		return key
	})
}

// CaseInsensitiveKeys returns a Formatter that formats keys to lowercase.
func CaseInsensitiveKeys() Formatter {
	return FormatterFunc(func(key string, delimiter string) string {
		elm := strings.Split(key, delimiter)

		for i, k := range elm {
			elm[i] = strings.ToLower(k)
		}

		return strings.Join(elm, delimiter)
	})
}

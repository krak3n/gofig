package gofig

import (
	"fmt"
	"log"
)

// A Logger can print log items.
type Logger interface {
	Print(values ...interface{})
	Printf(format string, values ...interface{})
}

// A LoggerFunc is an adapter function allowing regular methods to act as Loggers.
type LoggerFunc func(v ...interface{})

// Print calls the wrapped fn.
func (fn LoggerFunc) Print(v ...interface{}) {
	fn(v...)
}

// Printf calls the wrapped fn.
func (fn LoggerFunc) Printf(format string, v ...interface{}) {
	fn(fmt.Sprintf(format, v...))
}

// DefaultLogger returns a standard library logger.
func DefaultLogger() Logger {
	return LoggerFunc(func(v ...interface{}) {
		log.Println(v...)
	})
}

// NopLogger is a no operation that does nothing.
func NopLogger() Logger {
	return LoggerFunc(func(...interface{}) {})
}

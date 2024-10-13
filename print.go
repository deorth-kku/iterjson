package iterjson

import (
	"io"
	"os"
)

func IterFprint[K comparable, V any](f io.Writer, data any) (err error) {
	enc := NewEncoder[K, V](f)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	err = enc.Encode(data)
	return
}

func IterPrint[K comparable, V any](data any) error {
	return IterFprint[K, V](os.Stdout, data)
}

func Fprint(f io.Writer, data any) error {
	return IterFprint[string, map[string]any](f, data)
}

func Print(data any) error {
	return Fprint(os.Stdout, data)
}

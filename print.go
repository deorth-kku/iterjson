package iterjson

import (
	"fmt"
	"io"
	"os"
)

func IterFprint(f io.Writer, data any) (err error) {
	enc := NewEncoder(f)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	err = enc.Encode(data)
	return
}

func IterPrint(data any) error {
	return IterFprint(os.Stdout, data)
}

func Fprint(f io.Writer, data any) error {
	return IterFprint(f, data)
}

func Print(data any) error {
	return Fprint(os.Stdout, data)
}

func Println(data any) error {
	defer fmt.Print("\n")
	return Fprint(os.Stdout, data)
}

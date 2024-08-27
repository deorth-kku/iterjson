package iterjson

import (
	"encoding/json"
	"io"
	"iter"
	"maps"
	"slices"
)

type Encoder[K comparable, V any] struct {
	*json.Encoder
	w          *FormatWriter
	escapeHTML bool
}

func NewEncoder[K comparable, V any](w io.Writer) *Encoder[K, V] {
	fw := NewFormatWriter(w, "", "")
	return &Encoder[K, V]{json.NewEncoder(fw), fw, true}
}

func (e *Encoder[K, V]) encode(arg any) (err error) {
	switch v := arg.(type) {
	case iter.Seq[V]:
		return e.encodeSeq(v)
	case func(func(V) bool):
		return e.encodeSeq(v)
	case iter.Seq2[K, V]:
		return e.encodeSeq2(v)
	case func(func(K, V) bool):
		return e.encodeSeq2(v)
	case []V:
		return e.encodeSeq(slices.Values(v))
	case map[K]V:
		return e.encodeSeq2(maps.All(v))
	default:
		return e.Encoder.Encode(arg)
	}
}

func (e *Encoder[K, V]) Encode(arg any) (err error) {
	err = e.encode(arg)
	if err != nil {
		return
	}
	return
}

func (e *Encoder[K, V]) SetIndent(prefix, indent string) {
	if len(prefix) == 0 && len(indent) == 0 {
		return
	}
	e.w = NewFormatWriter(e.w.Writer, prefix, indent)
	e.Encoder = json.NewEncoder(e.w)
	e.Encoder.SetEscapeHTML(e.escapeHTML)
}

func (e *Encoder[K, V]) SetEscapeHTML(escapeHTML bool) {
	e.escapeHTML = escapeHTML
	e.Encoder.SetEscapeHTML(escapeHTML)
}

func (e *Encoder[K, V]) SetNewlines(newlines bool) {
	e.w.tailing_newline = newlines
}

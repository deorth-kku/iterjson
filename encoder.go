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
	w io.Writer
}

func NewEncoder[K comparable, V any](w io.Writer) *Encoder[K, V] {
	return &Encoder[K, V]{json.NewEncoder(NewFormatWriter(w, "", "")), w}
}

func (e *Encoder[K, V]) Encode(arg any) (err error) {
	switch v := arg.(type) {
	case iter.Seq[V]:
		return e.encodeSeq(v)
	case iter.Seq2[K, V]:
		return e.encodeSeq2(v)
	case []V:
		return e.encodeSeq(slices.Values(v))
	case map[K]V:
		return e.encodeSeq2(maps.All(v))
	default:
		return e.Encoder.Encode(arg)
	}
}

func (e *Encoder[K, V]) SetIndent(prefix, indent string) {
	if len(prefix) == 0 && len(indent) == 0 {
		return
	}
	if fw, ok := e.w.(*FormatWriter); ok {
		e.w = NewFormatWriter(fw.Writer, prefix, indent)
	} else {
		e.w = NewFormatWriter(e.w, prefix, indent)
	}
	e.Encoder = json.NewEncoder(e.w)
}

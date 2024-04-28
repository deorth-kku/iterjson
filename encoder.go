package iterjson

import (
	"encoding/json"
	"io"
	"iter"
)

type Encoder[K comparable, V any] struct {
	w io.Writer
}

func NewEncoder[K comparable, V any](w io.Writer) *Encoder[K, V] {
	return &Encoder[K, V]{w}
}

func (e *Encoder[K, V]) Encode(arg any) (err error) {
	switch v := arg.(type) {
	case iter.Seq[V]:
		return e.encodeSeq(v)
	case iter.Seq2[K, V]:
		return e.encodeSeq2(v)
	case []V:
		return e.encodeSeq(SliceSeq(v))
	case map[K]V:
		return e.encodeSeq2(MapSeq2(v))
	default:
		var data []byte
		data, err = json.Marshal(v)
		if err != nil {
			return
		}
		_, err = e.w.Write(data)
		if err != nil {
			return
		}
		return
	}
}

func (e *Encoder[K, V]) SetIndent(prefix, indent string) {
	if len(prefix) == 0 && len(indent) == 0 {
		return
	}
	e.w = NewFormatWriter(e.w, prefix, indent)
}

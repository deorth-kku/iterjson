package iterjson

import (
	"iter"
)

func MapSeq2[K comparable, V any](table map[K]V) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range table {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (e *Encoder[K, V]) encodeSeq2(iter iter.Seq2[K, V]) (err error) {
	_, err = e.w.Write([]byte("{"))
	if err != nil {
		return
	}
	first := true
	for k, v := range iter {
		if first {
			first = false
		} else {
			_, err = e.w.Write([]byte(","))
			if err != nil {
				return
			}
		}
		err = e.Encode(k)
		if err != nil {
			return
		}
		_, err = e.w.Write([]byte(":"))
		if err != nil {
			return
		}
		err = e.Encode(v)
		if err != nil {
			return
		}
	}
	_, err = e.w.Write([]byte("}"))
	return
}

package iterjson

import (
	"iter"
)

func (e *Encoder[K, V]) encodeSeq(iter iter.Seq[V]) (err error) {
	_, err = e.w.Write([]byte("["))
	if err != nil {
		return
	}
	first := true
	for line := range iter {
		if first {
			first = false
		} else {
			_, err = e.w.Write([]byte(","))
			if err != nil {
				return
			}
		}
		err = e.Encode(line)
		if err != nil {
			return
		}
	}
	_, err = e.w.Write([]byte("]"))
	return
}

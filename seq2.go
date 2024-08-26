package iterjson

import (
	"iter"
)

func (e *Encoder[K, V]) encodeSeq2(iter iter.Seq2[K, V]) (err error) {
	err = e.writeByte('{')
	if err != nil {
		return
	}
	first := true
	for k, v := range iter {
		if first {
			first = false
		} else {
			err = e.writeByte(',')
			if err != nil {
				return
			}
		}
		err = e.encode(k)
		if err != nil {
			return
		}
		err = e.writeByte(':')
		if err != nil {
			return
		}
		err = e.encode(v)
		if err != nil {
			return
		}
	}
	err = e.writeByte('}')
	return
}

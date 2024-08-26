package iterjson

import (
	"iter"
)

func (e *Encoder[K, V]) encodeSeq(iter iter.Seq[V]) (err error) {
	err = e.writeByte('[')
	if err != nil {
		return
	}
	first := true
	for line := range iter {
		if first {
			first = false
		} else {
			err = e.writeByte(',')
			if err != nil {
				return
			}
		}
		err = e.encode(line)
		if err != nil {
			return
		}
	}
	err = e.writeByte(']')
	return
}

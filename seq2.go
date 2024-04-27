package iterjson

import (
	"iter"
)

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

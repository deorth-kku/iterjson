package iterjson

import (
	"iter"
	"reflect"
)

func (e *Encoder) encodeSeq(iter iter.Seq[reflect.Value]) (err error) {
	err = e.w.WriteByte('[')
	if err != nil {
		return
	}
	first := true
	for line := range iter {
		if first {
			first = false
		} else {
			err = e.w.WriteByte(',')
			if err != nil {
				return
			}
		}
		err = e.encode(line)
		if err != nil {
			return
		}
	}
	return e.w.WriteByte(']')
}

package iterjson

import (
	"fmt"
	"iter"
	"reflect"
)

func (e *Encoder) encodeSeq2(iter iter.Seq2[reflect.Value, reflect.Value]) (err error) {
	err = e.w.WriteByte('{')
	if err != nil {
		return
	}
	first := true
	for k, v := range iter {
		if first {
			first = false
		} else {
			err = e.w.WriteByte(',')
			if err != nil {
				return
			}
		}
		err = e.enc.Encode(fmt.Sprint(k.Interface()))
		if err != nil {
			return
		}
		err = e.w.WriteByte(':')
		if err != nil {
			return
		}
		err = e.encode(v)
		if err != nil {
			return
		}
	}
	return e.w.WriteByte('}')
}

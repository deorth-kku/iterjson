package iterjson

import (
	"fmt"
	"iter"
	"reflect"
	"strings"
)

func rangeStruct(v reflect.Value) iter.Seq2[reflect.StructField, reflect.Value] {
	t := v.Type()
	return func(yield func(reflect.StructField, reflect.Value) bool) {
		for i := range t.NumField() {
			field := t.Field(i)
			if field.Anonymous {
				for f, v := range rangeStruct(v.Field(i)) {
					if !yield(f, v) {
						return
					}
				}
			} else {
				if !field.IsExported() {
					continue
				}
				if !yield(field, v.Field(i)) {
					return
				}
			}
		}
	}
}

func (e *Encoder) encodeStruct(v reflect.Value) (err error) {
	first := true
	for field, fv := range rangeStruct(v) {
		jsonTag := field.Tag.Get("json") // used as key
		var use_string bool
		if len(jsonTag) == 0 {
			jsonTag = field.Name
		} else if jsonTag == "-" {
			continue
		} else {
			var option string
			var ok bool
			jsonTag, option, ok = strings.Cut(jsonTag, ",")
			if ok {
				switch option {
				case "omitempty":
					if fv.IsZero() {
						continue
					}
				case "string":
					use_string = true
				}
			}
		}
		if first {
			first = false
			err = e.w.WriteByte('{')
		} else {
			err = e.w.WriteByte(',')
		}
		if err != nil {
			return err
		}
		err = e.enc.Encode(jsonTag)
		if err != nil {
			return err
		}
		err = e.w.WriteByte(':')
		if err != nil {
			return err
		}
		if use_string {
			err = e.enc.Encode(fmt.Sprint(fv.Interface()))
		} else {
			err = e.encode(fv)
		}
		if err != nil {
			return err
		}
	}
	if first {
		err = e.w.WriteByte('{')
		if err != nil {
			return
		}
	}
	return e.w.WriteByte('}')
}

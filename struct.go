package iterjson

import (
	"fmt"
	"reflect"
	"strings"
)

func (e *Encoder) encodeStruct(v reflect.Value) (err error) {
	first := true
	t := v.Type()
	for i := range t.NumField() {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		jsonTag := field.Tag.Get("json") // used as key
		fv := v.Field(i)                 // used as value
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

package iterjson

import (
	"reflect"
	"strconv"
	"strings"
)

func (e *Encoder) encodeStruct(v reflect.Value) error {
	err := e.w.WriteByte('{')
	if err != nil {
		return err
	}
	first := true
	t := v.Type()
	for i := range t.NumField() {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		jsonTag := field.Tag.Get("json") // used as key
		fv := v.Field(i)                 // used as value
		if len(jsonTag) == 0 {
			jsonTag = conv2snake_case(field.Name)
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
					fv = reflect.ValueOf(fv.String())
				}
			}
		}

		if first {
			first = false
		} else {
			err = e.w.WriteByte(',')
			if err != nil {
				return err
			}
		}
		e.w.write([]byte(strconv.Quote(jsonTag))...)
		err = e.w.WriteByte(':')
		if err != nil {
			return err
		}
		err = e.encode(fv)
		if err != nil {
			return err
		}
	}
	return e.w.WriteByte('}')
}

// Convert CamelCase to snake_case
func conv2snake_case(CamelCase string) string {
	result := make([]rune, len(CamelCase))
	for i, r := range CamelCase {
		if i > 0 && r >= 'A' && r <= 'Z' { // Capital letter found
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}

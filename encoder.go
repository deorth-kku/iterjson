package iterjson

import (
	"encoding/json"
	"io"
	"reflect"
)

type Encoder struct {
	*json.Encoder
	w          *FormatWriter
	escapeHTML bool
}

func NewEncoder(w io.Writer) *Encoder {
	fw := NewFormatWriter(w, "", "")
	return &Encoder{json.NewEncoder(fw), fw, true}
}

func (e *Encoder) encode(arg reflect.Value) error {
	for arg.Kind() == reflect.Pointer {
		arg = arg.Elem()
	}
	switch arg.Kind() {
	case reflect.Struct:
		if msl, ok := arg.Interface().(json.Marshaler); ok {
			data, err := msl.MarshalJSON()
			if err != nil {
				return err
			}
			e.w.write(data...)
			return nil
		}
		return e.encodeStruct(arg)
	case reflect.Map:
		return e.encodeSeq2(arg.Seq2())
	case reflect.Slice:
		return e.encodeSeq(arg.Seq())
	case reflect.Interface:
		return e.encode(reflect.ValueOf(arg.Interface()))
	case reflect.Func:
		ty := arg.Type()
		if ty.CanSeq() {
			return e.encodeSeq(arg.Seq())
		} else if ty.CanSeq2() {
			return e.encodeSeq2(arg.Seq2())
		}
		fallthrough
	default:
		return e.Encoder.Encode(arg.Interface())
	}
}

func (e *Encoder) Encode(arg any) error {
	return e.encode(reflect.ValueOf(arg))
}

func (e *Encoder) SetIndent(prefix, indent string) {
	if len(prefix) == 0 && len(indent) == 0 {
		return
	}
	e.w = NewFormatWriter(e.w.Writer, prefix, indent)
	e.Encoder = json.NewEncoder(e.w)
	e.Encoder.SetEscapeHTML(e.escapeHTML)
}

func (e *Encoder) SetEscapeHTML(escapeHTML bool) {
	e.escapeHTML = escapeHTML
	e.Encoder.SetEscapeHTML(escapeHTML)
}

func (e *Encoder) SetNewlines(newlines bool) {
	e.w.tailing_newline = newlines
}

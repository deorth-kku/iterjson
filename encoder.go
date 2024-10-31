package iterjson

import (
	"encoding/json"
	"io"
	"iter"
	"reflect"
)

type Encoder struct {
	enc        *json.Encoder
	w          *FormatWriter
	escapeHTML bool
}

func NewEncoder(w io.Writer) *Encoder {
	fw := NewFormatWriter(w, "", "")
	return &Encoder{json.NewEncoder(fw), fw, true}
}

func (e *Encoder) encode(arg reflect.Value) error {
	switch arg.Kind() {
	case reflect.Pointer:
		return e.encode(arg.Elem())
	case reflect.Interface:
		return e.encode(reflect.ValueOf(arg.Interface()))
	case reflect.Chan:
		return e.encodeSeq(arg.Seq())
	case reflect.Slice, reflect.Array:
		return e.encodeSeq(seq2Values(arg.Seq2()))
	case reflect.Map:
		return e.encodeSeq2(arg.Seq2())
	case reflect.Struct:
		if msl, ok := arg.Interface().(json.Marshaler); ok {
			return e.enc.Encode(msl)
		}
		return e.encodeStruct(arg)
	case reflect.Func:
		ty := arg.Type()
		if ty.CanSeq() {
			return e.encodeSeq(arg.Seq())
		} else if ty.CanSeq2() {
			return e.encodeSeq2(arg.Seq2())
		}
		fallthrough
	default:
		return e.enc.Encode(arg.Interface())
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
	e.enc = json.NewEncoder(e.w)
	e.enc.SetEscapeHTML(e.escapeHTML)
}

func (e *Encoder) SetEscapeHTML(escapeHTML bool) {
	e.escapeHTML = escapeHTML
	e.enc.SetEscapeHTML(escapeHTML)
}

func (e *Encoder) SetNewlines(newlines bool) {
	e.w.tailing_newline = newlines
}

func seq2Values[K any, V any](seq2 iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range seq2 {
			if !yield(v) {
				return
			}
		}
	}
}

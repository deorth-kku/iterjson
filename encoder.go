package iterjson

import (
	"encoding"
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

var (
	marshalerType     = reflect.TypeOf((*json.Marshaler)(nil)).Elem()
	textMarshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
)

func canMarshal(v reflect.Value) bool {
	t := v.Type()
	return t.Implements(marshalerType) || t.Implements(textMarshalerType)
}

func (e *Encoder) encode(arg reflect.Value) error {
	switch arg.Kind() {
	case reflect.Pointer:
		if arg.IsNil() {
			return e.enc.Encode(nil)
		}
		if !canMarshal(arg) {
			return e.encode(arg.Elem())
		}
	case reflect.Interface:
		return e.encode(reflect.ValueOf(arg.Interface()))
	case reflect.Invalid:
		return e.enc.Encode(nil)
	case reflect.Chan:
		if arg.IsNil() {
			return e.enc.Encode(nil)
		}
		if !canMarshal(arg) {
			return e.encodeSeq(arg.Seq())
		}
	case reflect.Slice:
		if arg.IsNil() {
			return e.enc.Encode(nil)
		}
		fallthrough
	case reflect.Array:
		if !canMarshal(arg) {
			return e.encodeSeq(iterSlice(arg))
		}
	case reflect.Map:
		if arg.IsNil() {
			return e.enc.Encode(nil)
		}
		if !canMarshal(arg) {
			return e.encodeSeq2(arg.Seq2())
		}
	case reflect.Struct:
		if !canMarshal(arg) {
			return e.encodeStruct(arg)
		}
	case reflect.Func:
		if arg.IsNil() {
			return e.enc.Encode(nil)
		}
		if !canMarshal(arg) {
			ty := arg.Type()
			if ty.CanSeq() {
				return e.encodeSeq(arg.Seq())
			} else if ty.CanSeq2() {
				return e.encodeSeq2(arg.Seq2())
			}
			return e.enc.Encode(arg.Interface())
		}
	}
	return e.enc.Encode(arg.Interface())
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

func iterSlice(arg reflect.Value) iter.Seq[reflect.Value] {
	return func(yield func(reflect.Value) bool) {
		for i := range arg.Len() {
			if !yield(arg.Index(i)) {
				return
			}
		}
	}
}

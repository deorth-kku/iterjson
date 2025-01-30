package iterjson

import (
	"reflect"
	"slices"
	"strconv"
	"strings"
	"sync"
)

type (
	getField   = func(stc reflect.Value) (reflect.Value, bool)
	writeValue = func(enc *Encoder, fv reflect.Value) error
)

type jsonOptions struct {
	tag        string
	use_string bool
	omitempty  bool
	omitzero   bool
	is_emb     bool
}

type run2 struct {
	jsonOptions
	getField
	writeValue
}

type canIsZero interface {
	IsZero() bool
}

var (
	canIsZeroType = reflect.TypeOf((*canIsZero)(nil)).Elem()
	cache         sync.Map
)

func getType(rt reflect.Type) []run2 {
	v, ok := cache.Load(rt)
	if ok {
		return v.([]run2)
	}
	rv := parseType(rt)
	cache.Store(rt, rv)
	return rv
}

func parseType(rt reflect.Type) []run2 {
	fs := make([]run2, 0, rt.NumField())
	for i := range cap(fs) {
		sf := rt.Field(i)
		if sf.Anonymous {
			sub := parseType(sf.Type)
			for j, run := range sub {
				sub[j].is_emb = true
				sub[j].getField = func(stc reflect.Value) (reflect.Value, bool) {
					return run.getField(stc.Field(i))
				}
			}
			fs = append(fs, sub...)
			continue
		}
		if !sf.IsExported() {
			continue
		}
		var opts jsonOptions
		opts.tag = sf.Tag.Get("json")
		switch opts.tag {
		case "":
			opts.tag = sf.Name
		case "-":
			continue
		default:
			opt0 := strings.Split(opts.tag, ",")
			opts.tag = opt0[0]
			for _, opt := range opt0[1:] {
				switch opt {
				case "omitempty":
					opts.omitempty = true
				case "omitzero":
					opts.omitzero = true
				case "string":
					opts.use_string = true
				}
			}
		}
		r := run2{
			jsonOptions: opts,
			getField:    func(stc reflect.Value) (reflect.Value, bool) { return stc.Field(i), false },
			writeValue: func(enc *Encoder, fv reflect.Value) error {
				return enc.encode(fv)
			},
		}
		switch {
		case opts.omitzero:
			if sf.Type.Implements(canIsZeroType) {
				r.getField = func(stc reflect.Value) (reflect.Value, bool) {
					fv := stc.Field(i)
					return fv, fv.Interface().(canIsZero).IsZero()
				}
			} else {
				r.getField = func(stc reflect.Value) (reflect.Value, bool) {
					fv := stc.Field(i)
					return fv, fv.IsZero()
				}
			}
		case opts.omitempty:
			switch sf.Type.Kind() {
			case reflect.String:
				r.getField = func(stc reflect.Value) (reflect.Value, bool) {
					fv := stc.Field(i)
					return fv, len(fv.String()) == 0
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				r.getField = func(stc reflect.Value) (reflect.Value, bool) {
					fv := stc.Field(i)
					return fv, fv.Int() == 0
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				r.getField = func(stc reflect.Value) (reflect.Value, bool) {
					fv := stc.Field(i)
					return fv, fv.Uint() == 0
				}
			case reflect.Float32, reflect.Float64:
				r.getField = func(stc reflect.Value) (reflect.Value, bool) {
					fv := stc.Field(i)
					return fv, fv.Float() == 0
				}
			case reflect.Bool:
				r.getField = func(stc reflect.Value) (reflect.Value, bool) {
					fv := stc.Field(i)
					return fv, !fv.Bool()
				}
			case reflect.Slice, reflect.Map, reflect.Chan, reflect.Pointer, reflect.UnsafePointer:
				r.getField = func(stc reflect.Value) (reflect.Value, bool) {
					fv := stc.Field(i)
					return fv, fv.IsNil()
				}
			}
		}
		if opts.use_string {
			switch sf.Type.Kind() {
			case reflect.String:
				r.writeValue = func(enc *Encoder, fv reflect.Value) error {
					return enc.writestring(fv.String())
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				r.writeValue = func(enc *Encoder, fv reflect.Value) error {
					return enc.writestring(strconv.FormatInt(fv.Int(), 10))
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				r.writeValue = func(enc *Encoder, fv reflect.Value) error {
					return enc.writestring(strconv.FormatUint(fv.Uint(), 10))
				}
			case reflect.Float32:
				r.writeValue = func(enc *Encoder, fv reflect.Value) error {
					return enc.writestring(strconv.FormatFloat(fv.Float(), 'f', -1, 32))
				}
			case reflect.Float64:
				r.writeValue = func(enc *Encoder, fv reflect.Value) error {
					return enc.writestring(strconv.FormatFloat(fv.Float(), 'f', -1, 64))
				}
			case reflect.Bool:
				r.writeValue = func(enc *Encoder, fv reflect.Value) error {
					return enc.writestring(strconv.FormatBool(fv.Bool()))
				}
			}
		}
		fs = append(fs, r)
	}
	tags := make(map[string]struct{}, len(fs))
	for _, opt := range fs {
		if opt.is_emb {
			continue
		}
		tags[opt.tag] = struct{}{}
	}
	return slices.DeleteFunc(fs, func(v run2) bool {
		if !v.is_emb {
			return false
		}
		if _, ok := tags[v.tag]; ok {
			return true
		}
		tags[v.tag] = struct{}{}
		return false
	})
}
func (e *Encoder) writestring(s string) (err error) {
	err = e.w.WriteByte('"')
	if err != nil {
		return err
	}
	_, err = e.w.Write([]byte(s))
	if err != nil {
		return err
	}
	err = e.w.WriteByte('"')
	if err != nil {
		return err
	}
	return
}

func (e *Encoder) encodeStruct(v reflect.Value) (err error) {
	first := true
	for _, f := range getType(v.Type()) {
		fv, skip := f.getField(v)
		if skip {
			continue
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
		err = e.writestring(f.tag)
		if err != nil {
			return err
		}
		err = e.w.WriteByte(':')
		if err != nil {
			return err
		}
		err = f.writeValue(e, fv)
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

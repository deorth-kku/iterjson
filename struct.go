package iterjson

import (
	"reflect"
	"slices"
	"strconv"
	"strings"
	"sync"
)

type (
	skipFunc   = func(stc reflect.Value) bool
	writeValue = func(enc *Encoder, stc reflect.Value) error
)

type jsonOptions struct {
	tag        string
	use_string bool
	omitempty  bool
	is_emb     bool
}

type run2 struct {
	jsonOptions
	skipFunc
	writeValue
}

var cache sync.Map

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
				sub[j].skipFunc = func(stc reflect.Value) bool {
					return run.skipFunc(stc.Field(i))
				}
				sub[j].writeValue = func(enc *Encoder, stc reflect.Value) error {
					return run.writeValue(enc, stc.Field(i))
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
			var option string
			var ok bool
			opts.tag, option, ok = strings.Cut(opts.tag, ",")
			if !ok {
				break
			}
			switch option {
			case "omitempty":
				opts.omitempty = true
			case "string":
				opts.use_string = true
			}
		}
		r := run2{
			jsonOptions: opts,
			skipFunc: func(stc reflect.Value) bool {
				if opts.omitempty && stc.Field(i).IsZero() {
					return true
				}
				return false
			},
		}
		if opts.use_string {
			switch sf.Type.Kind() {
			case reflect.String:
				r.writeValue = func(enc *Encoder, stc reflect.Value) error {
					return enc.writestring(stc.Field(i).String())
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				r.writeValue = func(enc *Encoder, stc reflect.Value) error {
					return enc.writestring(strconv.FormatInt(stc.Field(i).Int(), 10))
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				r.writeValue = func(enc *Encoder, stc reflect.Value) error {
					return enc.writestring(strconv.FormatUint(stc.Field(i).Uint(), 10))
				}
			case reflect.Float32:
				r.writeValue = func(enc *Encoder, stc reflect.Value) error {
					return enc.writestring(strconv.FormatFloat(stc.Field(i).Float(), 'f', -1, 32))
				}
			case reflect.Float64:
				r.writeValue = func(enc *Encoder, stc reflect.Value) error {
					return enc.writestring(strconv.FormatFloat(stc.Field(i).Float(), 'f', -1, 64))
				}
			case reflect.Bool:
				r.writeValue = func(enc *Encoder, stc reflect.Value) error {
					return enc.writestring(strconv.FormatBool(stc.Field(i).Bool()))
				}
			default:
				r.writeValue = func(enc *Encoder, stc reflect.Value) error {
					return enc.encode(stc.Field(i))
				}
			}
		} else {
			r.writeValue = func(enc *Encoder, stc reflect.Value) error {
				return enc.encode(stc.Field(i))
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
		if f.skipFunc(v) {
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
		err = f.writeValue(e, v)
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

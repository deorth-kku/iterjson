package iterjson

import (
	"io"
)

type FormatReader struct {
	io.Reader
	buf     []byte
	indent  []byte
	prefix  []byte
	level   int
	quoted  bool
	escaped bool
}

func NewFormatReader(reader io.Reader, prefix, indent string) (f *FormatReader) {
	f = new(FormatReader)
	f.indent = []byte(indent)
	f.prefix = []byte(prefix)
	f.Reader = reader
	return
}

func (f *FormatReader) Read(data []byte) (n int, err error) {
	t := make([]byte, 1)

	for {
		n, err = f.Reader.Read(t)
		if err != nil {
			break
		}
		if n == 0 {
			break
		}
		if f.escaped {
			f.buf = append(f.buf, t...)
			f.escaped = !f.escaped
			continue
		}
		if f.quoted {
			switch t[0] {
			case '\\':
				f.escaped = true
			case '"':
				f.quoted = !f.quoted
			}
			f.buf = append(f.buf, t...)
			continue
		}
		switch t[0] {
		case ' ', '\t', '\r', '\n': // skip whitespaces
		case '[', '{':
			f.buf = append(f.buf, t...)
			f.level++
			f.newline()
			f.writeindent()
		case ']', '}':
			f.newline()
			f.buf = append(f.buf, t...)
			f.level--
			f.newline()
			f.writeindent()
		case ':':
			f.buf = append(f.buf, t...)
			f.buf = append(f.buf, ' ')
		case ',':
			f.buf = append(f.buf, t...)
			f.newline()
			f.writeindent()
		case '"':
			f.quoted = !f.quoted
			fallthrough
		default:
			f.buf = append(f.buf, t...)
		}
		if len(f.buf) >= len(data) {
			break
		}
	}
	n = copy(data, f.buf)
	f.buf = f.buf[n:]
	return
}

func (f *FormatReader) writeindent() {
	for range f.level {
		f.buf = append(f.buf, f.indent...)
	}
}

func (f *FormatReader) newline() {
	f.buf = append(f.buf, '\n')
	f.buf = append(f.buf, f.prefix...)
}

package iterjson

import (
	"io"
	"slices"
)

type Formatter struct {
	buf             []byte
	indent          []byte
	prefix          []byte
	level           int
	has_element     bool
	quoted          bool
	escaped         bool
	compress        bool
	tailing_newline bool
}

func NewFormatter(prefix, indent string) *Formatter {
	return &Formatter{
		prefix:          []byte(prefix),
		indent:          []byte(indent),
		compress:        len(prefix) == 0 && len(indent) == 0,
		tailing_newline: false,
	}
}

func (f *Formatter) Format(p []byte) []byte {
	f.buf = slices.Grow(f.buf[:0], len(p))
	f.Write(p)
	return f.buf
}

func (f *Formatter) Write(p []byte) (int, error) {
	for _, b := range p {
		if f.escaped {
			f.write(b)
			f.escaped = false
			continue
		}
		if f.quoted {
			switch b {
			case '\\':
				f.escaped = true
			case '"':
				f.quoted = !f.quoted
			}
			f.write(b)
			continue
		}
		switch b {
		case ' ', '\t', '\r', '\n':
		case '[', '{':
			if !f.has_element && f.level != 0 {
				f.newline()
				f.writeindent()
			} else {
				f.has_element = false
			}
			f.write(b)
			f.level++
		case ']', '}':
			if f.has_element {
				f.newline()
				f.level--
				f.writeindent()
			} else {
				f.has_element = true
				f.level--
			}
			f.write(b)
			if f.level == 0 && f.tailing_newline {
				f.write('\n')
			}
		case ':':
			f.write(b)
			if !f.compress {
				f.write(' ')
			}
		case ',':
			f.write(b)
			f.newline()
			f.writeindent()
		case '"':
			f.quoted = !f.quoted
			fallthrough
		default:
			if !f.has_element {
				f.has_element = true
				if f.level != 0 {
					f.newline()
				}
				f.writeindent()
			}
			f.write(b)
		}
	}
	return len(p), nil
}

func (f *Formatter) write(b ...byte) {
	f.buf = append(f.buf, b...)
}

func (f *Formatter) Read(data []byte) (int, error) {
	n := copy(data, f.buf)
	f.buf = f.buf[n:]
	var err error
	if len(f.buf) == 0 {
		err = io.EOF
	}
	return int(n), err
}

func (f *Formatter) writeindent() {
	if f.compress || f.level < 0 {
		return
	}
	for range f.level {
		f.write(f.indent...)
	}
}

func (f *Formatter) newline() {
	if f.compress {
		return
	}
	f.write('\n')
	f.write(f.prefix...)
}

func Format(data []byte, prefix, indent string) []byte {
	return NewFormatter(prefix, indent).Format(data)
}

package iterjson

import "io"

type FormatWriter struct {
	io.Writer
	buf     []byte
	indent  []byte
	prefix  []byte
	level   int
	quoted  bool
	escaped bool
}

func (w *FormatWriter) Write(p []byte) (int, error) {
	var err error
	for _, b := range p {
		if w.escaped {
			w.write(b)
			w.escaped = false
			continue
		}
		if w.quoted {
			switch b {
			case '\\':
				w.escaped = true
			case '"':
				w.quoted = !w.quoted
			}
			w.write(b)
			continue
		}
		switch b {
		case ' ', '\t', '\r', '\n':
		case '[', '{':
			w.write(b)
			w.level++
			w.newline()
			w.writeindent()
		case ']', '}':
			w.newline()
			w.level--
			w.writeindent()
			w.write(b)
		case ':':
			w.write(b)
			w.write(' ')
		case ',':
			w.write(b)
			w.newline()
			w.writeindent()
		case '"':
			w.quoted = !w.quoted
			fallthrough
		default:
			w.write(b)
		}
	}
	n, err := w.Writer.Write(w.buf)
	if err != nil {
		return 0, err
	}
	w.buf = w.buf[n:]
	return len(p), nil
}

func (w *FormatWriter) write(b ...byte) {
	w.buf = append(w.buf, b...)
	return
}

func (w *FormatWriter) writeindent() {
	for range w.level {
		w.write(w.indent...)
	}
	return
}

func (w *FormatWriter) newline() {
	w.write('\n')
	w.write(w.prefix...)
	return
}

func NewFormatWriter(dst io.Writer, prefix, indent string) (w *FormatWriter) {
	w = new(FormatWriter)
	w.Writer = dst
	w.prefix = []byte(prefix)
	w.indent = []byte(indent)
	return
}

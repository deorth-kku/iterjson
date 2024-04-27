package iterjson

import "io"

type FormatWriter struct {
	io.Writer
	indent  []byte
	prefix  []byte
	level   int
	quoted  bool
	escaped bool
}

func (w *FormatWriter) Write(p []byte) (int, error) {
	var err error
	for i, b := range p {
		if w.escaped {
			err = w.write(b)
			if err != nil {
				return i, err
			}
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
			err = w.write(b)
			if err != nil {
				return i, err
			}
			continue
		}
		switch b {
		case ' ', '\t', '\r', '\n':
		case '[', '{':
			err = w.write(b)
			if err != nil {
				return i, err
			}
			w.level++
			err = w.newline()
			if err != nil {
				return i, err
			}
			err = w.writeindent()
			if err != nil {
				return i, err
			}
		case ']', '}':
			err = w.newline()
			if err != nil {
				return i, err
			}
			err = w.write(b)
			if err != nil {
				return i, err
			}
			w.level--
			err = w.newline()
			if err != nil {
				return i, err
			}
			err = w.writeindent()
			if err != nil {
				return i, err
			}
		case ':':
			err = w.write(b)
			if err != nil {
				return i, err
			}
			err = w.write(' ')
			if err != nil {
				return i, err
			}
		case ',':
			err = w.write(b)
			if err != nil {
				return i, err
			}
			err = w.newline()
			if err != nil {
				return i, err
			}
			err = w.writeindent()
			if err != nil {
				return i, err
			}
		case '"':
			w.quoted = !w.quoted
			fallthrough
		default:
			err = w.write(b)
			if err != nil {
				return i, err
			}
		}
	}
	return len(p), nil
}

func (w *FormatWriter) write(b byte) (err error) {
	_, err = w.Writer.Write([]byte{b})
	return
}

func (w *FormatWriter) writeindent() (err error) {
	for range w.level {
		_, err = w.Writer.Write(w.indent)
		if err != nil {
			return
		}
	}
	return
}

func (w *FormatWriter) newline() (err error) {
	err = w.write('\n')
	if err != nil {
		return
	}
	_, err = w.Write(w.prefix)
	return
}

func NewFormatWriter(dst io.Writer, prefix, indent string) (w *FormatWriter) {
	w = new(FormatWriter)
	w.Writer = dst
	w.prefix = []byte(prefix)
	w.indent = []byte(indent)
	return
}

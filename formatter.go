package iterjson

type Formatter struct {
	buf         []byte
	indent      []byte
	prefix      []byte
	level       int
	has_element bool
	quoted      bool
	escaped     bool
	compress    bool
}

func NewFormatter(prefix, indent string) *Formatter {
	return &Formatter{
		prefix:   []byte(prefix),
		indent:   []byte(indent),
		compress: len(prefix) == 0 && len(indent) == 0,
	}
}

func (f *Formatter) Format(p []byte) []byte {
	f.buf = nil
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
			if f.level == 0 {
				f.newline()
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
				f.newline()
				f.writeindent()
			}
			f.write(b)
		}
	}
	return f.buf
}

func (f *Formatter) write(b ...byte) {
	f.buf = append(f.buf, b...)
	return
}

func (f *Formatter) writeindent() {
	if f.compress || f.level < 0 {
		return
	}
	for range f.level {
		f.write(f.indent...)
	}
	return
}

func (f *Formatter) newline() {
	if f.compress {
		return
	}
	f.write('\n')
	f.write(f.prefix...)
	return
}

func Format(data []byte, prefix, indent string) []byte {
	return NewFormatter(prefix, indent).Format(data)
}

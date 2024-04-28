package iterjson

import (
	"io"
)

type FormatReader struct {
	io.Reader
	*Formatter
	buf []byte
}

func NewFormatReader(reader io.Reader, prefix, indent string) (f *FormatReader) {
	f = new(FormatReader)
	f.Reader = reader
	f.Formatter = NewFormatter(prefix, indent)
	return
}

func (f *FormatReader) Read(data []byte) (n int, err error) {
	datalen := len(data)
	for len(f.buf) < datalen {
		n, err = f.Reader.Read(data)
		if err == io.EOF {
			f.buf = append(f.buf, f.Format(data[:n])...)
			break
		} else if err != nil {
			return
		} else {
			f.buf = append(f.buf, f.Format(data[:n])...)
		}
	}
	n = copy(data, f.buf)
	clear(data[n:])
	f.buf = f.buf[n:]
	return
}

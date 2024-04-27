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
	temp := make([]byte, len(data))
	for len(f.buf) < datalen {
		n, err = f.Reader.Read(temp)
		if err == io.EOF {
			temp = temp[:n]
			f.buf = append(f.buf, f.Format(temp)...)
			break
		} else if err != nil {
			return
		} else {
			temp = temp[:n]
			f.buf = append(f.buf, f.Format(temp)...)
		}
	}
	n = copy(data, f.buf)
	f.buf = f.buf[n:]
	return
}

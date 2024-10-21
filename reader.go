package iterjson

import (
	"io"
)

type FormatReader struct {
	io.Reader
	*Formatter
}

func NewFormatReader(reader io.Reader, prefix, indent string) (f *FormatReader) {
	f = new(FormatReader)
	f.Reader = reader
	f.Formatter = NewFormatter(prefix, indent)
	return
}

func (f *FormatReader) Read(data []byte) (n int, err error) {
	datalen := len(data)
	for len(f.Formatter.buf) < datalen {
		n, err = f.Reader.Read(data)
		if err == io.EOF {
			f.Formatter.Write(data[:n])
			break
		} else if err != nil {
			return
		} else {
			f.Formatter.Write(data[:n])
		}
	}
	n, err = f.Formatter.Read(data)
	clear(data[n:])
	return
}

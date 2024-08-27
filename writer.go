package iterjson

import "io"

type FormatWriter struct {
	io.Writer
	*Formatter
}

func (w *FormatWriter) Write(p []byte) (int, error) {
	_, err := w.Writer.Write(w.Format(p))
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (w *FormatWriter) WriteByte(b byte) error {
	_, err := w.Write([]byte{b})
	return err
}

func NewFormatWriter(dst io.Writer, prefix, indent string) (w *FormatWriter) {
	w = new(FormatWriter)
	w.Writer = dst
	w.Formatter = NewFormatter(prefix, indent)
	return
}

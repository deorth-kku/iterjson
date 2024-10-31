package iterjson

import (
	"bytes"
	"io"
	"testing"
)

func BenchmarkReader(b *testing.B) {
	b.StopTimer()
	d0 := genDict(b.N)
	data := bytes.NewBuffer(nil)
	enc := NewEncoder(data)
	err := enc.Encode(d0)
	if err != nil {
		b.Error(err)
	}
	rd := NewFormatReader(data, "", "   ")
	b.StartTimer()
	_, err = io.Copy(bytes.NewBuffer(nil), rd)
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkWriter(b *testing.B) {
	b.StopTimer()
	d0 := genDict(b.N)
	data := bytes.NewBuffer(nil)
	enc := NewEncoder(data)
	err := enc.Encode(d0)
	if err != nil {
		b.Error(err)
	}
	wt := NewFormatWriter(bytes.NewBuffer(nil), "", "   ")
	b.StartTimer()
	_, err = io.Copy(wt, data)
	if err != nil {
		b.Error(err)
	}
}

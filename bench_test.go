package iterjson

import (
	"bytes"
	"encoding/json"
	"io"
	"maps"
	"math/rand/v2"
	"strconv"
	"testing"
)

type genDict struct {
	len int
}

func (g genDict) Range(yield func(string, int) bool) {
	for range g.len {
		if !yield(strconv.Itoa(rand.Int()), rand.Int()) {
			break
		}
	}
}

func BenchmarkEncoder(b *testing.B) {
	b.StopTimer()
	data := genDict{b.N}
	enc := NewEncoder[string, int](io.Discard)
	b.StartTimer()
	err := enc.Encode(data.Range)
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkStdlib(b *testing.B) {
	b.StopTimer()
	d0 := genDict{b.N}
	data := maps.Collect(d0.Range)
	enc := json.NewEncoder(io.Discard)
	b.StartTimer()
	err := enc.Encode(data)
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkReader(b *testing.B) {
	b.StopTimer()
	d0 := genDict{b.N}
	data := bytes.NewBuffer(nil)
	enc := NewEncoder[string, int](data)
	err := enc.Encode(d0.Range)
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
	d0 := genDict{b.N}
	data := bytes.NewBuffer(nil)
	enc := NewEncoder[string, int](data)
	err := enc.Encode(d0.Range)
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

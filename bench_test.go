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
	d1 := maps.Collect(d0.Range)
	d2, _ := json.Marshal(d1)
	data := bytes.NewBuffer(d2)
	rd := NewFormatReader(data, "", "   ")
	b.StartTimer()
	_, err := io.Copy(bytes.NewBuffer(nil), rd)
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkWriter(b *testing.B) {
	b.StopTimer()
	d0 := genDict{b.N}
	d1 := maps.Collect(d0.Range)
	d2, _ := json.Marshal(d1)
	data := bytes.NewBuffer(d2)
	wt := NewFormatWriter(bytes.NewBuffer(nil), "", "   ")
	b.StartTimer()
	_, err := io.Copy(wt, data)
	if err != nil {
		b.Error(err)
	}
}

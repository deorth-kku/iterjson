package iterjson

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"maps"
	"math/rand/v2"
	"net"
	"slices"
	"strconv"
	"testing"
)

func getSlice(n int) []string {
	out := make([]string, n)
	for i := range n {
		out[i] = strconv.Itoa(rand.Int())
	}
	return out
}

func BenchmarkSeq(b *testing.B) {
	b.StopTimer()
	data := getSlice(b.N)
	enc := NewEncoder(io.Discard)
	b.StartTimer()
	err := enc.Encode(slices.Values(data))
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkSlice(b *testing.B) {
	b.StopTimer()
	data := getSlice(b.N)
	enc := NewEncoder(io.Discard)
	b.StartTimer()
	err := enc.Encode(data)
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkSliceStd(b *testing.B) {
	b.StopTimer()
	data := getSlice(b.N)
	enc := json.NewEncoder(io.Discard)
	b.StartTimer()
	err := enc.Encode(data)
	if err != nil {
		b.Error(err)
	}
}

func genDict(n int) map[string]int {
	out := make(map[string]int, n)
	for range n {
		out[strconv.Itoa(rand.Int())] = rand.Int()
	}
	return out
}

func BenchmarkSeq2(b *testing.B) {
	b.StopTimer()
	data := genDict(b.N)
	enc := NewEncoder(io.Discard)
	b.StartTimer()
	err := enc.Encode(maps.All(data))
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkMap(b *testing.B) {
	b.StopTimer()
	data := genDict(b.N)
	enc := NewEncoder(io.Discard)
	b.StartTimer()
	err := enc.Encode(data)
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkMapStd(b *testing.B) {
	b.StopTimer()
	data := genDict(b.N)
	enc := json.NewEncoder(io.Discard)
	b.StartTimer()
	err := enc.Encode(data)
	if err != nil {
		b.Error(err)
	}
}

func randomIP() net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, rand.Uint32())
	return ip
}

func getStruct(n int) []testStruct {
	out := make([]testStruct, n)
	for i := range n {
		out[i].NormalField = strconv.Itoa(rand.Int())
		out[i].FieldToString = rand.Int()
		out[i].FieldWithTag = strconv.Itoa(rand.Int())
		out[i].FieldWithOmitEmpty = ""
		out[i].OmitField = false
		out[i].IP = randomIP()
	}
	return out
}

func BenchmarkStruct(b *testing.B) {
	b.StopTimer()
	arg := getStruct(b.N)
	enc := NewEncoder(io.Discard)
	b.StartTimer()
	err := enc.Encode(arg)
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkStructStd(b *testing.B) {
	b.StopTimer()
	arg := getStruct(b.N)
	enc := json.NewEncoder(io.Discard)
	b.StartTimer()
	err := enc.Encode(arg)
	if err != nil {
		b.Error(err)
	}
}

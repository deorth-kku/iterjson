package iterjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"os"
	"testing"
)

type testseq struct {
	list []string
}

func (it *testseq) Range() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, v := range it.list {
			if !yield(v) {
				return
			}
		}
	}
}

func TestSeq(t *testing.T) {
	l := &testseq{[]string{"a", "b", "c"}}
	data, err := Marshal[string, string](l.Range())
	if err != nil {
		t.Error(err)
		return
	}
	var list []string
	err = json.Unmarshal(data, &list)
	if err != nil {
		t.Error(err)
	}
}

type testseq2 struct {
	table map[string]string
}

func (it *testseq2) Range() iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		for k, v := range it.table {
			if !yield(k, v) {
				return
			}
		}
	}
}

func TestSeq2(t *testing.T) {
	l := &testseq2{map[string]string{
		"1": "a",
		"b": "2",
		"c": "3",
	}}
	data, err := Marshal[string, string](l.Range())
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(data))
	var table map[string]string
	err = json.Unmarshal(data, &table)
	if err != nil {
		t.Error(err)
	}
}

func TestFormatReader(t *testing.T) {
	l := map[string]string{
		"1\" ": "a",
		"b":    "2",
		"c":    "3",
	}
	data, err := json.Marshal(l)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(data))
	buf := bytes.NewReader(data)
	fm := NewFormatReader(buf, "", "    ")
	data, err = io.ReadAll(fm)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(data))

	var table map[string]string
	err = json.Unmarshal(data, &table)
	if err != nil {
		t.Error(err)
	}
}

func TestSetIndent(t *testing.T) {
	enc := NewEncoder[string, string](os.Stdout)
	enc.SetIndent("", "    ")
	l := &testseq2{map[string]string{
		"1\" ": "a",
		"b":    "2",
		"c":    "3",
	}}
	err := enc.Encode(l.Range())
	if err != nil {
		t.Error(err)
	}
}

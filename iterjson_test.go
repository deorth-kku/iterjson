package iterjson

import (
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"net/http"
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
	table map[string]any
}

func (it *testseq2) Range() iter.Seq2[string, any] {
	return func(yield func(string, any) bool) {
		for k, v := range it.table {
			if !yield(k, v) {
				return
			}
		}
	}
}

func TestSeq2(t *testing.T) {
	l := &testseq2{map[string]any{
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
	rsp, err := http.Get("https://api.github.com/repos/deorth-kku/iterjson")
	if err != nil {
		t.Error(err)
		return
	}
	defer rsp.Body.Close()
	r := NewFormatReader(rsp.Body, "", "")
	_, err = io.Copy(os.Stdout, r)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestSetIndent(t *testing.T) {
	enc := NewEncoder[string, any](os.Stdout)
	enc.SetIndent("", "    ")
	l := &testseq2{map[string]any{
		"a\" ": map[string]any{
			"x": "y",
		},
		"b": []string{},
		"c": []any{
			map[string]any{},
		},
	}}
	err := enc.Encode(l.Range())
	if err != nil {
		t.Error(err)
	}
}

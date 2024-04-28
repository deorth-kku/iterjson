package iterjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestSeq(t *testing.T) {
	l := SliceSeq([]string{"a", "b", "c"})
	data, err := Marshal[string, string](l)
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

func TestSeq2(t *testing.T) {
	l := map[string]any{
		"1": "a",
		"b": "2",
		"c": "3",
	}
	data, err := Marshal[string, any](MapSeq2(l))
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
	r := NewFormatReader(rsp.Body, "", "    ")
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, r)
	if err != nil {
		t.Error(err)
		return
	}
	r = NewFormatReader(buf, "", "")
	_, err = io.Copy(os.Stdout, r)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestSetIndent(t *testing.T) {
	enc := NewEncoder[string, any](os.Stdout)
	enc.SetIndent("", "    ")

	m := map[string]any{
		"a\" ": map[string]any{
			"x": "y",
		},
		"b": []string{},
		"c": []any{
			map[string]any{},
		},
	}
	err := enc.Encode(m)
	if err != nil {
		t.Error(err)
	}
}

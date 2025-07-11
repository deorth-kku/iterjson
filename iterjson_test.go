package iterjson

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"iter"
	"log/slog"
	"maps"
	"math/rand/v2"
	"net"
	"net/http"
	"os"
	"slices"
	"strconv"
	"testing"
)

func TestSeq(t *testing.T) {
	l := slices.Values([]string{"a", "b", "c"})
	data, err := Marshal(l)
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

func TestChannel(t *testing.T) {
	ch := make(chan int)
	go func() {
		for i := range 100 {
			ch <- i
		}
		close(ch)
	}()
	data, err := Marshal(ch)
	if err != nil {
		t.Error(err)
	}
	var list []int
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
	data, err := Marshal(maps.All(l))
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
	enc := NewEncoder(os.Stdout)
	enc.SetIndent("|test|", "    ")

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

func TestSetEscapeHTML(t *testing.T) {
	enc := NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "    ")

	m := map[string]any{
		"a\" ": map[string]any{
			"x": "y",
		},
		"b<>": []string{},
		"c": []any{
			map[string]any{},
		},
	}
	err := enc.Encode(m)
	if err != nil {
		t.Error(err)
	}
}

func TestSetNewlines(t *testing.T) {
	enc := NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetNewlines(true)

	m := map[string]any{
		"a\" ": map[string]any{
			"x": "y",
		},
		"b<>": []string{},
		"c": []any{
			map[string]any{},
		},
	}
	err := enc.Encode(m)
	if err != nil {
		t.Error(err)
	}
	enc.SetNewlines(false)
	err = enc.Encode(m)
	if err != nil {
		t.Error(err)
	}
}

type genDictAny struct {
	len int
}

func (g genDictAny) Range(yield func(string, any) bool) {
	for range g.len {
		if !yield(strconv.Itoa(rand.Int()), rand.Int()) {
			break
		}
	}
}

type genListAny struct {
	len int
}

func (g genListAny) Range(yield func(any) bool) {
	for range g.len {
		if !yield(rand.Int()) {
			break
		}
	}
}

func TestNested(t *testing.T) {
	d0 := genDictAny{10}
	d1 := genListAny{10}
	m := map[string]any{
		"dict": struct{ A any }{d0.Range},
		"list": d1.Range,
	}
	enc := NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	err := enc.Encode(m)
	if err != nil {
		t.Error(err)
	}
}

// verify that the result of Marshal is the same as [json.Marshal]
func verify(arg any) error {
	data0, err := Marshal(arg)
	if err != nil {
		return err
	}
	data1, err := json.Marshal(arg)
	if err != nil {
		return err
	}
	if !slices.Equal(data0, data1) {
		return errors.New("not equal")
	}
	return nil
}

func TestSlice(t *testing.T) {
	data := [4]string{
		"a",
		"b",
		"c",
		"d",
	}
	err := verify(data)
	if err != nil {
		t.Error(err)
	}
	err = verify(data[:])
	if err != nil {
		t.Error(err)
	}
}

func TestMap(t *testing.T) {
	err := verify(map[string]int{
		"a": 1, // only one key-value pair, otherwise because go's map is not ordered, there is no guarantee the result will be the same
	})
	if err != nil {
		t.Error(err)
	}
}

type embedStruct struct {
	F  int
	F2 uint
}

type testStruct struct {
	embedStruct
	F                  int
	NormalField        string
	FieldWithTag       string  `json:"this_is_tag"`
	FieldWithOmitEmpty string  `json:"om,omitempty"`
	OmitEmptyString    int     `json:"omstr,omitempty,string"`
	FieldToString      int     `json:"number,string"`
	FieldToStringFloat float32 `json:"float,string"`
	IP                 net.IP
	OmitField          bool `json:"-"`
	unexported         bool
}

func TestStruct(t *testing.T) {
	err := verify(testStruct{
		NormalField:        "a",
		FieldWithTag:       "b",
		FieldToString:      4,
		FieldToStringFloat: 1.24,
		OmitField:          true,
		unexported:         true,
		F:                  10,
		embedStruct: embedStruct{
			F:  1,
			F2: 2,
		},
		IP: net.IPv4(1, 3, 4, 5),
	})
	if err != nil {
		t.Error(err)
	}
}

func TestNil(t *testing.T) {
	err := verify(nil)
	if err != nil {
		t.Error(err)
	}

	var nilslice []string
	err = verify(nilslice)
	if err != nil {
		t.Error(err)
	}

	var nilmap map[string]any
	err = verify(nilmap)
	if err != nil {
		t.Error(err)
	}

	var nilptr *int
	err = verify(nilptr)
	if err != nil {
		t.Error(err)
	}

	var nilchan chan int
	_, err = Marshal(nilchan)
	if err != nil {
		t.Error(err)
	}

	var niliter iter.Seq[string]
	_, err = Marshal(niliter)
	if err != nil {
		t.Error(err)
	}
}

func TestEmpty(t *testing.T) {
	emptyslice := make([]string, 0)
	err := verify(emptyslice)
	if err != nil {
		t.Error(err)
	}
	emptymap := make(map[string]any)
	err = verify(emptymap)
	if err != nil {
		t.Error(err)
	}

	err = verify(struct{}{})
	if err != nil {
		t.Error(err)
	}
}

type example struct{}

func (example) MarshalJSON() ([]byte, error) {
	return []byte(`"example"`), nil
}

type example2 struct{}

func (*example2) MarshalJSON() ([]byte, error) {
	return []byte(`"example2"`), nil
}

type example3 func()

func (example3) MarshalText() ([]byte, error) {
	return []byte(`example3`), nil
}

func TestMarshaller(t *testing.T) {
	a := example{}
	err := verify(a)
	if err != nil {
		t.Error(err)
	}
	b := &example{}
	err = verify(b)
	if err != nil {
		t.Error(err)
	}

	c := example2{}
	err = verify(c)
	if err != nil {
		t.Error(err)
	}
	d := &example{}
	err = verify(d)
	if err != nil {
		t.Error(err)
	}

	e := example3(func() {})
	err = verify(e)
	if err != nil {
		t.Error(err)
	}

	f := &slog.LevelVar{}
	err = verify(f)
	if err != nil {
		t.Error(err)
	}
}

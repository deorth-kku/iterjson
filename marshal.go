package iterjson

import "bytes"

func Marshal[K comparable, V any](arg any) ([]byte, error) {
	return MarshalIndent[K, V](arg, "", "")
}

func MarshalIndent[K comparable, V any](arg any, prefix, indent string) (data []byte, err error) {
	w := bytes.NewBuffer(nil)
	e := NewEncoder[K, V](w)
	e.SetIndent(prefix, indent)
	err = e.Encode(arg)
	if err != nil {
		return
	}
	data = w.Bytes()
	return
}

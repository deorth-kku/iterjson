package iterjson

import "bytes"

func Marshal(arg any) ([]byte, error) {
	return MarshalIndent(arg, "", "")
}

func MarshalIndent(arg any, prefix, indent string) (data []byte, err error) {
	w := bytes.NewBuffer(nil)
	e := NewEncoder(w)
	e.SetIndent(prefix, indent)
	err = e.Encode(arg)
	if err != nil {
		return
	}
	data = w.Bytes()
	return
}

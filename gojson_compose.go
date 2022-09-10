package gojson

import (
	"bytes"
	"encoding/json"
)

func Compose(val *Value) (data []byte, err error) {
	w := bytes.NewBuffer([]byte{})

	encoder := json.NewEncoder(w)
	err = encoder.Encode(val.Interface())
	if err != nil {
		return
	}

	data = w.Bytes()
	return
}

package gojson

import (
	"bytes"
	"encoding/json"
)

func Parse(data []byte) (val *Value, err error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	var i interface{}
	err = decoder.Decode(&i)
	if err != nil {
		return
	}

	val, err = FromInterface(i)
	return
}

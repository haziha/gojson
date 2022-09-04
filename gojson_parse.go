package gojson

import (
	"bytes"
	"encoding/json"
	"fmt"
)

/*
现阶段先用官方json库解析, 再转成Value格式
*/

func Parse(data []byte) (val *Value, err error) {
	var j interface{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	err = decoder.Decode(&j)
	if err != nil {
		return
	}

	val, err = interface2value(j)

	return
}

func interface2value(j interface{}) (val *Value, err error) {
	switch v := j.(type) {
	case string:
		val, err = NewValue(String, j)
		return
	case json.Number:
		val, err = NewValue(Number, j)
		return
	case bool:
		val, err = NewValue(Boolean, j)
		return
	case nil:
		val, err = NewValue(Null, nil)
		return
	case []interface{}:
		array := make([]*Value, 0, len(v))
		var element *Value
		for i := range v {
			element, err = interface2value(v[i])
			if err != nil {
				return
			}
			array = append(array, element)
		}
		val, err = NewValue(Array, array)
		return
	case map[string]interface{}:
		object := make(map[string]*Value)
		var element *Value
		for k := range v {
			element, err = interface2value(v[k])
			if err != nil {
				return
			}
			object[k] = element
		}
		val, err = NewValue(Object, object)
		return
	default:
		err = fmt.Errorf("unknown type: %v", v)
		return
	}
}

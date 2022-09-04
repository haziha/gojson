package gojson

import (
	"encoding/json"
	"fmt"
)

/*
现阶段先转换成interface, 再用官方json库编码
*/

func Compose(val *Value) (data []byte, err error) {
	j, err := value2interface(val)
	if err != nil {
		return
	}
	data, err = json.Marshal(j)
	return
}

func value2interface(val *Value) (j interface{}, err error) {
	switch val.Type() {
	case String:
		j, err = val.String()
		return
	case Number:
		j, err = val.Number()
		return
	case Boolean:
		j, err = val.Boolean()
		return
	case Null:
		j, err = val.Null()
		return
	case Array:
		jSlice := make([]interface{}, 0, val.MustLen())
		for i := 0; i < val.MustLen(); i++ {
			var ele interface{}
			ele, err = value2interface(val.MustIndex(i))
			if err != nil {
				return
			}
			jSlice = append(jSlice, ele)
		}
		j = jSlice
		return
	case Object:
		jMap := make(map[string]interface{})
		for _, k := range val.MustKeys() {
			var v interface{}
			v, err = value2interface(val.MustValue(k))
			if err != nil {
				return
			}
			jMap[k] = v
		}
		j = jMap
		return
	default:
		err = fmt.Errorf("unknown type: %v", val.Type())
		return
	}
}

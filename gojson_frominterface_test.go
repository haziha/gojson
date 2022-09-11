package gojson

import (
	"fmt"
	"testing"
)

func TestFromInterface(t *testing.T) {
	{
		val, err := FromInterface(1)
		fmt.Println(val.typ, val.str, err)
	}
	{
		val, err := FromInterface(nil)
		fmt.Println(val.typ, err)
	}
	{
		val, err := FromInterface(true)
		fmt.Println(val.typ, val.boolean, err)
	}
	{
		val, err := FromInterface("string")
		fmt.Println(val.typ, val.str, err)
	}
	{
		val, err := FromInterface(map[string]interface{}{
			"number": 123,
			"string": "this is string",
			"bool":   true,
			"null":   nil,
		})
		fmt.Println(val.typ, val.obj, err)
	}
	{
		val, err := FromInterface([]interface{}{
			456, "slice string", true, nil,
		})
		fmt.Println(val.typ, val.arr, err)
	}
}

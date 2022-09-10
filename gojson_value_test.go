package gojson

import (
	"encoding/json"
	"fmt"
	"testing"
)

var demoString = `{
	"array": [
		"array.string",
		111111,
		true,
		null
	],
	"bool": true,
	"null": null,
	"number": 123456,
	"object": {
		"object.bool": true,
		"object.null": null,
		"object.number": 654321,
		"object.string": "string"
	},
	"string": "this is string"
}`

var demoInterface = map[string]interface{}{
	"string": "this is string",
	"number": 123456,
	"bool":   true,
	"null":   nil,
	"object": map[string]interface{}{
		"object.string": "string",
		"object.number": 654321,
		"object.bool":   true,
		"object.null":   nil,
	},
	"array": []interface{}{
		"array.string",
		111111,
		true,
		nil,
	},
}

func TestValue_Get(t *testing.T) {
	val, err := FromInterface(demoInterface)
	if err != nil {
		panic(err)
	}
	v, err := val.Get("object", "object.number")
	if err != nil {
		panic(err)
	}
	fmt.Println(v)
	v, err = val.Get("array", 1)
	if err != nil {
		panic(err)
	}
	fmt.Println(v)
}

func TestValue_Interface(t *testing.T) {
	val, err := FromInterface(demoInterface)
	fmt.Println(val, err)
	fmt.Println(val.Interface())
	data, _ := json.Marshal(val.Interface())
	fmt.Println(string(data))
}

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

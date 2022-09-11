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

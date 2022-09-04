package gojson

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	j := `{"string": "this is string", "bool": true, "number": 3.14, "null": null, "object": {
"obj.string": "obj.string", "obj.bool": false, "obj.number": 123.456, "obj.null": null
}, "array": [
"arr.string", true, false, 6457.111, null
]}`
	v, err := Parse([]byte(j))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(v.Len())
	fmt.Println(v.Get("object", "obj.string"))
	fmt.Println(v.Get("array", 1))
	fmt.Println(v.Get("object", "empty"))
	fmt.Println(v.Get("array", -1))
}

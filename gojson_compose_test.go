package gojson

import (
	"fmt"
	"testing"
)

func TestCompose(t *testing.T) {
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

	data, err := Compose(v)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(data))
}

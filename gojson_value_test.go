package gojson

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNewValue(t *testing.T) {
	fmt.Println(NewValue(Boolean, true))
	fmt.Println(NewValue(Boolean, false))
	fmt.Println(NewValue(Boolean, 1))
	fmt.Println(NewValue(Boolean, nil))

	fmt.Println(NewValue(Number, 1))
	fmt.Println(NewValue(Number, json.Number("1233")))
	fmt.Println(NewValue(Number, 3.14))
	fmt.Println(NewValue(Number, "3.14"))
	fmt.Println(NewValue(Number, "abc"))
	fmt.Println(NewValue(Number, nil))

	fmt.Println(NewValue(String, ""))
	fmt.Println(NewValue(String, []byte("123444")))
	var test []byte
	fmt.Println(NewValue(String, test))
	fmt.Println(NewValue(String, 123))
	fmt.Println(NewValue(String, true))
	fmt.Println(NewValue(String, nil))

	fmt.Println(NewValue(Null, nil))
	fmt.Println(NewValue(Null, 1))

	fmt.Println(NewValue(Array, []*Value{}))
	fmt.Println(NewValue(Array, nil))

	fmt.Println(NewValue(Object, map[string]*Value{}))
	fmt.Println(NewValue(Object, nil))
}

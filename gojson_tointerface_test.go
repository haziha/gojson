package gojson

import (
	"fmt"
	"testing"
)

func TestValue_ToInterface(t *testing.T) {
	v, _ := FromInterface(map[string]interface{}{
		"A": 1,
		"B": 2,
		"C": []interface{}{
			1, "2333", true, nil,
		},
	})
	var a struct {
		A interface{}
		B interface{}
		C interface{}
	}
	v.ToInterface(&a)
	fmt.Println(a)
}

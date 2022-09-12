package gojson

import (
	"fmt"
	"testing"
)

func TestValue_ToInterface3(t *testing.T) {
	v, _ := FromInterface(map[string]interface{}{
		"id":   1,
		"b_id": 2,
		"C-Id": []interface{}{
			1, "2333", true, nil,
		},
	})
	var a struct {
		Id  interface{}
		BId interface{}
		CId interface{}
	}
	v.ToInterface(&a)
	fmt.Println(a)
}

func TestValue_ToInterface2(t *testing.T) {
	v, _ := FromInterface(map[string]interface{}{
		"A": 1,
		"B": 2,
		"C": []interface{}{
			1, "2333", true, nil,
		},
	})
	var a struct {
		A interface{} `gojson:"/C/2"`
	}
	v.ToInterface(&a)
	fmt.Println(a)
}

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

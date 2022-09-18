package gojson

import (
	"fmt"
	"reflect"
	"testing"
)

func TestValue_ToInterface(t *testing.T) {
	newInt := func(n int) (p *int) {
		p = new(int)
		*p = n
		return
	}
	src := map[string]interface{}{
		"struct": struct {
			Struct         struct{}
			Map            map[string]interface{}
			MapPointer     map[string]*int
			ArrayPointer   [5]*int
			SliceInterface []interface{}
			SlicePointer   []*int
			Int            int64
			Uint           uint64
			Float64        float64
			Float32        float32
			Bool           bool
			String         string
			Nil            interface{}
		}{
			Struct: struct{}{},
			Map: map[string]interface{}{
				"struct": struct {
				}{},
				"map": map[string]interface{}{},
				"map pointer": map[string]*int{
					"11111": newInt(11111),
					"22222": newInt(22222),
				},
				"array pointer": [5]*int{
					newInt(-9), newInt(-99), newInt(-999), newInt(-999), newInt(-9999),
				},
				"slice interface": []interface{}{
					"this is string", true, 3.14, -999, nil, []*int{
						newInt(-1), newInt(-2), newInt(-3),
					},
				},
				"slice pointer": []*int{
					newInt(1), newInt(2), newInt(3), newInt(4),
				},
				"int":     int64(2333),
				"uint64":  uint64(99999999),
				"float64": 3.14159,
				"float32": 2.78,
				"bool":    true,
				"string":  "this is string",
				"null":    nil,
			},
			MapPointer: map[string]*int{
				"11111": newInt(11111),
				"22222": newInt(22222),
			},
			ArrayPointer: [5]*int{
				newInt(-9), newInt(-99), newInt(-999), newInt(-999), newInt(-9999),
			},
			SliceInterface: []interface{}{
				"this is string", true, 3.14, -999, nil, []*int{
					newInt(-1), newInt(-2), newInt(-3),
				},
			},
			SlicePointer: []*int{
				newInt(1), newInt(2), newInt(3), newInt(4),
			},
			Int:     int64(2333),
			Uint:    uint64(999999999),
			Float64: 3.14159,
			Float32: 2.78,
			Bool:    true,
			String:  "this is string",
			Nil:     nil,
		},
		"map": map[string]interface{}{
			"struct": struct {
			}{},
			"map": map[string]interface{}{},
			"map pointer": map[string]*int{
				"11111": newInt(11111),
				"22222": newInt(22222),
			},
			"array pointer": [5]*int{
				newInt(-9), newInt(-99), newInt(-999), newInt(-999), newInt(-9999),
			},
			"slice interface": []interface{}{
				"this is string", true, 3.14, -999, nil, []*int{
					newInt(-1), newInt(-2), newInt(-3),
				},
			},
			"slice pointer": []*int{
				newInt(1), newInt(2), newInt(3), newInt(4),
			},
			"int":     int64(2333),
			"uint64":  uint64(99999999),
			"float64": 3.14159,
			"float32": 2.78,
			"bool":    true,
			"string":  "this is string",
			"null":    nil,
		},
		"map pointer": map[string]*int{
			"11111": newInt(11111),
			"22222": newInt(22222),
		},
		"array pointer": [5]*int{
			newInt(-9), newInt(-99), newInt(-999), newInt(-999), newInt(-9999),
		},
		"slice interface": []interface{}{
			"this is string", true, 3.14, -999, nil, []*int{
				newInt(-1), newInt(-2), newInt(-3),
			},
		},
		"slice pointer": []*int{
			newInt(1), newInt(2), newInt(3), newInt(4),
		},
		"int":     int64(2333),
		"uint64":  uint64(99999999),
		"float64": 3.14159,
		"float32": 2.78,
		"bool":    true,
		"string":  "this is string",
		"null":    nil,
	}
	j, err := FromInterface(src)
	if err != nil {
		panic(err)
	}
	fmt.Println(j)
	dst := struct {
		Struct struct {
			Struct         struct{}
			Map            map[string]interface{}
			MapPointer     map[string]*int
			ArrayPointer   [5]*int
			SliceInterface []interface{}
			SlicePointer   []*int
			Int            int64
			Uint           uint64
			Float64        float64
			Float32        float32
			Bool           bool
			String         string
			Nil            interface{}
		}
		Map            map[string]interface{}
		MapPointer     map[string]*int `gk:"map pointer"`
		ArrayPointer   [5]*int         `gk:"array pointer"`
		SliceInterface []interface{}   `gk:"slice interface"`
		SlicePointer   []*int          `gk:"slice pointer"`
		Int            int64
		Uint           uint64
		Float64        float64
		Float32        float32
		Bool           bool
		String         string
		Nil            interface{} `gk:"null"`
	}{
		Map: map[string]interface{}{"other": "other"},
		Nil: "this is not nil",
	}
	err = j.ToInterface(&dst)
	if err != nil {
		panic(err)
	}
	fmt.Println(dst)
}

func TestValue_ToInterface_Struct(t *testing.T) {
	var src = `{
	"a": "b", "c": "d", "e": "f", "g": "h"
}`
	s := struct {
		A interface{}
		B interface{} `gk:"c"`
		C interface{} `gp:"/e" gk:"g"`
		D interface{} `gp:"/z" gk:"y"`
		E interface{} `gp:"/x"`
	}{}

	j, err := Parse([]byte(src))
	if err != nil {
		panic(err)
	}
	err = j.ToInterface(&s)
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
}

func TestValue_ToInterface_Interface(t *testing.T) {
	var src = "this is string"
	j, err := FromInterface(src)
	if err != nil {
		panic(err)
	}
	var a interface{}
	err = j.ToInterface(&a)
	if err != nil {
		panic(err)
	}
	fmt.Println(a)

	a = nil
	var b interface{} = &a
	err = j.ToInterface(&b)
	if err != nil {
		panic(err)
	}
	fmt.Println(a, b)
}

func TestValue_ToInterface_Pointer(t *testing.T) {
	src := 2333
	j, err := FromInterface(src)
	if err != nil {
		panic(err)
	}

	a := new(*int)

	err = j.ToInterface(&a)
	if err != nil {
		panic(err)
	}
	fmt.Println(a, *a)
}

func TestValue_ToInterface_Object(t *testing.T) {
	src := map[string]string{
		"a": "b", "c": "d", "e": "f",
	}
	j, err := FromInterface(src)
	if err != nil {
		panic(err)
	}

	a := make(map[string]string)
	err = j.ToInterface(&a)
	if err != nil {
		panic(err)
	}
	fmt.Println(a)

	a = make(map[string]string)
	b := &a
	err = j.ToInterface(b)
	if err != nil {
		panic(err)
	}
	fmt.Println(*b)
}

func TestValue_ToInterface_Array(t *testing.T) {
	src := []string{
		"a", "b", "c", "d",
	}
	a := new([]string)
	j, err := FromInterface(src)
	if err != nil {
		panic(err)
	}
	err = j.ToInterface(a)
	if err != nil {
		panic(err)
	}
	fmt.Println("a := new([]string)", a)

	b := make([]string, 0)
	b = append(b, "e", "f", "g", "h")
	err = j.ToInterface(&b)
	if err != nil {
		panic(err)
	}
	fmt.Println("b := make([]string, 0)", b)
	if !reflect.DeepEqual(src, *a) || !reflect.DeepEqual(src, b) {
		panic("unequal")
	}
}

func TestValue_ToInterface_BaseType(t *testing.T) {
	a := new(int)
	j, err := FromInterface(2333)
	if err != nil {
		panic(err)
	}
	err = j.ToInterface(a)
	if err != nil {
		panic(err)
	}
	fmt.Println("a := new(int):", *a)

	var b bool
	j, err = FromInterface(true)
	if err != nil {
		panic(err)
	}
	err = j.ToInterface(&b)
	if err != nil {
		panic(err)
	}
	fmt.Println("var b bool:", b)
}

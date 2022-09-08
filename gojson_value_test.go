package gojson

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestValue_Interface(t *testing.T) {
	val, _ := NewValue(Boolean, true)
	fmt.Println(val.Interface())
	_ = val.SetValue(String, "2333")
	fmt.Println(val.Interface())
	_ = val.SetValue(Number, 23333)
	fmt.Println(val.Interface())
	_ = val.SetValue(Null, nil)
	fmt.Println(val.Interface())
	_ = val.SetValue(Object, map[string]*Value{
		"test": MustNewValue(String, "value"), "test2": MustNewValue(Number, 1),
		"test3": MustNewValue(Boolean, true), "test4": MustNewValue(Null, nil),
		"test5": MustNewValue(Array, []*Value{
			MustNewValue(String, "str"), MustNewValue(Boolean, false), MustNewValue(Number, 233333),
			MustNewValue(Null, nil),
		}),
	})
	fmt.Println(val.Interface())
}

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

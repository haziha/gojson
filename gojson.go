package gojson

import (
	"encoding/json"
	"reflect"
)

var stringType reflect.Type
var sliceType reflect.Type
var mapType reflect.Type
var interfaceType reflect.Type
var jsonNumberType reflect.Type
var byteSliceType reflect.Type
var mapKind2Type map[reflect.Kind]Type

func init() {
	stringType = reflect.TypeOf((*string)(nil)).Elem()
	sliceType = reflect.TypeOf((*[]interface{})(nil)).Elem()
	mapType = reflect.TypeOf((*map[string]interface{})(nil)).Elem()
	interfaceType = reflect.TypeOf((*interface{})(nil)).Elem()
	jsonNumberType = reflect.TypeOf((*json.Number)(nil)).Elem()
	byteSliceType = reflect.TypeOf((*[]byte)(nil)).Elem()

	mapKind2Type = map[reflect.Kind]Type{
		reflect.Invalid:       Null,
		reflect.Bool:          Boolean,
		reflect.Int:           Number,
		reflect.Int8:          Number,
		reflect.Int16:         Number,
		reflect.Int32:         Number,
		reflect.Int64:         Number,
		reflect.Uint:          Number,
		reflect.Uint8:         Number,
		reflect.Uint16:        Number,
		reflect.Uint32:        Number,
		reflect.Uint64:        Number,
		reflect.Uintptr:       Unknown,
		reflect.Float32:       Number,
		reflect.Float64:       Number,
		reflect.Complex64:     Unknown,
		reflect.Complex128:    Unknown,
		reflect.Array:         Array,
		reflect.Chan:          Unknown,
		reflect.Func:          Unknown,
		reflect.Interface:     Unknown,
		reflect.Map:           Object,
		reflect.Pointer:       Unknown,
		reflect.Slice:         Array,
		reflect.String:        String,
		reflect.Struct:        Object,
		reflect.UnsafePointer: Unknown,
	}
}

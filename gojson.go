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

func init() {
	stringType = reflect.TypeOf((*string)(nil)).Elem()
	sliceType = reflect.TypeOf((*[]interface{})(nil)).Elem()
	mapType = reflect.TypeOf((*map[string]interface{})(nil)).Elem()
	interfaceType = reflect.TypeOf((*interface{})(nil)).Elem()
	jsonNumberType = reflect.TypeOf((*json.Number)(nil)).Elem()
	byteSliceType = reflect.TypeOf((*[]byte)(nil)).Elem()
}

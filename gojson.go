package gojson

import "reflect"

var byteType reflect.Type
var sliceType reflect.Type
var objectType reflect.Type

func init() {
	byteType = reflect.TypeOf([]byte{})
	sliceType = reflect.TypeOf([]*Value{})
	objectType = reflect.TypeOf(map[string]*Value{})
}

package gojson

import "fmt"

type Type uint

const (
	Unknown Type = 0
	String  Type = 1 << iota
	Number
	Boolean
	Null
	Array
	Object
)

func (_this Type) String() string {
	switch _this {
	case String:
		return "string"
	case Number:
		return "number"
	case Boolean:
		return "boolean"
	case Null:
		return "null"
	case Array:
		return "array"
	case Object:
		return "object"
	case String | Number:
		return "string or number"
	case Array | Object:
		return "array or object"
	default:
		return "unknown"
	}
}

func (_this Type) Error(mustType Type) error {
	return fmt.Errorf("type must be %s, but %s", mustType.String(), _this.String())
}

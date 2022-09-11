package gojson

import (
	"fmt"
	"github.com/haziha/golist"
	"reflect"
	"strconv"
)

func FromInterface(val interface{}) (v *Value, err error) {
	type tempPair struct {
		src reflect.Value
		dst *Value
	}

	v = new(Value)
	vList := golist.New[tempPair]()
	vList.PushBack(tempPair{src: reflect.ValueOf(val), dst: v})

	for vList.Len() != 0 {
		back := vList.Back()
		vList.Remove(back)

		pair := back.Value
		src := pair.src
		dst := pair.dst

		switch src.Kind() {
		case reflect.Invalid:
			dst.typ = Null
		case reflect.Bool:
			dst.typ = Boolean
			dst.boolean = src.Bool()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			dst.typ = Number
			dst.str = strconv.FormatInt(src.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			dst.typ = Number
			dst.str = strconv.FormatUint(src.Uint(), 10)
		case reflect.Float32, reflect.Float64:
			dst.typ = Number
			dst.str = strconv.FormatFloat(src.Float(), 'f', -1, 64)
		case reflect.Array:
			fallthrough
		case reflect.Slice:
			dst.typ = Array
			dst.arr = make([]*Value, src.Len())
			for i := 0; i < src.Len(); i++ {
				dst.arr[i] = new(Value)
				vList.PushBack(tempPair{src.Index(i), dst.arr[i]})
			}
		case reflect.Interface:
			fallthrough
		case reflect.Pointer:
			vList.PushBack(tempPair{src.Elem(), dst})
		case reflect.Map:
			dst.typ = Object
			if src.Type().Key() != stringType {
				err = fmt.Errorf("map key must be string, but %v", src.Type().Key())
				return nil, err
			}
			dst.obj = make(map[string]*Value)
			iter := src.MapRange()
			for iter.Next() {
				key := iter.Key()
				value := iter.Value()

				nV := new(Value)
				dst.obj[key.String()] = nV
				vList.PushBack(tempPair{value, nV})
			}
		case reflect.String:
			if src.Type() == jsonNumberType {
				dst.typ = Number
			} else {
				dst.typ = String
			}
			dst.str = src.String()
		case reflect.Struct:
			dst.typ = Object
			dst.obj = make(map[string]*Value)

			for i := 0; i < src.NumField(); i++ {
				if !src.Type().Field(i).IsExported() { // 过滤私有成员
					continue
				}
				nV := new(Value)
				dst.obj[src.Type().Field(i).Name] = nV
				vList.PushBack(tempPair{src.Field(i), nV})
			}
		default:
			v = nil
			err = fmt.Errorf("cannot convert %v to json type", src.Kind())
			return
		}
	}
	return
}

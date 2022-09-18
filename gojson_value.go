package gojson

import (
	"encoding/json"
	"fmt"
	"github.com/haziha/golist"
	"reflect"
	"strconv"
)

type Value struct {
	typ     Type
	str     string
	boolean bool
	arr     []*Value
	obj     map[string]*Value
}

func (_this *Value) Get(k ...interface{}) (val *Value, err error) {
	defer func() {
		err1 := recover()
		if err1 != nil {
			err = fmt.Errorf("%v", err1)
		}
	}()

	kVal := reflect.ValueOf(k)
	ptrList := golist.New[*Value]()
	ptrList.PushBack(_this)

	for i := range k {
		ptrBack := ptrList.Back().Value
		switch ptrBack.Type() {
		case Object:
			key := kVal.Index(i)
			if !key.IsValid() {
				err = fmt.Errorf("key must be string, but invalid")
				return
			}
			key = key.Elem()
			if key.Kind() == reflect.String {
				ptrList.PushBack(ptrBack.MustValue(key.String()))
			} else if key.Kind() == reflect.Slice {
				if key.IsZero() {
					err = fmt.Errorf("key must be string, but empty slice")
					return
				} else if key.Type() != byteSliceType {
					err = fmt.Errorf("key must be string, but not byte slice")
					return
				}
				ptrList.PushBack(ptrBack.MustValue(string(key.Bytes())))
			} else {
				err = fmt.Errorf("key must be string, but %v", key.Kind())
				return
			}
		case Array:
			index := kVal.Index(i)
			if !index.IsValid() {
				err = fmt.Errorf("index must be integer, but mepty")
				return
			}
			index = index.Elem()
			if index.CanInt() {
				ptrList.PushBack(ptrBack.MustIndex(int(index.Int())))
			} else if index.CanUint() {
				ptrList.PushBack(ptrBack.MustIndex(int(index.Uint())))
			} else if index.CanFloat() {
				ptrList.PushBack(ptrBack.MustIndex(int(index.Float())))
			} else if index.Kind() == reflect.String {
				jN := json.Number(index.String())
				var i64 int64
				if i64, err = jN.Int64(); err != nil {
					err = fmt.Errorf("index must be integer, but string and cannot conver to integer: %s", jN)
					return
				}
				ptrList.PushBack(ptrBack.MustIndex(int(i64)))
			} else {
				err = fmt.Errorf("index must be integer, but %v", index.Kind())
				return
			}
		default:
			err = fmt.Errorf("cannot get element in %v", ptrBack.Type())
			return
		}
	}

	return ptrList.Back().Value, nil
}

func (_this *Value) MustUint64() (v uint64) {
	v, err := _this.Uint64()
	if err != nil {
		panic(err)
	}
	return
}

func (_this *Value) MustInt64() (v int64) {
	v, err := _this.Int64()
	if err != nil {
		panic(err)
	}
	return
}

func (_this *Value) MustFloat64() (v float64) {
	v, err := _this.Float64()
	if err != nil {
		panic(err)
	}
	return
}

func (_this *Value) MustKeys() (v map[string]struct{}) {
	v, err := _this.Keys()
	if err != nil {
		panic(err)
	}
	return
}

func (_this *Value) MustValue(k string) (v *Value) {
	v, err := _this.Value(k)
	if err != nil {
		panic(err)
	}
	return
}

func (_this *Value) MustLen() (v int) {
	v, err := _this.Len()
	if err != nil {
		panic(err)
	}
	return
}

func (_this *Value) MustIndex(i int) (v *Value) {
	v, err := _this.Index(i)
	if err != nil {
		panic(err)
	}
	return
}

func (_this *Value) MustStr() (v string) {
	v, err := _this.Str()
	if err != nil {
		panic(err)
	}
	return
}

func (_this *Value) MustNumber() (v json.Number) {
	v, err := _this.Number()
	if err != nil {
		panic(err)
	}
	return
}

func (_this *Value) MustBool() (v bool) {
	v, err := _this.Bool()
	if err != nil {
		panic(err)
	}
	return
}

func (_this *Value) MustNull() (v interface{}) {
	v, err := _this.Null()
	if err != nil {
		panic(err)
	}
	return
}

func (_this *Value) Uint64() (v uint64, err error) {
	v, err = strconv.ParseUint(_this.str, 10, 64)
	return
}

func (_this *Value) Int64() (v int64, err error) {
	v, err = strconv.ParseInt(_this.str, 10, 64)
	return
}

func (_this *Value) Float64() (v float64, err error) {
	v, err = strconv.ParseFloat(_this.str, 64)
	return
}

func (_this *Value) Keys() (v map[string]struct{}, err error) {
	if _this.typ&Object == 0 {
		err = _this.typ.Error(Object)
	} else {
		v = make(map[string]struct{})

		for k := range _this.obj {
			v[k] = struct{}{}
		}
	}

	return
}

func (_this *Value) Value(key string) (v *Value, err error) {
	if _this.typ&Object == 0 {
		err = _this.typ.Error(Object)
	} else {
		v = _this.obj[key]
	}
	return
}

/*
Len
Just object or array
*/
func (_this *Value) Len() (v int, err error) {
	if _this.typ&(Object|Array) == 0 {
		err = _this.typ.Error(Object | Array)
	} else if _this.typ == Object {
		v = len(_this.obj)
	} else {
		v = len(_this.arr)
	}
	return
}

func (_this *Value) Index(i int) (v *Value, err error) {
	if _this.typ&Array == 0 {
		err = _this.typ.Error(Array)
	} else if i > len(_this.arr) {
		err = fmt.Errorf("out of range")
	} else {
		v = _this.arr[i]
	}
	return
}

/*
Str
Just string or number
*/
func (_this *Value) Str() (v string, err error) {
	if _this.typ&(String|Number) == 0 {
		err = _this.typ.Error(String | Number)
	} else {
		v = _this.str
	}
	return
}

func (_this *Value) Number() (v json.Number, err error) {
	if _this.typ&Number == 0 {
		err = _this.typ.Error(Number)
	} else {
		v = json.Number(_this.str)
	}
	return
}

func (_this *Value) Bool() (v bool, err error) {
	if _this.typ&Boolean == 0 {
		err = _this.typ.Error(Boolean)
	} else {
		v = _this.boolean
	}
	return
}

func (_this *Value) Null() (v interface{}, err error) {
	if _this.typ&Null == 0 {
		err = _this.typ.Error(Null)
	} else {
		v = nil
	}
	return
}

func (_this *Value) Type() Type {
	return _this.typ
}

func (_this *Value) Interface() (val interface{}) {
	type tempPair struct {
		src *Value
		dst reflect.Value

		flag      bool
		key       reflect.Value
		parentMap reflect.Value
	}
	vList := golist.New[tempPair]()
	vList.PushBack(tempPair{
		_this, reflect.ValueOf(&val).Elem(),
		false, reflect.Value{}, reflect.Value{}})

	for vList.Len() != 0 {
		back := vList.Back()
		vList.Remove(back)

		pair := back.Value
		src := pair.src
		dst := pair.dst

		switch src.typ {
		case String:
			dst.Set(reflect.ValueOf(src.str))
		case Number:
			dst.Set(reflect.ValueOf(json.Number(src.str)))
		case Boolean:
			dst.Set(reflect.ValueOf(src.boolean))
		case Null:
			dst.Set(reflect.New(interfaceType).Elem())
		case Array:
			dst.Set(reflect.MakeSlice(sliceType, len(src.arr), len(src.arr)))
			for i := range src.arr {
				vList.PushBack(tempPair{
					src.arr[i], dst.Elem().Index(i),
					false, reflect.Value{}, reflect.Value{}})
			}
		case Object:
			dst.Set(reflect.MakeMap(mapType))
			for k := range src.obj {
				key := reflect.ValueOf(k)
				vList.PushBack(tempPair{
					src.obj[k], reflect.New(interfaceType).Elem(),
					true, key, dst})
			}
		default:
			panic(fmt.Errorf("unknown type, cannot convert to json type"))
		}

		if pair.flag {
			pair.parentMap.Elem().SetMapIndex(pair.key, dst)
		}
	}

	return
}

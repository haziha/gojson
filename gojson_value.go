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
				ptrList.PushBack(ptrBack.Value(key.String()))
			} else if key.Kind() == reflect.Slice {
				if key.IsNil() {
					err = fmt.Errorf("key must be string, but empty slice")
					return
				} else if key.Type() != byteSliceType {
					err = fmt.Errorf("key must be string, but not byte slice")
					return
				}
				ptrList.PushBack(ptrBack.Value(string(key.Bytes())))
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
				ptrList.PushBack(ptrBack.Index(int(index.Int())))
			} else if index.CanUint() {
				ptrList.PushBack(ptrBack.Index(int(index.Uint())))
			} else if index.CanFloat() {
				ptrList.PushBack(ptrBack.Index(int(index.Float())))
			} else if index.Kind() == reflect.String {
				jN := json.Number(index.String())
				var i64 int64
				if i64, err = jN.Int64(); err != nil {
					err = fmt.Errorf("index must be integer, but string and cannot conver to integer: %s", jN)
					return
				}
				ptrList.PushBack(ptrBack.Index(int(i64)))
			} else {
				err = fmt.Errorf("index must be integer, but %v", index.Kind())
			}
		default:
			err = fmt.Errorf("cannot get element in %v", ptrBack.Type())
		}
	}

	return ptrList.Back().Value, nil
}

func (_this *Value) Int64() int64 {
	num := _this.Number()
	i64, err := num.Int64()
	if err != nil {
		panic(err)
	}
	return i64
}

func (_this *Value) Float64() float64 {
	num := _this.Number()
	f64, err := num.Float64()
	if err != nil {
		panic(err)
	}
	return f64
}

func (_this *Value) Keys() []string {
	if _this.typ&Object == 0 {
		panic(_this.typ.Error(Object))
	}
	strList := make([]string, 0, len(_this.obj))

	for k := range _this.obj {
		strList = append(strList, k)
	}

	return strList
}

func (_this *Value) Value(key string) *Value {
	if _this.typ&Object == 0 {
		panic(_this.typ.Error(Object))
	}
	return _this.obj[key]
}

/*
Len
Just object or array
*/
func (_this *Value) Len() int {
	if _this.typ&(Object|Array) == 0 {
		panic(_this.typ.Error(Object | Array))
	}
	if _this.typ == Object {
		return len(_this.obj)
	} else {
		return len(_this.arr)
	}
}

func (_this *Value) Index(i int) *Value {
	if _this.typ&Array == 0 {
		panic(_this.typ.Error(Array))
	}
	return _this.arr[i]
}

/*
Str
Just string or number
*/
func (_this *Value) Str() string {
	if _this.typ&(String|Number) == 0 {
		panic(_this.typ.Error(String | Number))
	}
	return _this.str
}

func (_this *Value) Number() json.Number {
	if _this.typ&Number == 0 {
		panic(_this.typ.Error(Number))
	}
	return json.Number(_this.str)
}

func (_this *Value) Bool() bool {
	if _this.typ&Boolean == 0 {
		panic(_this.typ.Error(Boolean))
	}
	return _this.boolean
}

func (_this *Value) Null() interface{} {
	if _this.typ&Null == 0 {
		panic(_this.typ.Error(Null))
	}
	return nil
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

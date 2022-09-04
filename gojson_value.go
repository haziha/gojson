package gojson

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"sync"
)

var valuePool = sync.Pool{
	New: func() interface{} {
		return new(Value)
	},
}

func NewValue(typ Type, val interface{}) (v *Value, err error) {
	v = valuePool.Get().(*Value)
	err = v.SetValue(typ, val)
	return v, err
}

type Value struct {
	typ     Type
	str     string            // String / Number
	boolean bool              // Boolean
	arr     []*Value          // Array
	obj     map[string]*Value // Object
}

func (_this *Value) Get(k ...interface{}) (val *Value, err error) {
	ptr := _this
	kVal := reflect.ValueOf(k)

	for i := range k {
		switch ptr.Type() {
		case Object:
			k := kVal.Index(i)
			if !k.IsValid() {
				err = fmt.Errorf("key must be string, but empty")
				return
			}
			k = k.Elem()
			var key string
			if k.Kind() == reflect.String {
				key = k.String()
			} else if k.Kind() == reflect.Slice && k.Index(i).IsValid() && k.Type() == byteType {
				key = string(k.Bytes())
			} else {
				err = fmt.Errorf("key must be string, but %s: %v", k.Kind(), k.Interface())
				return
			}
			ptr, err = ptr.Value(key)
			if err != nil {
				return
			}
		case Array:
			k := kVal.Index(i)
			if !k.IsValid() {
				err = fmt.Errorf("key must be integer, but empty")
				return
			}
			k = k.Elem()
			var index int
			if k.CanInt() {
				index = int(k.Int())
			} else if k.CanUint() {
				index = int(k.Uint())
			} else if k.CanFloat() {
				index = int(k.Float())
			} else {
				err = fmt.Errorf("index must be integer, but %s: %v", kVal.Index(i).Kind(), kVal.Index(i).Interface())
				return
			}
			ptr, err = ptr.Index(index)
			if err != nil {
				return
			}
		default:
			err = fmt.Errorf("cannot get element in %s", ptr.Type())
			return
		}
	}

	return ptr, nil
}

func (_this *Value) Len() (length int, err error) {
	if _this.typ&(Object|Array) != 0 {
		if _this.typ == Object {
			length = len(_this.obj)
		} else {
			length = len(_this.arr)
		}
		return
	}
	return -1, _this.typ.Error(Object | Array)
}

func (_this *Value) Keys() (keys []string, err error) {
	if _this.typ&Object != 0 {
		keys = make([]string, 0, len(_this.obj))
		for k := range _this.obj {
			keys = append(keys, k)
		}
		return
	}
	return nil, _this.typ.Error(Object)
}

func (_this *Value) Value(key string) (val *Value, err error) {
	if _this.typ&Object != 0 {
		var ok bool
		val, ok = _this.obj[key]
		if !ok {
			err = fmt.Errorf("not found \"%s\" in object", key)
			return
		}
		return
	}
	return nil, _this.typ.Error(Object)
}

func (_this *Value) Index(index int) (val *Value, err error) {
	if _this.typ&Array != 0 {
		if index >= 0 && index < len(_this.arr) {
			val = _this.arr[index]
			return
		}
		err = fmt.Errorf("out of slice")
		return
	}
	return nil, _this.typ.Error(Array)
}

func (_this *Value) Null() (interface{}, error) {
	if _this.typ&Null != 0 {
		return nil, nil
	}
	return nil, _this.typ.Error(Null)
}

func (_this *Value) Boolean() (bool, error) {
	if _this.typ&Boolean != 0 {
		return _this.boolean, nil
	}
	return false, _this.typ.Error(Boolean)
}

func (_this *Value) Float64() (float64, error) {
	if _this.typ&Number != 0 {
		f64, _ := json.Number(_this.str).Float64()
		return f64, nil
	}
	return 0, _this.typ.Error(Number)
}

func (_this *Value) Int64() (int64, error) {
	if _this.typ&Number != 0 {
		i64, _ := json.Number(_this.str).Int64()
		return i64, nil
	}
	return 0, _this.typ.Error(Number)
}

func (_this *Value) Number() (json.Number, error) {
	if _this.typ&Number != 0 {
		return json.Number(_this.str), nil
	}
	return "", _this.typ.Error(Number)
}

func (_this *Value) String() (string, error) {
	if _this.typ&(String|Number) != 0 {
		return _this.str, nil
	}
	return "", _this.typ.Error(String)
}

func (_this *Value) SetValue(typ Type, val interface{}) error {
	vVal := reflect.ValueOf(val)
	switch typ {
	case String:
		if vVal.Kind() == reflect.String {
			_this.str = vVal.String()
		} else if vVal.Kind() == reflect.Slice && vVal.IsValid() && vVal.Type() == byteType {
			_this.str = string(vVal.Bytes())
		} else {
			return fmt.Errorf("val must be string, but %v", vVal.Kind())
		}
	case Number:
	case Boolean:
		if vVal.Kind() == reflect.Bool {
			_this.boolean = vVal.Bool()
		} else {
			return fmt.Errorf("val must be bool, but %v", vVal.Kind())
		}
	case Null:
	case Array:
		if vVal.Kind() == sliceType.Kind() {
			if vVal.Type() == sliceType {
				_this.arr = val.([]*Value)
			} else {
				return fmt.Errorf("val must be %v, but %v", sliceType, vVal.Type())
			}
		} else {
			return fmt.Errorf("val must be %v, but %v", sliceType.Kind(), vVal.Kind())
		}
	case Object:
		if vVal.Kind() == objectType.Kind() {
			if vVal.Type() == objectType {
				_this.obj = val.(map[string]*Value)
			} else {
				return fmt.Errorf("val must be %v, but %v", objectType, vVal.Type())
			}
		} else {
			return fmt.Errorf("val must be %v, but %v", objectType.Kind(), vVal.Kind())
		}
	default:
		return fmt.Errorf("unknown type: %d", typ)
	}
	_this.typ = typ

	if _this.typ&Number != 0 {
		if vVal.CanInt() {
			_this.str = strconv.FormatInt(vVal.Int(), 10)
		} else if vVal.CanUint() {
			_this.str = strconv.FormatUint(vVal.Uint(), 10)
		} else if vVal.CanFloat() {
			_this.str = strconv.FormatFloat(vVal.Float(), 'f', -1, 64)
		} else if vVal.Kind() == reflect.String {
			if _, err := json.Number(vVal.String()).Float64(); err == nil {
				_this.str = vVal.String()
			} else if _, err = json.Number(vVal.String()).Int64(); err == nil {
				_this.str = vVal.String()
			} else {
				return fmt.Errorf("value must be convertible to integer or float, but \"%v\" cannot", val)
			}
		} else if vVal.Kind() == reflect.Slice && vVal.IsValid() && vVal.Type() == byteType {
			bytes := vVal.Bytes()
			jN := json.Number(bytes)
			if _, err := jN.Float64(); err == nil {
				_this.str = jN.String()
			} else if _, err = jN.Int64(); err == nil {
				_this.str = jN.String()
			} else {
				return fmt.Errorf("value must be convertible to integer or float, but \"%v\" cannot", jN.String())
			}
		} else {
			return fmt.Errorf("value must be convertible to integer or float, but \"%v\" cannot", val)
		}
	}

	return nil
}

func (_this *Value) MustLen() (length int) {
	length, err := _this.Len()
	if err != nil {
		panic(err)
	}
	return
}

func (_this *Value) MustKeys() (keys []string) {
	keys, err := _this.Keys()
	if err != nil {
		panic(err)
	}
	return
}

func (_this *Value) MustValue(key string) (val *Value) {
	val, err := _this.Value(key)
	if err != nil {
		panic(err)
	}
	return
}

func (_this *Value) MustIndex(index int) (val *Value) {
	val, err := _this.Index(index)
	if err != nil {
		panic(err)
	}
	return
}

func (_this *Value) MustNull() interface{} {
	_, err := _this.Null()
	if err != nil {
		panic(err)
	}
	return nil
}

func (_this *Value) MustBoolean() bool {
	b, err := _this.Boolean()
	if err != nil {
		panic(err)
	}
	return b
}

func (_this *Value) MustFloat64() float64 {
	f64, err := _this.Float64()
	if err != nil {
		panic(err)
	}
	return f64
}

func (_this *Value) MustInt64() int64 {
	i64, err := _this.Int64()
	if err != nil {
		panic(err)
	}
	return i64
}

func (_this *Value) MustNumber() json.Number {
	jN, err := _this.Number()
	if err != nil {
		panic(err)
	}
	return jN
}

func (_this *Value) MustString() string {
	str, err := _this.String()
	if err != nil {
		panic(err)
	}
	return str
}

func (_this *Value) IsObject() bool {
	return _this.typ&Object != 0
}

func (_this *Value) IsArray() bool {
	return _this.typ&Array != 0
}

func (_this *Value) IsNull() bool {
	return _this.typ&Null != 0
}

func (_this *Value) IsBoolean() bool {
	return _this.typ&Boolean != 0
}

func (_this *Value) IsNumber() bool {
	return _this.typ&Number != 0
}

func (_this *Value) IsString() bool {
	return _this.typ&String != 0
}

func (_this *Value) Type() Type {
	return _this.typ
}

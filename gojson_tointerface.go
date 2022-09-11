package gojson

import (
	"github.com/haziha/golist"
	"reflect"
	"strings"
)

func (_this *Value) equalReflectKind(v reflect.Kind) bool {
	return mapKind2Type[v] == _this.typ
}

func (_this *Value) ToInterface(s any) {
	if reflect.ValueOf(s).Kind() != reflect.Pointer {
		return
	}

	type tempPair struct {
		dst *Value
		src reflect.Value

		flag      bool
		key       reflect.Value
		parentMap reflect.Value
	}
	vList := golist.New[tempPair]()
	vList.PushBack(tempPair{_this, reflect.ValueOf(s).Elem(), false, reflect.Value{}, reflect.Value{}})

	for vList.Len() != 0 {
		back := vList.Back()
		vList.Remove(back)
		dst := back.Value.dst
		src := back.Value.src

		if dst.equalReflectKind(src.Kind()) {
			switch dst.typ {
			case String, Number, Boolean, Null: // base type
				if !src.CanSet() {
					continue
				}
				// can set
				switch dst.typ {
				case String:
					src.SetString(dst.str)
				case Number:
					if src.CanUint() {
						src.SetUint(uint64(dst.MustInt64()))
					} else if src.CanInt() {
						src.SetInt(dst.MustInt64())
					} else if src.CanFloat() {
						src.SetFloat(dst.MustFloat64())
					}
				case Boolean:
					src.SetBool(dst.boolean)
				case Null:
					src.Set(reflect.Zero(src.Type()))
				}
			case Object:
				switch src.Kind() {
				case reflect.Map:
					if src.CanSet() {
						src.Set(reflect.MakeMap(src.Type()))
						keys := dst.MustKeys()
						for k := range keys {
							vList.PushBack(tempPair{dst.MustValue(k), reflect.New(src.Type().Elem()).Elem(), true, reflect.ValueOf(k), src})
						}
					}
				case reflect.Struct:
					sType := src.Type()
					for i := 0; i < sType.NumField(); i++ {
						field := sType.Field(i)
						if !field.IsExported() {
							continue
						}
						if gj, ok := field.Tag.Lookup("gojson"); ok && gj != "" && len(strings.Split(gj, "/")) >= 2 {
							paths := strings.Split(gj, "/")
							paths = paths[1:]
							p := make([]interface{}, 0, len(paths))
							for j := range paths {
								p = append(p, strings.ReplaceAll(strings.ReplaceAll(paths[j], "~1", "/"), "~0", "~"))
							}
							dVal, err := dst.Get(p...)
							if err != nil {
								continue
							}
							vList.PushBack(tempPair{dVal, src.Field(i), false, reflect.Value{}, reflect.Value{}})
						} else {
							keys := dst.MustKeys()
							if _, ok = keys[field.Name]; !ok {
								continue
							}
							vList.PushBack(tempPair{dst.MustValue(field.Name), src.Field(i), false, reflect.Value{}, reflect.Value{}})
						}
					}
				}
			case Array:
				switch src.Kind() {
				case reflect.Array:
					for i := 0; i < src.Len() && i < dst.MustLen(); i++ {
						vList.PushBack(tempPair{dst.MustIndex(i), src.Index(i), false, reflect.Value{}, reflect.Value{}})
					}
				case reflect.Slice:
					if src.CanSet() {
						src.Set(reflect.MakeSlice(src.Type(), dst.MustLen(), dst.MustLen()))
						for i := 0; i < dst.MustLen(); i++ {
							vList.PushBack(tempPair{dst.MustIndex(i), src.Index(i), false, reflect.Value{}, reflect.Value{}})
						}
					}
				}
			}
		} else {
			switch src.Kind() {
			case reflect.Interface:
				if src.CanSet() {
					vDst := dst.Interface()
					if reflect.ValueOf(vDst).CanConvert(src.Type()) {
						src.Set(reflect.ValueOf(vDst).Convert(src.Type()))
					}
				}
			case reflect.Pointer:
				if src.CanSet() {
					if dst.equalReflectKind(src.Type().Elem().Kind()) {
						src.Set(reflect.New(src.Type().Elem()))
						vList.PushBack(tempPair{dst, src.Elem(), false, reflect.Value{}, reflect.Value{}})
					} else {
						src.Set(reflect.Zero(src.Type()))
					}
				}
			}
		}

		if back.Value.flag {
			back.Value.parentMap.SetMapIndex(back.Value.key, src)
		}
	}
}

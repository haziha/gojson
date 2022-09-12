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
		src *Value
		dst reflect.Value

		flag      bool
		key       reflect.Value
		parentMap reflect.Value
	}
	vList := golist.New[tempPair]()
	vList.PushBack(tempPair{_this, reflect.ValueOf(s).Elem(), false, reflect.Value{}, reflect.Value{}})

	for vList.Len() != 0 {
		back := vList.Back()
		vList.Remove(back)
		src := back.Value.src
		dst := back.Value.dst

		if src.equalReflectKind(dst.Kind()) {
			switch src.typ {
			case String, Number, Boolean, Null: // base type
				if !dst.CanSet() {
					continue
				}
				// can set
				switch src.typ {
				case String:
					dst.SetString(src.str)
				case Number:
					if dst.CanUint() {
						dst.SetUint(uint64(src.MustInt64()))
					} else if dst.CanInt() {
						dst.SetInt(src.MustInt64())
					} else if dst.CanFloat() {
						dst.SetFloat(src.MustFloat64())
					}
				case Boolean:
					dst.SetBool(src.boolean)
				case Null:
					dst.Set(reflect.Zero(dst.Type()))
				}
			case Object:
				switch dst.Kind() {
				case reflect.Map:
					if dst.CanSet() {
						dst.Set(reflect.MakeMap(dst.Type()))
						keys := src.MustKeys()
						for k := range keys {
							vList.PushBack(tempPair{src.MustValue(k), reflect.New(dst.Type().Elem()).Elem(), true, reflect.ValueOf(k), dst})
						}
					}
				case reflect.Struct:
					sType := dst.Type()
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
							dVal, err := src.Get(p...)
							if err != nil {
								continue
							}
							vList.PushBack(tempPair{dVal, dst.Field(i), false, reflect.Value{}, reflect.Value{}})
						} else {
							keys := src.MustKeys()
							if _, ok = keys[field.Name]; ok {
								vList.PushBack(tempPair{src.MustValue(field.Name), dst.Field(i), false, reflect.Value{}, reflect.Value{}})
							} else {
								fReplace := strings.ReplaceAll(field.Name, "_", "")
								fReplace = strings.ReplaceAll(fReplace, "-", "")
								fReplace = strings.ToLower(fReplace)
								for k := range keys {
									kReplace := strings.ReplaceAll(k, "_", "")
									kReplace = strings.ReplaceAll(kReplace, "-", "")
									kReplace = strings.ToLower(kReplace)

									if fReplace == kReplace {
										vList.PushBack(tempPair{src.MustValue(k), dst.Field(i), false, reflect.Value{}, reflect.Value{}})
										break
									}
								}
							}
						}
					}
				}
			case Array:
				switch dst.Kind() {
				case reflect.Array:
					for i := 0; i < dst.Len() && i < src.MustLen(); i++ {
						vList.PushBack(tempPair{src.MustIndex(i), dst.Index(i), false, reflect.Value{}, reflect.Value{}})
					}
				case reflect.Slice:
					if dst.CanSet() {
						dst.Set(reflect.MakeSlice(dst.Type(), src.MustLen(), src.MustLen()))
						for i := 0; i < src.MustLen(); i++ {
							vList.PushBack(tempPair{src.MustIndex(i), dst.Index(i), false, reflect.Value{}, reflect.Value{}})
						}
					}
				}
			}
		} else {
			switch dst.Kind() {
			case reflect.Interface:
				if dst.CanSet() {
					vDst := src.Interface()
					if reflect.ValueOf(vDst).CanConvert(dst.Type()) {
						dst.Set(reflect.ValueOf(vDst).Convert(dst.Type()))
					}
				}
			case reflect.Pointer:
				if dst.CanSet() {
					if src.equalReflectKind(dst.Type().Elem().Kind()) {
						dst.Set(reflect.New(dst.Type().Elem()))
						vList.PushBack(tempPair{src, dst.Elem(), false, reflect.Value{}, reflect.Value{}})
					} else {
						dst.Set(reflect.Zero(dst.Type()))
					}
				}
			}
		}

		if back.Value.flag {
			back.Value.parentMap.SetMapIndex(back.Value.key, dst)
		}
	}
}

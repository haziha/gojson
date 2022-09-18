package gojson

import (
	"fmt"
	"github.com/haziha/golist"
	"reflect"
	"strings"
)

/*
变量为零值规则
	1.	基础类型(非复合类型, number/bool/string):
			可能确实是零值(不告知父级), 也可能赋值失败(会告知父级)
	2.	指针/interface:
			源数据不存在与之对应的数据(告知父级), 或子元素赋值失败(不告知父级)
	3.	map/slice:
			源数据类型不匹配(非 object / array)[会告知父级]
	4.	struct/array:
			源数据类型不匹配(非 object / array)[会告知父级], 或确实是零值(不告知父级)
	5.	struct field/array element/map element/slice element:
			成员/子元素赋值失败, 或确实是零值
			或不存在(
				仅限struct;
				map只要有key, 则 if ok 一定是 true;
				array 和 slice 则需要对应 index, 所以不能把元素去掉, 只能置空值
			)
			均不告知父级
*/

type taskType int

const (
	baseTask taskType = 1 << iota
	mapTask
	checkTask
)

type taskPair struct {
	typ taskType

	// baseTask
	src    *Value
	dst    reflect.Value
	parent *taskPair

	// mapTask
	key       reflect.Value
	parentMap reflect.Value

	// mapTask or checkTask
	childFlag bool
}

type KeyTagName string
type PathTagName string

func newMapCheckTask(key reflect.Value, parentMap reflect.Value) (tp *taskPair) {
	return &taskPair{
		typ:       mapTask | checkTask,
		key:       key,
		parentMap: parentMap,
		childFlag: true,
	}
}

func newMapTask(src *Value, dst reflect.Value, parent *taskPair, key reflect.Value, parentMap reflect.Value) (tp *taskPair) {
	return &taskPair{
		typ:       mapTask,
		src:       src,
		dst:       dst,
		parent:    parent,
		key:       key,
		parentMap: parentMap,
	}
}

func newCheckTask(dst reflect.Value) (tp *taskPair) {
	return &taskPair{
		typ:       checkTask,
		dst:       dst,
		childFlag: true,
	}
}

func newBaseTask(src *Value, dst reflect.Value, parent *taskPair) (tp *taskPair) {
	return &taskPair{
		typ:    baseTask,
		src:    src,
		dst:    dst,
		parent: parent,
	}
}

func (_this *Value) ToInterface(s any, options ...interface{}) (err error) {
	keyTagName := "gk"
	pathTagName := "gp"
	for _, arg := range options {
		switch v := arg.(type) {
		case KeyTagName:
			keyTagName = string(v)
		case PathTagName:
			pathTagName = string(v)
		}
	}

	if reflect.ValueOf(s).Kind() != reflect.Pointer {
		err = fmt.Errorf("ToInterface(non-pointer)")
		return
	}
	if reflect.ValueOf(s).Elem().Kind() == reflect.Invalid {
		err = fmt.Errorf("ToInterface(nil)")
		return
	}

	tList := golist.New[*taskPair]()
	tList.PushBack(newCheckTask(reflect.ValueOf(s).Elem()))
	tList.PushBack(newBaseTask(_this, reflect.ValueOf(s).Elem(), tList.Back().Value))

	for tList.Len() != 0 {
		back := tList.Back()
		tList.Remove(back)
		task := back.Value

		if task.typ == baseTask || task.typ == mapTask {
			src := task.src
			dst := task.dst

			if src.typ == Boolean && dst.Kind() == reflect.Bool {
				dst.SetBool(src.boolean)
			} else if src.typ == String && dst.Kind() == reflect.String {
				dst.SetString(src.str)
			} else if src.typ == Number && dst.Kind() >= reflect.Int && dst.Kind() <= reflect.Int64 {
				dst.SetInt(src.MustInt64())
			} else if src.typ == Number && dst.Kind() >= reflect.Uint && dst.Kind() <= reflect.Uint64 {
				dst.SetUint(src.MustUint64())
			} else if src.typ == Number && dst.Kind() >= reflect.Float32 && dst.Kind() <= reflect.Float64 {
				dst.SetFloat(src.MustFloat64())
			} else if src.typ == Object && dst.Kind() == reflect.Map {
				if dst.IsZero() {
					dst.Set(reflect.MakeMap(dst.Type()))
				}
				for k := range src.obj {
					tList.PushBack(newMapCheckTask(reflect.ValueOf(k), dst))
					tList.PushBack(newMapTask(src.MustValue(k), reflect.New(dst.Type().Elem()).Elem(), tList.Back().Value, reflect.ValueOf(k), dst))
				}
			} else if src.typ == Array && dst.Kind() == reflect.Slice {
				if dst.IsZero() {
					dst.Set(reflect.MakeSlice(dst.Type(), src.MustLen(), src.MustLen()))
				}
				for i := 0; i < src.MustLen() && i < dst.Len(); i++ {
					tList.PushBack(newBaseTask(src.MustIndex(i), dst.Index(i), nil))
				}
			} else if src.typ == Object && dst.Kind() == reflect.Struct {
				// 需过滤私有成员, 虽然能用 reflect 和 unsafe 强行修改, 但没必要
				dstTyp := dst.Type()
				keys := src.MustKeys()
				replacer := strings.NewReplacer("_", "", "-", "")
				for i := 0; i < dstTyp.NumField(); i++ {
					if !dstTyp.Field(i).IsExported() {
						continue
					}
					// gp 优先于 gk
					// 是否指定路径
					if gp := dstTyp.Field(i).Tag.Get(pathTagName); len(gp) > 0 && gp[0] == '/' {
						k := strings.Split(gp, "/")[1:]
						ki := make([]interface{}, 0, len(k))
						for j := range k {
							ki = append(ki, strings.ReplaceAll(strings.ReplaceAll(k[j], "~1", "/"), "~0", "~"))
						}
						if src, err := src.Get(ki...); err == nil && src != nil {
							tList.PushBack(newCheckTask(dst.Field(i)))
							tList.PushBack(newBaseTask(src, dst.Field(i), tList.Back().Value))
							continue
						}
					}
					// 是否有指定 key
					if gk, ok := dstTyp.Field(i).Tag.Lookup(keyTagName); ok {
						tList.PushBack(newCheckTask(dst.Field(i)))
						if _, ok = keys[gk]; ok {
							tList.PushBack(newBaseTask(src.MustValue(gk), dst.Field(i), tList.Back().Value))
						} else {
							tList.Back().Value.childFlag = false
						}
						continue
					} else {
						gk = strings.ToLower(replacer.Replace(dstTyp.Field(i).Name))
						for j := range keys {
							k := strings.ToLower(replacer.Replace(j))
							if gk == k {
								tList.PushBack(newCheckTask(dst.Field(i)))
								tList.PushBack(newBaseTask(src.MustValue(j), dst.Field(i), tList.Back().Value))
								break
							}
						}
					}
				}
			} else if src.typ == Array && dst.Kind() == reflect.Array {
				for i := 0; i < src.MustLen() && i < dst.Len(); i++ {
					tList.PushBack(newCheckTask(dst.Index(i)))
					tList.PushBack(newBaseTask(src.MustIndex(i), dst.Index(i), tList.Back().Value))
				}
			} else if src.typ == Null {
				if task.parent != nil {
					task.parent.childFlag = false
				}
				dst.Set(reflect.Zero(dst.Type()))
			} else if dst.Kind() == reflect.Pointer {
				if dst.Elem().Kind() == reflect.Invalid {
					dst.Set(reflect.New(dst.Type().Elem()))
				}
				tList.PushBack(newCheckTask(dst))
				tList.PushBack(newBaseTask(src, dst.Elem(), tList.Back().Value))
			} else if dst.Kind() == reflect.Interface {
				if dst.Elem().Kind() == reflect.Invalid {
					src_ := reflect.ValueOf(src.Interface())
					if src_.CanConvert(dst.Type()) {
						dst.Set(src_.Convert(dst.Type()))
					} else if task.parent != nil {
						task.parent.childFlag = false
					}
				} else {
					tList.PushBack(newCheckTask(dst))
					tList.PushBack(newBaseTask(src, dst.Elem(), tList.Back().Value))
				}
			} else {
				if task.parent != nil {
					task.parent.childFlag = false
				}
				continue
			}

			if task.typ == mapTask {
				task.parentMap.SetMapIndex(task.key, dst)
			}
		} else if task.typ == (mapTask | checkTask) {
			if !task.childFlag {
				task.parentMap.SetMapIndex(task.key, reflect.Zero(task.parentMap.Type().Elem()))
			}
		} else if task.typ == checkTask {
			if !task.childFlag {
				task.dst.Set(reflect.Zero(task.dst.Type()))
			}
		} else {
			panic("unknown type")
		}
	}

	return
}

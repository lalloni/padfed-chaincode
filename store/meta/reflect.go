package meta

import (
	"reflect"
)

func FieldSetter(name string) SetterFunc {
	return func(v interface{}, w interface{}) {
		reflect.ValueOf(v).Elem().FieldByName(name).Set(reflect.ValueOf(w))
	}
}

func FieldGetter(name string) GetterFunc {
	return func(v interface{}) interface{} {
		return reflect.ValueOf(v).Elem().FieldByName(name).Interface()
	}
}

func MapEnumerator(mapGetter GetterFunc) EnumeratorFunc {
	return func(v interface{}) []Item {
		items := []Item{}
		mv := reflect.ValueOf(mapGetter(v))
		vs := mv.MapKeys()
		for _, v := range vs {
			items = append(items, NewItem(v.String(), mv.MapIndex(v).Interface()))
		}
		return items
	}
}

func MapCollector(mapGetter GetterFunc) CollectorFunc {
	return func(v interface{}, item Item) {
		reflect.ValueOf(mapGetter(v)).SetMapIndex(reflect.ValueOf(item.Identifier), reflect.ValueOf(item.Value))
	}
}

func ValueCreator(ref interface{}) CreatorFunc {
	t := reflect.TypeOf(ref).Elem()
	return func() interface{} {
		return reflect.New(t).Interface()
	}
}

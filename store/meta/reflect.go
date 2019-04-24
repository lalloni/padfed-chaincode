package meta

import "reflect"

func FieldSetter(name string) SetterFunc {
	return func(v interface{}, w interface{}) {
		reflect.ValueOf(v).Elem().FieldByName(name).Set(reflect.ValueOf(w))
	}
}

func FieldGetter(name string) GetterFunc {
	return func(v interface{}) interface{} {
		return reflect.ValueOf(v).FieldByName(name).Interface()
	}
}

func MapEnumerator(mapGetter GetterFunc) EnumeratorFunc {
	return func(v interface{}) []Item {
		items := []Item{}
		mi := reflect.ValueOf(mapGetter(v)).MapRange()
		for mi.Next() {
			items = append(items, NewItem(mi.Key().String(), mi.Value().Interface()))
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
		return reflect.New(t)
	}
}

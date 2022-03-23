package reflectricity

import (
	"encoding/gob"
	"reflect"
)

// RegisterGob recursively registers a value.
// Useful when dealing with polymorphic values that need to be written across
// the wire.
func RegisterGob(value any) {
	registerRecursiveGob(value, make(map[reflect.Type]struct{}))
}

func registerRecursiveGob(i any, visited map[reflect.Type]struct{}) {
	if i == nil {
		return
	}

	t := reflect.ValueOf(i)
	kind := t.Kind()
	if kind == reflect.Ptr {
		if k := t.Elem().Kind(); k == reflect.Struct {
			gob.Register(i)
			registerStructFields(t.Elem().Interface(), visited)
		}
	}

	if _, ok := visited[t.Type()]; ok {
		return
	}
	visited[t.Type()] = struct{}{}
	switch kind {
	case reflect.Map:
		gob.Register(i)
		iter := t.MapRange()
		for iter.Next() {
			if iter.Key().Kind() == reflect.Struct {
				registerRecursiveGob(iter.Key().Interface(), visited)
			}
			if iter.Value().Kind() == reflect.Struct {
				registerRecursiveGob(iter.Value().Interface(), visited)
			}
		}
	case reflect.Array, reflect.Slice:
		gob.Register(i)
		for i := 0; i < t.Len(); i++ {
			registerRecursiveGob(t.Index(i).Interface(), visited)
		}
	case reflect.Struct:
		gob.Register(i)
		registerStructFields(i, visited)
	}
}

func registerStructFields(i any, visited map[reflect.Type]struct{}) {
	t := reflect.ValueOf(i)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		registerRecursiveGob(field.Interface(), visited)
	}
}

package reflectricity

import (
	"encoding/gob"
	"reflect"
	"sync"
)

var registeredNames sync.Map

// RegisterGob recursively registers a value.
// Useful when dealing with polymorphic values that need to be written across
// the wire.
func RegisterGob(value any) {
	registerRecursiveGob(value)
}

func registerRecursiveGob(i any) {
	if i == nil {
		return
	}

	t := reflect.ValueOf(i)
	kind := t.Kind()
	switch kind {
	case reflect.Ptr:
		if k := t.Elem().Kind(); k == reflect.Struct {
			if _, ok := registeredNames.Load(t.Type()); !ok {
				gob.Register(i)
				registeredNames.Store(t.Type(), struct{}{})
			}
		}
		registerStructFields(t.Elem().Interface())
	case reflect.Map:
		iter := t.MapRange()
		for iter.Next() {
			registerRecursiveGob(iter.Key().Interface())
			registerRecursiveGob(iter.Value().Interface())
		}
	case reflect.Array, reflect.Slice:
		for n := 0; n < t.Len(); n++ {
			registerRecursiveGob(t.Index(n).Interface())
		}
	case reflect.Struct:
		if _, ok := registeredNames.Load(t.Type()); !ok {
			gob.Register(i)
			registeredNames.Store(t.Type(), struct{}{})
		}
		registerStructFields(i)
	}
}

func registerStructFields(i any) {
	t := reflect.ValueOf(i)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		registerRecursiveGob(field.Interface())
	}
}

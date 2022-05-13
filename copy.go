package reflectricity

import "reflect"

func DeepCopy[T any](value T, private bool) T {
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr && val.IsNil() {
		return value
	}
	return deepCopy(val, private).Interface().(T)
}

func deepCopy(val reflect.Value, p bool) reflect.Value {
	if val.Kind() == reflect.Ptr && val.IsNil() {
		return val
	}

	val, depth := ptrUnwrap(val)
	t := val.Type()
	switch t.Kind() {
	case reflect.Map:
		s := reflect.MakeMap(t)
		iter := val.MapRange()
		for iter.Next() {
			s.SetMapIndex(deepCopy(iter.Key(), p), deepCopy(iter.Value(), p))
		}
		val = s
	case reflect.Chan:
		val = reflect.MakeChan(t, val.Cap())
	case reflect.Array, reflect.Slice:
		s := reflect.MakeSlice(t, val.Len(), val.Len())
		for i := 0; i < val.Len(); i++ {
			s.Index(i).Set(deepCopy(val.Index(i), p))
		}
		val = s
	case reflect.Struct:
		s := reflect.New(t).Elem()
		for i := 0; i < val.NumField(); i++ {
			orig := val.Field(i)
			dest := s.Field(i)
			if orig.CanInterface() {
				dest.Set(deepCopy(orig, p))
			} else {
				// Is this an unholy space?
				if canExposeInterface() {
					// Perform black magic
					orig = exposePrivateValue(orig)
					dest = exposePrivateValue(dest)
					dest.Set(deepCopy(orig, p))
				}
			}
		}
		val = s
	}
	return ptrWrap(val, depth)
}

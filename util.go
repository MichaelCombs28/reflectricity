package reflectricity

import "reflect"

// Nests a type into a pointer.
// Useful for nested pointer definitions when printing a value
// or a generic way to nest something as a pointer
func ToPtr[T any](t T) *T {
	return &t
}

func ptrWrap(v reflect.Value, depth int) reflect.Value {
	for i := 0; i < depth; i++ {
		pt := reflect.PtrTo(v.Type())
		pv := reflect.New(pt.Elem())
		pv.Elem().Set(v)
		v = pv
	}
	return v
}

func ptrUnwrap(v reflect.Value) (va reflect.Value, depth int) {
	va = v
	for va.Kind() == reflect.Ptr && !va.IsNil() {
		va = va.Elem()
		depth++
	}
	return
}

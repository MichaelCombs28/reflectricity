package reflectricity

import (
	"reflect"
	"unsafe"
)

type arrayMergeStrategy int

const (
	CONCAT arrayMergeStrategy = iota
	REPLACE
)

// Merges 2 structs, field by field replacing all existing
// field in the right structure into the left. If left is nil
// returns right. If the types don't match up, returns right
func (r *Reflector) Merge(left any, right any) any {
	return r.merge(reflect.ValueOf(left), reflect.ValueOf(right)).Interface()
}

func (r *Reflector) merge(va reflect.Value, vb reflect.Value) reflect.Value {
	ta := va.Type()
	tb := vb.Type()
	if ta != tb {
		return vb
	}

	// Deref until not pointer
	va, pdepth := ptrUnwrap(va)
	vb, _ = ptrUnwrap(vb)

	if va.Interface() == nil || !va.IsValid() || va.IsZero() {
		return ptrWrap(vb, pdepth)
	}

	if vb.Interface() == nil || !vb.IsValid() || vb.IsZero() {
		return ptrWrap(va, pdepth)
	}

	switch ta.Kind() {
	case reflect.Map:
		mp := reflect.MakeMap(va.Type())
		iter := va.MapRange()
		for iter.Next() {
			v := va.MapIndex(iter.Key())
			mp.SetMapIndex(iter.Key(), r.merge(iter.Value(), v))
		}
		iter = vb.MapRange()
		for iter.Next() {
			v := vb.MapIndex(iter.Key())
			mp.SetMapIndex(iter.Key(), r.merge(iter.Value(), v))
		}
		va = mp
	case reflect.Array, reflect.Slice:
		va = r.mergeArrays(va, vb)
	case reflect.Struct:
		out := reflect.New(ta)
		for i := 0; i < va.NumField(); i++ {
			avalue := va.Field(i)
			bvalue := vb.Field(i)
			field := out.Elem().Field(i)

			if field.CanInterface() {
				m := r.merge(avalue, bvalue)

				field.Set(m)
			} else {
				//assume private
				if r.private && canExposeInterface() {
					avalue = exposePrivateValue(avalue)
					bvalue = exposePrivateValue(bvalue)
					rf := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
					rf.Set(r.merge(avalue, bvalue))
				}
			}
		}
		va = out.Elem()
	default:
		va = vb
	}
	return ptrWrap(va, pdepth)
}

func (r *Reflector) mergeArrays(a reflect.Value, b reflect.Value) reflect.Value {
	var result reflect.Value
	switch r.arrayMergeStrategy {
	case CONCAT:
		le := a.Len() + b.Len()
		result = reflect.MakeSlice(a.Type(), le, le)
		var i int
		for i = 0; i < a.Len(); i++ {
			result.Index(i).Set(a.Index(i))
		}
		for x := 0; x < b.Len(); x++ {
			result.Index(i).Set(b.Index(x))
			i++
		}
	case REPLACE:
		result = reflect.MakeSlice(a.Type(), max(a.Len(), b.Len()), max(a.Len(), b.Len()))
		var i int
		for i = 0; i < a.Len(); i++ {
			result.Index(i).Set(a.Index(i))
		}

		for x := 0; x < b.Len(); x++ {
			if x < a.Len() {
				result.Index(x).Set(r.merge(a.Index(x), b.Index(x)))
			} else {
				result.Index(x).Set(b.Index(x))
			}
		}
	}
	return result
}

func max(i1 int, i2 int) int {
	if i1 > i2 {
		return i1
	}
	return i2
}

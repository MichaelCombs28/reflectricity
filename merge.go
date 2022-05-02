package reflectricity

import (
	"reflect"
	"unsafe"
)

type arrayMergeStrategy int

const (
	APPEND arrayMergeStrategy = iota
	FULL_REPLACE
)

func MergeLeft(left any, right any) any {
	return MergeLeftWithOptions(left, right, false, APPEND)
}

// Merges 2 structs, field by field replacing all existing
// field in the right structure into the left. If left is nil
// returns right. If the types don't match up, returns right.
func MergeLeftWithOptions(left any, right any, mergePrivate bool, arrayMerge arrayMergeStrategy) any {
	return mergeLeftWithOptions(reflect.ValueOf(left), reflect.ValueOf(right), mergePrivate, arrayMerge).Interface()
}

func mergeLeftWithOptions(va reflect.Value, vb reflect.Value, mergePrivate bool, arrayMerge arrayMergeStrategy) reflect.Value {
	pdepth := 0

	ta := va.Type()
	tb := va.Type()
	if ta != tb {
		return va
	}

	// Deref until not pointer
	for va.Kind() == reflect.Ptr {
		if !va.IsNil() {
			va = va.Elem()
			pdepth++
		} else {
			break
		}
	}

	for vb.Kind() == reflect.Ptr {
		if !vb.IsNil() {
			vb = vb.Elem()
		} else {
			break
		}
	}

	if va.Interface() == nil || !va.IsValid() || va.IsZero() {
		return ptrWrap(vb, pdepth)
	}

	if vb.Interface() == nil || !vb.IsValid() || vb.IsZero() {
		return ptrWrap(va, pdepth)
	}

	switch ta.Kind() {
	case reflect.Map:
		r := reflect.MakeMap(va.Type())
		iter := va.MapRange()
		for iter.Next() {
			v := va.MapIndex(iter.Key())
			r.SetMapIndex(iter.Key(), mergeLeftWithOptions(iter.Value(), v, mergePrivate, arrayMerge))
		}
		iter = vb.MapRange()
		for iter.Next() {
			v := vb.MapIndex(iter.Key())
			r.SetMapIndex(iter.Key(), mergeLeftWithOptions(iter.Value(), v, mergePrivate, arrayMerge))
		}
		va = r
	case reflect.Array, reflect.Slice:
		va = mergeArrays(va, vb, mergePrivate, arrayMerge)
	case reflect.Struct:
		out := reflect.New(ta)
		for i := 0; i < va.NumField(); i++ {
			avalue := va.Field(i)
			bvalue := vb.Field(i)
			field := out.Elem().Field(i)

			if field.CanInterface() {
				m := mergeLeftWithOptions(avalue, bvalue, mergePrivate, arrayMerge)

				field.Set(m)
			} else {
				//assume private
				if mergePrivate && canExposeInterface() {
					avalue = exposePrivateValue(avalue)
					bvalue = exposePrivateValue(bvalue)
					rf := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
					rf.Set(mergeLeftWithOptions(avalue, bvalue, mergePrivate, arrayMerge))
				}
			}
		}
		va = out.Elem()
	default:
		va = vb
	}
	return ptrWrap(va, pdepth)
}

func mergeArrays(a reflect.Value, b reflect.Value, mergePrivate bool, mergeStrategy arrayMergeStrategy) reflect.Value {
	var result reflect.Value
	switch mergeStrategy {
	case APPEND:
		result = reflect.MakeSlice(a.Type(), 0, 0)
		for i := 0; i < a.Len(); i++ {
			result = reflect.Append(result, a.Index(i))
		}
		for i := 0; i < b.Len(); i++ {
			result = reflect.Append(result, b.Index(i))
		}
	case FULL_REPLACE:
		result = reflect.MakeSlice(a.Type(), max(a.Len(), b.Len()), max(a.Len(), b.Len()))
		var i int
		for i = 0; i < a.Len(); i++ {
			result.Index(i).Set(a.Index(i))
		}

		for x := 0; x < b.Len(); x++ {
			if x < a.Len() {
				result.Index(x).Set(mergeLeftWithOptions(a.Index(x), b.Index(x), mergePrivate, mergeStrategy))
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

func ptrWrap(v reflect.Value, depth int) reflect.Value {
	for i := 0; i < depth; i++ {
		pt := reflect.PtrTo(v.Type())
		pv := reflect.New(pt.Elem())
		pv.Elem().Set(v)
		v = pv
	}
	return v
}
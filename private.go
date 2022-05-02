package reflectricity

// Taken directly
// from https://github.com/kstenerud/go-describe/blob/master/describe_unsafe.go

import (
	"log"
	"reflect"
	"unsafe"
)

type flag uintptr

type flagROTester struct {
	A   int
	a   int // reflect/value.go:flagStickyRO
	int     // reflect/value.go:flagEmbedRO
	// Note: flagRO = flagStickyRO | flagEmbedRO
}

var flagOffset uintptr
var maskFlagRO flag
var hasExpectedReflectStruct bool

func init() {
	if field, ok := reflect.TypeOf(reflect.Value{}).FieldByName("flag"); ok {
		flagOffset = field.Offset
	} else {
		log.Println("reflectricity: exposeInterface() is disabled because the " +
			"reflect.Value struct no longer has a flag field. Please open an " +
			"issue at https://github.com/MichaelCombs28/reflectricity/issues")
		hasExpectedReflectStruct = false
		return
	}

	rv := reflect.ValueOf(flagROTester{})
	getFlag := func(v reflect.Value, name string) flag {
		return flag(reflect.ValueOf(v.FieldByName(name)).FieldByName("flag").Uint())
	}
	flagRO := (getFlag(rv, "a") | getFlag(rv, "int")) ^ getFlag(rv, "A")
	maskFlagRO = ^flagRO

	if flagRO == 0 {
		log.Println("reflectricity: exposeInterface() is disabled because the " +
			"reflect flag type no longer has a flagEmbedRO or flagStickyRO bit. " +
			"Please open an issue at https://github.com/MichaelCombs28/reflectricity/issues")
		hasExpectedReflectStruct = false
		return
	}

	hasExpectedReflectStruct = true
}

func canExposeInterface() bool {
	return hasExpectedReflectStruct
}

func exposePrivateValue(v reflect.Value) reflect.Value {
	pFlag := (*flag)(unsafe.Pointer(uintptr(unsafe.Pointer(&v)) + flagOffset))
	*pFlag &= maskFlagRO
	return v
}

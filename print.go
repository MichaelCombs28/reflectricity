package reflectricity

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func (r *Reflector) AddAlias(pkg, alias string) {
	r.packages[pkg] = alias
}

// Print a value as code. Very useful for code generation when a library
// provides the types being generated
func (r *Reflector) PrintValue(i any) string {
	// TODO: Use runtime to pull pkg info
	// Adding pacakge import for ToPtr usage
	sb := new(strings.Builder)
	r.printValueWithDepth(sb, reflect.ValueOf(i), 0)
	return sb.String()
}

func (r *Reflector) printValueWithDepth(sb *strings.Builder, t reflect.Value, depth int) {
	t, pdepth := ptrUnwrap(t)
	if t.Kind() == reflect.Pointer && t.IsNil() {
		sb.WriteString("nil")
		return
	}
	var after []string
	typeOf := t.Type()
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		after = wrapPtr(sb, pdepth, false)
		sb.WriteString(strconv.FormatInt(t.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		sb.WriteString(strconv.FormatUint(t.Uint(), 10))
	case reflect.Bool:
		after = wrapPtr(sb, pdepth, false)
		isBool := typeOf.Name() != "bool"
		if isBool {
			r.writeType(sb, typeOf)
			sb.WriteRune('(')
		}
		sb.WriteString(fmt.Sprint(t.Bool()))
		if isBool {
			sb.WriteRune(')')
		}
	case reflect.Float64:
		after = wrapPtr(sb, pdepth, false)
		sb.WriteString(fmt.Sprint(t.Float()))
	case reflect.String:
		after = wrapPtr(sb, pdepth, false)
		isStr := typeOf.Name() != "string"
		if isStr {
			r.writeType(sb, typeOf)
			sb.WriteRune('(')
		}
		sb.WriteRune('"')
		sb.WriteString(t.String())
		sb.WriteRune('"')
		if isStr {
			sb.WriteRune(')')
		}
	case reflect.Chan:
		after = wrapPtr(sb, pdepth, false)
		sb.WriteString("make(")
		r.writeType(sb, typeOf)
		if t.Cap() > 0 {
			sb.WriteString(", ")
			sb.WriteString(strconv.Itoa(t.Cap()))
		}
		sb.WriteRune(')')
	case reflect.Map:
		after = wrapPtr(sb, pdepth, true)
		r.writeType(sb, typeOf)
		isMap := strings.HasPrefix(typeOf.String(), "map")
		if !isMap {
			sb.WriteRune('(')
			kt := typeOf.Key()
			vt := typeOf.Elem()
			sb.WriteString("map[")
			r.writeType(sb, kt)
			sb.WriteRune(']')
			r.writeType(sb, vt)
		}
		sb.WriteString("{")
		iter := t.MapRange()
		for iter.Next() {
			sb.WriteRune('\n')
			sb.WriteString(strings.Repeat("\t", depth+1))
			// Key
			r.printValueWithDepth(sb, iter.Key(), depth+1)
			sb.WriteString(": ")
			r.printValueWithDepth(sb, iter.Value(), depth+1)
			sb.WriteRune(',')
		}
		sb.WriteRune('\n')
		sb.WriteString(strings.Repeat("\t", depth))
		sb.WriteRune('}')
		if !isMap {
			sb.WriteRune(')')
		}
	case reflect.Array, reflect.Slice:
		after = wrapPtr(sb, pdepth, true)
		r.writeType(sb, typeOf)
		sb.WriteRune('{')
		for i := 0; i < t.Len(); i++ {
			sb.WriteRune('\n')
			sb.WriteString(strings.Repeat("\t", depth+1))
			r.printValueWithDepth(sb, t.Index(i), depth+1)
			sb.WriteRune(',')
		}
		sb.WriteRune('\n')
		sb.WriteString(strings.Repeat("\t", depth))
		sb.WriteRune('}')
	case reflect.Struct:
		after = wrapPtr(sb, pdepth, true)
		r.writeType(sb, typeOf)
		sb.WriteString("{\n")

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if !field.CanInterface() && !r.private {
				continue
			}

			if !field.CanInterface() {
				field = exposePrivateValue(field)
			}
			key := typeOf.Field(i).Name
			sb.WriteString(strings.Repeat("\t", depth+1))
			sb.WriteString(key)
			sb.WriteString(": ")
			r.printValueWithDepth(sb, field, depth+1)
			sb.WriteString(",\n")
		}

		sb.WriteString(strings.Repeat("\t", depth))
		sb.WriteRune('}')
	}
	// If struct

	for _, rune_ := range after {
		sb.WriteString(rune_)
	}
}

func (r *Reflector) writeType(sb *strings.Builder, t reflect.Type) {
	typeStr := t.String()
	l := len(r.packages)

	if l > 0 && strings.Contains(typeStr, ".") {
		var layers []reflect.Type
		x := t
		for x.Name() == "" {
			layers = append(layers, x)
			x = x.Elem()
		}
		for i := 0; i < len(layers); i++ {
			layer := layers[i]
			switch layer.Kind() {
			case reflect.Pointer:
				sb.WriteRune('*')
			case reflect.Array, reflect.Slice:
				sb.WriteString("[]")
			case reflect.Map:
				sb.WriteString("map[")
				r.writeType(sb, layer.Key())
				sb.WriteRune(']')
			case reflect.Chan:
				sb.WriteString("chan ")
			}
		}
		if pkg, ok := r.packages[x.PkgPath()]; ok {
			sb.WriteString(pkg)
			sb.WriteRune('.')
			sb.WriteString(x.Name())
		} else {
			sb.WriteString(t.String())
		}
	} else {
		sb.WriteString(t.String())
	}
}

func wrapPtr(sb *strings.Builder, pdepth int, wantsAmp bool) []string {
	var result []string
	for p := pdepth; p > 0; p-- {
		if p == 1 && wantsAmp {
			sb.WriteRune('&')
		} else {
			sb.WriteString("reflectricity.ToPtr(")
			result = append(result, ")")
		}
	}
	return result
}

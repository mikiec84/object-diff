package obj_diff

import (
	"reflect"
)

func CopyValueReflectively(oldValue interface{}) interface{} {
	return CopyReflectValue(reflect.ValueOf(oldValue)).Interface()
}

func CopyReflectValue(oldVal reflect.Value) (newVal reflect.Value) {
	newType := oldVal.Type()
	switch newType.Kind() {
	case reflect.Struct:
		newVal = reflect.New(newType).Elem()
		// 	newVal = reflect.Zero(newType)
		for f := 0; f < newType.NumField(); f++ {
			newVal.Field(f).Set(CopyReflectValue(oldVal.Field(f)))
		}

	case reflect.Map:
		newVal = reflect.MakeMapWithSize(newType, oldVal.Len())
		for _, key := range oldVal.MapKeys() {
			newVal.SetMapIndex(CopyReflectValue(key), CopyReflectValue(oldVal.MapIndex(key)))
		}

	case reflect.Array:
		newVal = reflect.New(newType).Elem()
		for i := 0; i < oldVal.Len(); i++ {
			newVal.Index(i).Set(CopyReflectValue(oldVal.Index(i)))
		}

	case reflect.Slice:
		newVal = reflect.MakeSlice(newType, oldVal.Len(), oldVal.Cap())
		for i := 0; i < oldVal.Len(); i++ {
			newVal.Index(i).Set(CopyReflectValue(oldVal.Index(i)))
		}

	case reflect.Ptr:
		if oldVal.IsNil() {
			newVal = reflect.Zero(newType)
		} else {
			newVal = reflect.New(newType.Elem())
			newVal.Elem().Set(CopyReflectValue(oldVal.Elem()))
		}

	default:
		newVal = copyBasic(oldVal)
	}

	return
}

func copyBasic(oldVal reflect.Value) (newVal reflect.Value) {
	switch oldVal.Kind() {
	case reflect.String:
		newVal = reflect.ValueOf(oldVal.String())

	case reflect.Int64:
		newVal = reflect.ValueOf(oldVal.Int())
	case reflect.Int32:
		newVal = reflect.ValueOf(int32(oldVal.Int()))
	case reflect.Int16:
		newVal = reflect.ValueOf(int16(oldVal.Int()))
	case reflect.Int8:
		newVal = reflect.ValueOf(int8(oldVal.Int()))
	case reflect.Int:
		newVal = reflect.ValueOf(int(oldVal.Int()))

	case reflect.Uint64:
		newVal = reflect.ValueOf(oldVal.Uint())
	case reflect.Uint32:
		newVal = reflect.ValueOf(uint32(oldVal.Uint()))
	case reflect.Uint16:
		newVal = reflect.ValueOf(uint16(oldVal.Uint()))
	case reflect.Uint8:
		newVal = reflect.ValueOf(uint8(oldVal.Uint()))
	case reflect.Uint:
		newVal = reflect.ValueOf(uint(oldVal.Uint()))

	case reflect.Float64:
		newVal = reflect.ValueOf(oldVal.Float())
	case reflect.Float32:
		newVal = reflect.ValueOf(float32(oldVal.Float()))

	case reflect.Complex128:
		newVal = reflect.ValueOf(oldVal.Complex())
	case reflect.Complex64:
		newVal = reflect.ValueOf(complex64(oldVal.Complex()))

	case reflect.Bool:
		newVal = reflect.ValueOf(oldVal.Bool())

	default:
		panic(NewPatchError("unhandled basic kind '%v'\n", oldVal.Kind()))
	}

	return
}
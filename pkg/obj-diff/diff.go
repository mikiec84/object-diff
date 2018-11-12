package obj_diff

import (
	"fmt"
	"reflect"
)

// BUG(11xor6) Renaming of Map keys results in a deletion and addition.
// BUG(11xor6) Lists with different orders but the same elements will generate changes.
func Diff(obj1 interface{}, obj2 interface{}) (*changeSet, error) {
	v1 := reflect.ValueOf(obj1)
	v2 := reflect.ValueOf(obj2)

	if v1.Type() != v2.Type() {
		return nil, fmt.Errorf("type of obj1(%T) not equal to obj2(%T)", obj1, obj2)
	}

	changeSet := &changeSet{BaseType: v1.Type()}
	return changeSet, doDiff(v1.Type(), v1, v2, changeSet, []PathElement{})
}

func doDiff(currType reflect.Type, v1 reflect.Value, v2 reflect.Value, cs *changeSet, ctx []PathElement) error {

	switch currType.Kind() {
	case reflect.Struct:
		for f := 0; f < currType.NumField(); f++ {
			currField := currType.Field(f)
			newCtx := extendContext(ctx, NewNameElem(f, currField.Name))
			err := doDiff(currField.Type, v1.Field(f), v2.Field(f), cs, newCtx)
			if err != nil {
				return err
			}
		}
	case reflect.Map:
		for _, key := range v1.MapKeys() {
			val2 := v2.MapIndex(key)
			newCtx := extendContext(ctx, NewKeyElem(key))
			if val2.IsValid() {
				err := doDiff(currType.Elem(), v1.MapIndex(key), v2.MapIndex(key), cs, newCtx)
				if err != nil {
					return err
				}
			} else {
				cs.AddPathDelete(newCtx)
			}
		}

		for _, key := range v2.MapKeys() {
			val1 := v1.MapIndex(key)
			if !val1.IsValid() {
				newCtx := extendContext(ctx, NewKeyElem(key))
				cs.AddPathValue(newCtx, v2.MapIndex(key))
			}
		}
	case reflect.Array:
		for i := 0; i < currType.Len(); i++ {
			newCtx := extendContext(ctx, NewIndexElem(i))
			err := doDiff(currType.Elem(), v1.Index(i), v2.Index(i), cs, newCtx)
			if err != nil {
				return err
			}
		}
	case reflect.Slice:
		minLen := intMin(v1.Len(), v2.Len())
		maxLen := intMax(v1.Len(), v2.Len())
		for i := 0; i < minLen; i++ {
			newCtx := extendContext(ctx, NewIndexElem(i))
			err := doDiff(currType.Elem(), v1.Index(i), v2.Index(i), cs, newCtx)
			if err != nil {
				return err
			}
		}

		if minLen != maxLen {
			if maxLen == v1.Len() {
				for i := minLen; i < maxLen; i++ {
					newCtx := extendContext(ctx, NewIndexElem(i))
					cs.AddPathDelete(newCtx)
				}
			} else { // maxLen == v2.Len()
				for i := minLen; i < maxLen; i++ {
					newCtx := extendContext(ctx, NewIndexElem(i))
					cs.AddPathValue(newCtx, v2.Index(i))
				}

			}
		}
	case reflect.Ptr:
		newCtx := extendContext(ctx, NewPtrElem())
		if v1.IsNil() {
			cs.AddPathValue(newCtx, v2.Elem())
		} else if v2.IsNil() {
			cs.AddPathDelete(newCtx)
		} else {
			err := doDiff(currType.Elem(), v1.Elem(), v2.Elem(), cs, newCtx)
			if err != nil {
				return err
			}
		}
	default:
		return compareBasicType(currType, v1, v2, cs, ctx)
	}

	return nil
}

func extendContext(ctx []PathElement, pe PathElement) []PathElement {
	newCtx := make([]PathElement, len(ctx), len(ctx) + 1)
	copy(newCtx, ctx)
	return append(newCtx, pe)
}

func intMin(x int, y int) int {
	if x < y {
		return x
	}

	return y
}
func intMax(x int, y int) int {
	if x > y {
		return x
	}

	return y
}

func compareBasicType(currType reflect.Type, v1 reflect.Value, v2 reflect.Value, cs *changeSet, ctx []PathElement) error {
	switch currType.Kind() {
	case reflect.String:
		if v1.String() != v2.String() {
			cs.AddPathValue(ctx, v2)
		}
	case reflect.Int64:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int:
		if v1.Int() != v2.Int() {
			cs.AddPathValue(ctx, v2)
		}

	case reflect.Uint64:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint:
		if v1.Uint() != v2.Uint() {
			cs.AddPathValue(ctx, v2)
		}

	case reflect.Float64:
		fallthrough
	case reflect.Float32:
		if v1.Float() != v2.Float() {
			cs.AddPathValue(ctx, v2)
		}

	case reflect.Complex128:
		fallthrough
	case reflect.Complex64:
		if v1.Complex() != v2.Complex() {
			cs.AddPathValue(ctx, v2)
		}

	case reflect.Bool:
		if v1.Bool() != v2.Bool() {
			cs.AddPathValue(ctx, v2)
		}

	default:
		return fmt.Errorf("unhandled kind '%v'\n", currType.Kind())
	}

	return nil
}
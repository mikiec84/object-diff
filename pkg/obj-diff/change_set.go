package obj_diff

import (
	"fmt"
	"reflect"
)

type changeSet struct {
	BaseType reflect.Type
	Changes  []Change
}

func (cs changeSet) String() string {
	return fmt.Sprintf("BaseType: %v Changes: %v", cs.BaseType, cs.Changes)
}

func (cs *changeSet) AddPathValue(ctx []PathElement, v reflect.Value) {
	cs.Changes = append(cs.Changes, Change{Path: ctx, NewValue: v})
}

func (cs *changeSet) AddPathDelete(ctx []PathElement) {
	cs.Changes = append(cs.Changes, Change{Path: ctx, Deleted: true})
}

func (cs changeSet) Patch(obj interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			e, ok := r.(PatchError)
			if ok {
				fmt.Println("Recovered in Patch", r)
				err = e
			} else {
				panic(r)
			}
		}
	}()

	root := reflect.ValueOf(obj)
	if root.Kind() != reflect.Ptr || !root.Elem().CanSet() {
		return NewPatchError("can not set obj1 of Type: %v", root.Type())
	}

	if root.Elem().Type() != cs.BaseType {
		return NewPatchError("obj1 (%T) is not of type %v", obj, cs.BaseType)
	}

	opConfig := ObjectPathConfig{true, true}

	for i := 0; i < len(cs.Changes); i++ {
		fmt.Println()
		fmt.Printf("##################\n")
		fmt.Printf("### New Change ###\n")
		fmt.Printf("##################\n")
		// fmt.Printf("Root: %+v\n", root)
		change := cs.Changes[i]
		fmt.Printf("Change: %+v\n", change)
		op := NewObjectPathWithConfig(root, change.Path, opConfig)
		for op.Next() {
			fmt.Printf("Types lastVal: %T, currVal: %T\n", op.LastVal().Interface(), op.Interface())
			fmt.Printf("Kinds lastVal: %v, currVal: %v\n", op.LastVal().Kind(), op.Kind())
			fmt.Println("==================")

			switch op.Kind() {
			case reflect.Struct:
				// NO-OP
			case reflect.Map:
				// NO-OP
			case reflect.Array:
				// NO-OP
			case reflect.Slice:
				// NO-OP
			case reflect.Ptr:
				// NO-OP
			}
		}

		if change.Deleted {
			op.Delete()
		} else {
			op.Set(change.NewValue)
		}
	}

	return nil
}

func (cs changeSet) Equals(other changeSet) bool {
	if cs.BaseType !=  other.BaseType {
		return false
	}

	if len(cs.Changes) != len(other.Changes) {
		return false
	}

	for i := 0; i < len(cs.Changes); i++ {
		c1 := cs.Changes[i]
		c2 := other.Changes[i]

		if !c1.Equals(c2) {
			return false
		}
	}

	return true
}
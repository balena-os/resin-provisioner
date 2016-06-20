package pinejs

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/resin-io/pinejs-client-go/Godeps/_workspace/src/github.com/fatih/structs"
)

// Retrieve resource name from input struct - if contains a pinejs tag, use
// that, otherwise use the lowercase of the struct name.
func resourceNameFromStruct(v interface{}) string {
	// Only called from a function that asserts input is a struct.

	// Look for pinejs tag, use it if we find it.
	for _, f := range structs.Fields(v) {
		if name := f.Tag("pinejs"); name != "" {
			return name
		}
	}

	// Otherwise, we default to the name of the struct in lower case.
	return strings.ToLower(structs.Name(v))
}

func resourceNameFromMap(v interface{}) (string, error) {
	if m, ok := v.(map[string]interface{}); !ok {
		return "", errors.New("Invalid map type")
	} else if val, ok := m["pinejs"].(string); ok && val != "" {
		return val, nil
	} else {
		return "", errors.New("Failed to retrieve resource name from map")
	}
}

// Unwinds pointers, slices, and slices of pointers, etc. until we get to a
// struct then we hand off to resourceNameFromStruct, or the equivalent with a map.
func resourceName(v interface{}) (string, error) {
	ty := reflect.TypeOf(v)

	switch ty.Kind() {
	case reflect.Struct:
		return resourceNameFromStruct(v), nil
	case reflect.Map:
		return resourceNameFromMap(v)
	case reflect.Ptr:
		return resourceName(reflect.Indirect(reflect.ValueOf(v)).Interface())
	case reflect.Slice:
		if ty.Elem().Kind() == reflect.Struct {
			// Create new pointer to pointer/slice type.
			ptr := reflect.New(ty.Elem())
			// Deref the pointer and recurse on that value until we get to a struct.
			el := ptr.Elem().Interface()
			return resourceName(el)
		} else if ty.Elem().Kind() == reflect.Map {
			theSlice := reflect.ValueOf(v)
			if theSlice.Len() == 0 {
				return "", errors.New("Slice of maps must have non-zero length")
			} else {
				return resourceName(theSlice.Index(0).Interface())
			}
		}
	}

	return "", fmt.Errorf("tried to retrieve resource name from non-struct and non-map %s",
		ty.Kind())
}

func getResourceField(v interface{}) (f *structs.Field, err error) {
	var ok bool

	if !structs.IsStruct(v) {
		err = errors.New("not a struct")
	} else if f, ok = structs.New(v).FieldOk("Id"); !ok {
		err = errors.New("no 'Id' field")
	} else if _, ok = f.Value().(int); !ok {
		err = errors.New("Id field is not an int")
	}

	return
}

func getIdFromMap(m map[string]interface{}) (ret int, err error) {
	if val, ok := m["id"].(int); !ok {
		ret = 0
	} else {
		ret = val
	}
	return
}

// Retrieve Id field from interface.
func resourceId(v interface{}) (ret int, err error) {
	var f *structs.Field
	invalidTypeErr := errors.New("Not a struct, map, or pointer to either")
	ty := reflect.TypeOf(v)

	switch ty.Kind() {
	case reflect.Map:
		if m, ok := v.(map[string]interface{}); !ok {
			err = errors.New("Invalid map")
		} else {
			return getIdFromMap(m)
		}
	case reflect.Struct:
		if f, err = getResourceField(v); err == nil {
			ret = f.Value().(int)
		}
	case reflect.Ptr:
		return resourceId(reflect.Indirect(reflect.ValueOf(v)).Interface())
	default:
		err = invalidTypeErr
	}

	return
}

// Determine whether the struct's id will be omitted in json encoding.
func isIdOmitted(v interface{}) (bool, error) {

	if structs.IsStruct(v) {
		if f, err := getResourceField(v); err != nil {
			return false, err
		} else if jsonTag := f.Tag("json"); jsonTag == "" {
			// No json tag means the id field won't be ommitted.
			return false, nil
		} else {
			// getResourceField() ensures this is an int.
			id := f.Value().(int)
			// Json tags are comma separated. 'omitempty' means id == 0 -> id field
			// not included in generated json. Also note spacing, like "id,
			// omitempty" is significant and prevents a tag from taking effect so no
			// need for trimming.
			//
			// See http://golang.org/pkg/encoding/json/#Marshal
			for _, field := range strings.Split(jsonTag, ",") {
				if field == "omitempty" && id == 0 {
					return true, nil
				}
			}
		}
	} else {
		ty := reflect.TypeOf(v)
		switch ty.Kind() {
		case reflect.Map:
			if _, ok := v.(map[string]interface{}); !ok {
				return false, errors.New("Invalid map")
			} else {
				return false, nil
			}
		case reflect.Ptr:
			return isIdOmitted(reflect.Indirect(reflect.ValueOf(v)).Interface())
		}
	}

	// Id field exists, no 'omitempty' tag so not omitted.
	return false, nil
}

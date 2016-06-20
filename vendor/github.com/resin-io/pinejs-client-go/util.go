package pinejs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"strings"

	"github.com/resin-io/pinejs-client-go/Godeps/_workspace/src/github.com/bitly/go-simplejson"
)

// oDataEncodeVals URL Encode values and separate with commas.
func oDataEncodeVals(strs []string) string {
	encoded := make([]string, len(strs))

	for i, str := range strs {
		encoded[i] = strings.Replace(url.QueryEscape(str), "+", "%20", -1)
	}

	return strings.Join(encoded, ",")
}

// encodeQuery encodes query values, working around a net/url issue whereby keys
// get encoded as well as values. We only want values encoded, otherwise OData
// dies.
func encodeQuery(query map[string][]string) string {
	if len(query) == 0 {
		return ""
	}

	var tuples []string
	for key, vals := range query {
		if len(vals) == 0 {
			continue
		}

		tuple := key + "=" + oDataEncodeVals(vals)
		tuples = append(tuples, tuple)
	}

	return "?" + strings.Join(tuples, "&")
}

// getSinglePath generates a path for accessing a specific resource and ID.
func getSinglePath(v interface{}) (string, error) {
	if name, err := resourceName(v); err != nil {
		return "", err
	} else if id, err := resourceId(v); err != nil {
		return "", err
	} else {
		return fmt.Sprintf("%s(%d)", name, id), nil
	}
}

// ptrType extracts the pointer type of a provided interface, e.g. *Foo -> Foo.
func ptrType(v interface{}) (reflect.Type, error) {
	if ty := reflect.TypeOf(v); ty.Kind() != reflect.Ptr {
		return nil, errors.New("not a pointer")
	} else {
		return ty.Elem(), nil
	}
}

func isPointerToStruct(v interface{}) (bool, error) {
	if el, err := ptrType(v); err != nil {
		return false, err
	} else if el.Kind() != reflect.Struct {
		return false, errors.New("not a pointer to a struct")
	}

	return true, nil
}

func isPointerToStructOrMap(v interface{}) (bool, error) {
	if el, err := ptrType(v); err != nil {
		return false, err
	} else if el.Kind() != reflect.Struct && el.Kind() != reflect.Map {
		return false, errors.New("not a pointer to a struct or map")
	}

	return true, nil
}

func isPointerToSliceStructsOrMaps(v interface{}) (bool, error) {
	if el, err := ptrType(v); err != nil {
		return false, err
	} else if el.Kind() != reflect.Slice {
		return false, errors.New("not a pointer to a slice")
	} else if el.Elem().Kind() != reflect.Struct && el.Elem().Kind() != reflect.Map {
		return false, errors.New("not a pointer to a slice of structs or maps")
	}

	return true, nil
}

func toJsonReader(v interface{}) (io.Reader, error) {
	if v == nil {
		return nil, nil
	}

	if buf, err := json.Marshal(v); err != nil {
		return nil, err
	} else {
		return bytes.NewReader(buf), nil
	}
}

// Some functionality that is strangely lacking from simplejson...

type jsonNodeType int

const (
	jsonObject jsonNodeType = iota
	jsonArray
	jsonValue // Anything else.
)

func getJsonNodeType(j *simplejson.Json) jsonNodeType {
	// TODO: Reuse returned values.
	if _, err := j.Map(); err == nil {
		return jsonObject
	} else if _, err := j.Array(); err == nil {
		return jsonArray
	} else {
		return jsonValue
	}
}

func getJsonFieldNames(j *simplejson.Json) (ret []string, err error) {
	var obj map[string]interface{}

	if obj, err = j.Map(); err == nil {
		for name, _ := range obj {
			ret = append(ret, name)
		}
	}

	return
}

func getJsonFields(j *simplejson.Json) (map[string]*simplejson.Json, error) {
	if ns, err := getJsonFieldNames(j); err != nil {
		return nil, err
	} else {
		ret := make(map[string]*simplejson.Json)

		for _, name := range ns {
			ret[name] = j.Get(name)
		}

		return ret, nil
	}
}

func getJsonArray(j *simplejson.Json) (ret []*simplejson.Json, err error) {
	var arr []interface{}

	if arr, err = j.Array(); err == nil {
		// TODO: This sucks. Don't want to remarshal just to use returned data
		// though.
		for i := 0; i < len(arr); i++ {
			ret = append(ret, j.GetIndex(i))
		}
	}

	return
}

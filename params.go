package httpc

import (
	"fmt"
	"net/url"
	"reflect"
)

const (
	TAG_NAME_JSON    = "json"
	TAG_VALUE_IGNORE = "-" //ignore
)

//make query params from struct and json tag only
func MakeQueryParams(v interface{}) url.Values {

	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)

LOOP:
	for {
		kind := typ.Kind()
		switch kind {
		case reflect.Ptr:
			typ = typ.Elem()
			val = val.Elem()
			continue
		default:
			break LOOP
		}
	}
	return parseStructValues(typ, val)
}

// get struct field's tag value
func getTag(sf reflect.StructField, tagName string) (strValue string, ignore bool) {

	strValue = sf.Tag.Get(tagName)
	if strValue == TAG_VALUE_IGNORE {
		return "", true
	}
	return
}

// parse struct fields
func parseStructValues(typ reflect.Type, val reflect.Value) url.Values {
	var values = make(url.Values)
	kind := typ.Kind()
	if kind == reflect.Struct {
		NumField := val.NumField()
		for i := 0; i < NumField; i++ {
			typField := typ.Field(i)
			valField := val.Field(i)

			if typField.Type.Kind() == reflect.Ptr {
				typField.Type = typField.Type.Elem()
				valField = valField.Elem()
			}
			if !valField.IsValid() || !valField.CanInterface() {
				continue
			}
			strTagVal, ignore := getTag(typField, TAG_NAME_JSON)
			if ignore {
				continue
			}
			values[strTagVal] = []string{fmt.Sprintf("%v", valField.Interface())}
		}
	}
	return values
}

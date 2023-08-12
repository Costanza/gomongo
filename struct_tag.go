package gomongo

import (
	"reflect"
	"strings"
)

func FieldToJSONTag(entity interface{}, field string) (tag string) {
	st := reflect.TypeOf(entity)
	f, found := st.FieldByName(field)
	if found {
		t := f.Tag.Get("bson")

		parts := strings.Split(t, ",")
		if len(parts) > 0 {
			tag = parts[0]
		}
	}

	return
}

func FieldToBSONTag(entity interface{}, field string) (tag string) {
	st := reflect.TypeOf(entity)
	f, found := st.FieldByName(field)
	if found {
		t := f.Tag.Get("bson")

		parts := strings.Split(t, ",")
		if len(parts) > 0 {
			tag = parts[0]
		}
	}

	return
}

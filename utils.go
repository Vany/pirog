package pirog

import (
	"bytes"
	"encoding/json"
	"reflect"
)

func ToJson(in any) string {
	buff := bytes.Buffer{}
	MUST(json.NewEncoder(&buff).Encode(in))
	return buff.String()
}

// PutToStruct search for corresponding type field in structure and put obj there
func PutToStruct(c any, obj any) {
	ft := reflect.TypeOf(c).Elem()
	ot := reflect.TypeOf(obj)
	for i := 0; i < ft.NumField(); i++ {
		ftc := ft.Field(i).Type
		if ftc.Kind() == reflect.Pointer {
			ftc = ftc.Elem()
		}
		if ot.Elem().AssignableTo(ftc) || (ftc.Kind() == reflect.Interface && ot.Implements(ftc)) {
			reflect.ValueOf(c).Elem().Field(i).Set(reflect.ValueOf(obj))
		}
	}
}

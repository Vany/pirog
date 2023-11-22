package pirog

import (
	"bytes"
	"context"
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

// ExecuteOnAllFields - On all interface fields run method by name
func ExecuteOnAllFields(ctx context.Context, a any, mname string) error {
	v := reflect.ValueOf(a).Elem()
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).IsNil() {
			continue
		}
		f := v.Field(i).Elem()
		m := f.MethodByName(mname)
		if !m.IsValid() {
			continue
		}
		ret := m.Call([]reflect.Value{reflect.ValueOf(ctx)})[0]
		if !ret.IsNil() {
			return ret.Interface().(error)
		}
	}
	return nil
}

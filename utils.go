package pirog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/sync/errgroup"
	. "reflect"
)

func ToJson(in any) string {
	buff := bytes.Buffer{}
	MUST(json.NewEncoder(&buff).Encode(in))
	return buff.String()
}

// ExecuteOnAllFields - On all interface fields run method by name
func ExecuteOnAllFields(ctx context.Context, a any, mname string) error {
	v := ValueOf(a).Elem()
	wg := errgroup.Group{}

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Type().Kind() != Pointer && f.Type().Kind() != Interface || f.IsNil() {
			continue
		}
		if f.Type().Kind() == Interface {
			f = f.Elem()
		}
		m := f.MethodByName(mname)
		if !m.IsValid() {
			continue
		}
		wg.Go(func(i int) func() error {
			return func() error {
				ret := m.Call([]Value{ValueOf(ctx)})[0]
				if !ret.IsNil() {
					return fmt.Errorf("%s => %w", v.Type().Field(i).Name, ret.Interface().(error))
				}
				return nil
			}
		}(i))
	}
	return wg.Wait()
}

// InjectComponents - search for corresponding fields in fields and put references there
func InjectComponents(a any) {
	v := ValueOf(a).Elem()
	fields := make(map[string]Value)
	for i := 0; i < v.NumField(); i++ {
		t := v.Type().Field(i)
		tag, have := t.Tag.Lookup("injectable")
		if !have || !t.IsExported() || t.Type.Kind() != Pointer && t.Type.Kind() != Interface {
			continue
		}
		fname := COALESCE(tag, t.Name)
		if _, ok := fields[fname]; ok {
			panic("Injecting table already contains field '" + fname + "'")
		}
		fields[fname] = v.Field(i)
	}

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		t := v.Type().Field(i)
		if !t.IsExported() || t.Type.Kind() != Pointer && t.Type.Kind() != Interface || f.IsNil() {
			continue
		}
		f = f.Elem()
		if t.Type.Kind() == Interface {
			f = f.Elem()
		}

		for j := 0; j < f.NumField(); j++ {
			ff := f.Field(j)
			tt := f.Type().Field(j)
			tag, have := tt.Tag.Lookup("inject")
			if !have || !tt.IsExported() {
				continue
			}
			tag = COALESCE(tag, tt.Name)
			if vv, ok := fields[tag]; !ok && ff.IsZero() {
				panic(`want inject unexistent field name into empty field"` + tag + `" into "` + t.Name + "." + tt.Name + `"`)
			} else if !ok {
				if DEBUG {
					println("> ", t.Name+"."+tt.Name, "<< Already have a value and have no candidates")
				}
			} else {
				if DEBUG {
					println("> ", t.Name+"."+tt.Name, "<<", tag, vv.Type().Name())
				}
				if vv.CanConvert(ff.Type()) {
					ff.Set(vv)
				} else {
					ff.Set(vv.Elem())
				}
			}
		}
	}
}

package copier

import (
	"fmt"
	"reflect"
	"strings"
)

// ToMap ...
func ToMap(obj interface{}) (map[string]interface{}, error) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	t := v.Type()
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("only instance of struct and pointer to struct can be converted to map (has %s)", t)
	}
	m := make(map[string]interface{}, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		js := f.Tag.Get("json")
		if js == "-" {
			continue
		}
		fn := f.Name
		if n := strings.Split(js, ",")[0]; n != "" {
			fn = n
		}
		v2 := v.FieldByIndex(f.Index).Interface()
		if f.Type.Kind() == reflect.Struct || f.Type.Kind() == reflect.Pointer && f.Type.Elem().Kind() == reflect.Struct {
			m2, err := ToMap(v2)
			if err != nil {
				return nil, err
			}
			m[fn] = m2
		} else {
			m[fn] = v2
		}
	}
	return m, nil
}

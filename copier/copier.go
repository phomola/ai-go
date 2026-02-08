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
		return nil, fmt.Errorf("only instances of struct and pointer to struct can be converted to map (has %s)", t)
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

// FromMap ...
func FromMap[T any](m map[string]interface{}) (*T, error) {
	obj := new(T)
	if err := fromMap(m, reflect.ValueOf(obj).Elem()); err != nil {
		return nil, err
	}
	return obj, nil
}

func fromMap(m map[string]interface{}, v reflect.Value) error {
	t := v.Type()
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
		if f.Type.Kind() == reflect.Struct {
			m2, ok := m[fn].(map[string]interface{})
			if !ok {
				return fmt.Errorf("expected map for '%s'", fn)
			}
			v2 := v.FieldByIndex(f.Index)
			if err := fromMap(m2, v2); err != nil {
				return err
			}
		} else if f.Type.Kind() == reflect.Pointer && f.Type.Elem().Kind() == reflect.Struct {
			m2, ok := m[fn].(map[string]interface{})
			if !ok {
				return fmt.Errorf("expected map for '%s'", fn)
			}
			v2 := reflect.New(f.Type.Elem())
			if err := fromMap(m2, v2.Elem()); err != nil {
				return err
			}
			v.FieldByIndex(f.Index).Set(v2)
		} else {
			switch f.Type.Kind() {
			case reflect.Int:
				switch x := m[fn].(type) {
				case int:
					v.FieldByIndex(f.Index).SetInt(int64(x))
				case float32:
					v.FieldByIndex(f.Index).SetInt(int64(x))
				case float64:
					v.FieldByIndex(f.Index).SetInt(int64(x))
				default:
					return fmt.Errorf("expected number for '%s'", fn)
				}
			case reflect.Float32, reflect.Float64:
				switch x := m[fn].(type) {
				case int:
					v.FieldByIndex(f.Index).SetFloat(float64(x))
				case float32:
					v.FieldByIndex(f.Index).SetFloat(float64(x))
				case float64:
					v.FieldByIndex(f.Index).SetFloat(x)
				default:
					return fmt.Errorf("expected number for '%s'", fn)
				}
			case reflect.String:
				x, ok := m[fn].(string)
				if !ok {
					return fmt.Errorf("expected string for '%s'", fn)
				}
				v.FieldByIndex(f.Index).SetString(x)
			case reflect.Bool:
				x, ok := m[fn].(bool)
				if !ok {
					return fmt.Errorf("expected boolean for '%s'", fn)
				}
				v.FieldByIndex(f.Index).SetBool(x)
			}
		}
	}
	return nil
}

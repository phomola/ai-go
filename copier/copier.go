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
		v2, err := toMapValue(v.FieldByIndex(f.Index))
		if err != nil {
			return nil, err
		}
		m[fn] = v2
	}
	return m, nil
}

func toMapValue(v reflect.Value) (interface{}, error) {
	t := v.Type()
	if t.Kind() == reflect.Struct || t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Struct {
		return ToMap(v.Interface())
	}
	if t.Kind() == reflect.Slice {
		sl := make([]interface{}, 0, v.Len())
		for i := 0; i < v.Len(); i++ {
			v2, err := toMapValue(v.Index(i))
			if err != nil {
				return nil, err
			}
			sl = append(sl, v2)
		}
		return sl, nil
	}
	return v.Interface(), nil
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
		x, ok := m[fn]
		if !ok {
			continue
		}
		if err := fromMapValue(x, v.FieldByIndex(f.Index)); err != nil {
			return err
		}
	}
	return nil
}

func fromMapValue(x interface{}, v reflect.Value) error {
	t := v.Type()
	if t.Kind() == reflect.Struct {
		m, ok := x.(map[string]interface{})
		if !ok {
			return fmt.Errorf("expected map instead of '%s'", x)
		}
		if err := fromMap(m, v); err != nil {
			return err
		}
	} else if t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Struct {
		m, ok := x.(map[string]interface{})
		if !ok {
			return fmt.Errorf("expected map instead of '%s'", x)
		}
		v2 := reflect.New(t.Elem())
		if err := fromMap(m, v2.Elem()); err != nil {
			return err
		}
		v.Set(v2)
	} else if t.Kind() == reflect.Slice {
		sl, ok := x.([]interface{})
		if !ok {
			return fmt.Errorf("expected slice instead of '%s'", x)
		}
		v2 := reflect.MakeSlice(t, len(sl), len(sl))
		for i, x := range sl {
			if err := fromMapValue(x, v2.Index(i)); err != nil {
				return err
			}
		}
		v.Set(v2)
	} else {
		switch t.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			switch x2 := x.(type) {
			case int:
				v.SetInt(int64(x2))
			case int8:
				v.SetInt(int64(x2))
			case int16:
				v.SetInt(int64(x2))
			case int32:
				v.SetInt(int64(x2))
			case int64:
				v.SetInt(int64(x2))
			case float32:
				v.SetInt(int64(x2))
			case float64:
				v.SetInt(int64(x2))
			default:
				return fmt.Errorf("expected number instead of '%s'", x)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			switch x2 := x.(type) {
			case uint:
				v.SetUint(uint64(x2))
			case uint8:
				v.SetUint(uint64(x2))
			case uint16:
				v.SetUint(uint64(x2))
			case uint32:
				v.SetUint(uint64(x2))
			case uint64:
				v.SetUint(uint64(x2))
			default:
				return fmt.Errorf("expected number instead of '%s'", x)
			}
		case reflect.Float32, reflect.Float64:
			switch x2 := x.(type) {
			case int:
				v.SetFloat(float64(x2))
			case int8:
				v.SetFloat(float64(x2))
			case int16:
				v.SetFloat(float64(x2))
			case int32:
				v.SetFloat(float64(x2))
			case int64:
				v.SetFloat(float64(x2))
			case float32:
				v.SetFloat(float64(x2))
			case float64:
				v.SetFloat(x2)
			default:
				return fmt.Errorf("expected number instead of '%s'", x)
			}
		case reflect.String:
			x2, ok := x.(string)
			if !ok {
				return fmt.Errorf("expected string instead of '%s'", x)
			}
			v.SetString(x2)
		case reflect.Bool:
			x2, ok := x.(bool)
			if !ok {
				return fmt.Errorf("expected boolean instead of '%s'", x)
			}
			v.SetBool(x2)
		default:
			return fmt.Errorf("unhandled type '%s'", t)
		}
	}
	return nil
}

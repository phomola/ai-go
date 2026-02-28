package infer

import (
	"context"
	"fmt"
	"reflect"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/phomola/ai-go/copier"
)

var (
	contextType = reflect.TypeFor[context.Context]()
	errorType   = reflect.TypeFor[error]()
)

// Function ...
type Function struct {
	Name        string
	Description string
	InSchema    *jsonschema.Schema
	OutSchema   *jsonschema.Schema
	Fn          func(context.Context, map[string]any) (map[string]any, error)
}

// Functions ...
func Functions[T any](obj *T) ([]*Function, error) {
	typ := reflect.TypeFor[*T]()
	funcs := make([]*Function, 0, typ.NumMethod())
	objPtr := reflect.ValueOf(obj)
	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)
		if m.Type.NumIn() != 4 {
			return nil, fmt.Errorf("method '%s' must have 3 arguments", m.Name)
		}
		if m.Type.NumOut() != 2 {
			return nil, fmt.Errorf("method '%s' must have 2 return values", m.Name)
		}
		if m.Type.In(1) != contextType {
			return nil, fmt.Errorf("the first argument of method '%s' must be a context", m.Name)
		}
		if m.Type.In(2).Kind() != reflect.Pointer || m.Type.In(2).Elem().Kind() != reflect.Struct {
			return nil, fmt.Errorf("the second argument of method '%s' must be a pointer to a struct", m.Name)
		}
		if m.Type.In(3).Kind() != reflect.Pointer || m.Type.In(3).Elem().Kind() != reflect.Struct {
			return nil, fmt.Errorf("the third argument of method '%s' must be a pointer to a struct", m.Name)
		}
		if m.Type.Out(0).Kind() != reflect.Pointer || m.Type.Out(0).Elem().Kind() != reflect.Struct {
			return nil, fmt.Errorf("the first return value of method '%s' must be a pointer to a struct", m.Name)
		}
		if m.Type.Out(1) != errorType {
			return nil, fmt.Errorf("the second return value of method '%s' must be an erro", m.Name)
		}
		infoType := m.Type.In(3)
		var methodDesc string
		f, ok := infoType.Elem().FieldByName("Info")
		if ok {
			methodDesc = f.Tag.Get("guide")
		}
		inType := m.Type.In(2).Elem()
		outType := m.Type.Out(0).Elem()
		inSchema, err := jsonschema.ForType(inType, nil)
		if err != nil {
			return nil, err
		}
		outSchema, err := jsonschema.ForType(outType, nil)
		if err != nil {
			return nil, err
		}
		funcs = append(funcs, &Function{
			Name:        m.Name,
			Description: methodDesc,
			InSchema:    inSchema,
			OutSchema:   outSchema,
			Fn: func(ctx context.Context, inMap map[string]any) (map[string]any, error) {
				in, err := copier.FromMapAny(inMap, inType)
				if err != nil {
					return nil, err
				}
				ret := m.Func.Call([]reflect.Value{objPtr, reflect.ValueOf(ctx), reflect.ValueOf(in), reflect.Zero(infoType)})
				if err := ret[1].Interface(); err != nil {
					return nil, err.(error)
				}
				return copier.ToMap(ret[0].Interface())
			},
		})
	}
	return funcs, nil
}

package infer

import (
	"context"
	"fmt"
	"reflect"

	"github.com/google/jsonschema-go/jsonschema"
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
	Fn          func(context.Context, any) (any, error)
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

		styp := m.Type.In(3)
		var methodDesc string
		f, ok := styp.FieldByName("Info")
		if ok {
			methodDesc = f.Tag.Get("guide")
		}
		funcs = append(funcs, &Function{
			Name:        m.Name,
			Description: methodDesc,
			Fn: func(ctx context.Context, in any) (any, error) {
				ret := m.Func.Call([]reflect.Value{objPtr, reflect.ValueOf(ctx), reflect.ValueOf(in), reflect.New(styp).Elem()})
				return ret[0].Interface(), ret[1].Interface().(error)
			},
		})
	}
	return funcs, nil
}

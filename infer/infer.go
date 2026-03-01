package infer

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/phomola/ai-go/copier"
	"github.com/phomola/ai-go/gemini/ai"
)

var (
	contextType = reflect.TypeFor[context.Context]()
	errorType   = reflect.TypeFor[error]()
)

// GeminiTool ...
func GeminiTool(funcs []*Function) (*ai.Tool, error) {
	var tool ai.Tool
	for _, f := range funcs {
		if err := tool.AddFunction(f.Name, f.FullDescription(), f.InSchema, f.OutSchema, f.Fn); err != nil {
			return nil, err
		}
	}
	return &tool, nil
}

// Argument ...
type Argument struct {
	Name  string
	Guide string
}

// Function ...
type Function struct {
	Name        string
	Description string
	Arguments   []Argument
	InSchema    *jsonschema.Schema
	OutSchema   *jsonschema.Schema
	Fn          func(context.Context, map[string]any) (map[string]any, error)
}

// FullDescription ...
func (f *Function) FullDescription() string {
	var sb strings.Builder
	if f.Description != "" {
		sb.WriteString(f.Description)
		sb.WriteString("\n\n")
	}
	sb.WriteString("Arguments:\n")
	for i, arg := range f.Arguments {
		sb.WriteString(arg.Name)
		sb.WriteString(": ")
		sb.WriteString(arg.Guide)
		if i+1 < len(f.Arguments) {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

// Functions ...
func Functions[T any](obj *T) ([]*Function, error) {
	typ := reflect.TypeFor[*T]()
	typName := typ.Elem().Name()
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
		args := make([]Argument, 0, inType.NumField())
		for j := 0; j < inType.NumField(); j++ {
			field := inType.Field(j)
			args = append(args, Argument{
				Name:  field.Name,
				Guide: field.Tag.Get("jsonschema"),
			})
		}
		inSchema, err := jsonschema.ForType(inType, nil)
		if err != nil {
			return nil, err
		}
		outSchema, err := jsonschema.ForType(outType, nil)
		if err != nil {
			return nil, err
		}
		funcs = append(funcs, &Function{
			Name:        typName + "::" + m.Name,
			Description: methodDesc,
			Arguments:   args,
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

package ai

import (
	"github.com/phomola/ai-go/copier"
	"google.golang.org/genai"
)

// Tool is an LLM tool with functions.
type Tool struct {
	funcDecls []*genai.FunctionDeclaration
	functions map[string]func(map[string]any) (map[string]any, error)
}

// AddFunction adds a function to a tool.
func AddFunction[I, O any](tool *Tool, name, description string, f func(*I) (*O, error)) error {
	inSchema, err := schemaFor[I]()
	if err != nil {
		return err
	}
	outSchema, err := schemaFor[O]()
	if err != nil {
		return err
	}
	tool.funcDecls = append(tool.funcDecls, &genai.FunctionDeclaration{
		Name:                 name,
		Description:          description,
		ParametersJsonSchema: inSchema,
		ResponseJsonSchema:   outSchema,
	})
	if tool.functions == nil {
		tool.functions = make(map[string]func(map[string]any) (map[string]any, error))
	}
	tool.functions[name] = func(inMap map[string]any) (map[string]any, error) {
		in, err := copier.FromMap[I](inMap)
		if err != nil {
			return nil, err
		}
		out, err := f(in)
		if err != nil {
			return nil, err
		}
		return copier.ToMap(out)
	}
	return nil
}

func (t *Tool) tool() *genai.Tool {
	return &genai.Tool{FunctionDeclarations: t.funcDecls}
}

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
	inRs, err := inSchema.Resolve(nil)
	if err != nil {
		return err
	}
	outRs, err := outSchema.Resolve(nil)
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
		if err := inRs.Validate(inMap); err != nil {
			return nil, err
		}
		in, err := copier.FromMap[I](inMap)
		if err != nil {
			return nil, err
		}
		out, err := f(in)
		if err != nil {
			return nil, err
		}
		outMap, err := copier.ToMap(out)
		if err != nil {
			return nil, err
		}
		if err := outRs.Validate(outMap); err != nil {
			return nil, err
		}
		return outMap, nil
	}
	return nil
}

func (t *Tool) tool() *genai.Tool {
	return &genai.Tool{FunctionDeclarations: t.funcDecls}
}

package proxy

import (
	"context"
	"fmt"
	"math"

	"github.com/phomola/ai-go/gemini/ai"
	"github.com/phomola/ai-go/infer"
	"github.com/phomola/ai-go/nlp"
)

type proxyInput struct {
	Prompt string `json:"prompt" jsonschema:"The original prompt."`
}

type proxyOutput struct {
	Output string `json:"output"`
}

type function struct {
	vec nlp.Vector
	fn  *infer.Function
}

// Tool creates a proxy tool.
func Tool(functions []*infer.Function, emb nlp.Embedding) (*ai.Tool, error) {
	funcs := make([]*function, 0, len(functions))
	for _, f := range functions {
		vec, err := emb.Vector(f.Description)
		if err != nil {
			return nil, err
		}
		vec.Normalise()
		funcs = append(funcs, &function{
			vec: vec,
			fn:  f,
		})
	}
	var tool ai.Tool
	if err := ai.AddFunction(&tool, "proxyTool", "A tool for answering prompts that the LLM alone can't handle.", func(ctx context.Context, in *proxyInput) (*proxyOutput, error) {
		vec, err := emb.Vector(in.Prompt)
		if err != nil {
			return nil, err
		}
		vec.Normalise()
		var (
			fn      *infer.Function
			minDist = math.MaxFloat64
		)
		for _, f := range funcs {
			if dist := 1 - nlp.DotProd(vec, f.vec); dist < minDist {
				minDist = dist
				fn = f.fn
			}
		}
		fmt.Println("proxy picked:", fn.Name)
		var tool ai.Tool
		if err := tool.AddFunction(fn.Name, fn.Description, fn.InSchema, fn.OutSchema, fn.Fn); err != nil {
			return nil, err
		}
		cl, err := ai.NewClient(ctx, ai.Gemini3FlashPreview)
		if err != nil {
			return nil, err
		}
		resp, err := cl.GenerateText(ctx, ai.NewText(in.Prompt), []*ai.Tool{&tool})
		if err != nil {
			return nil, err
		}
		return &proxyOutput{
			Output: resp.String(),
		}, nil
	}); err != nil {
		return nil, err
	}
	return &tool, nil
}

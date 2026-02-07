package main

import (
	"context"
	"fmt"
	"log"

	"github.com/phomola/ai-go/gemini/ai"
)

func main() {
	ctx := context.Background()

	cl, err := ai.NewClient(ctx, ai.Gemini3FlashPreview)
	if err != nil {
		log.Fatal(err)
	}

	type Output struct {
		Name string
	}

	var nameTool ai.Tool
	ai.AddFunction(&nameTool, "nameTool", "Provides the current user's name.", func(_ *struct{}) (*Output, error) {
		return &Output{Name: "Sean"}, nil
	})

	resp, err := cl.GenerateText(ctx, ai.NewText("What's my name?"), []*ai.Tool{&nameTool})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)
}

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
	resp, err := cl.Generate(ctx, ai.NewText("What is an LLM?"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}

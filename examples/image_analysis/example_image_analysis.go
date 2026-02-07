package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"mime"
	"os"
	"path/filepath"

	"github.com/phomola/ai-go/gemini/ai"
)

type (
	// Menu ...
	Menu struct {
		Items []Item `json:"items" jsonschema:"The list of items on the menu."`
	}

	// Item ...
	Item struct {
		Name              string  `json:"name"`
		NameInEnglish     string  `json:"nameInEnglish" jsonschema:"The name of the item in English."`
		Price             float64 `json:"price"`
		Currency          string  `json:"currency,omitempty"`
		Category          string  `json:"category,omitempty"`
		CategoryInEnglish string  `json:"categoryInEnglish,omitempty" jsonschema:"The category of the item in English."`
		Portion           string  `json:"portion,omitempty" jsonschema:"The size or weight of the item."`
	}
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatal("a file name is required")
	}
	fileName := flag.Arg(0)

	ext := filepath.Ext(fileName)
	if ext == "" {
		log.Fatal("no file extension")
	}
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		log.Fatal("can't infer mime type")
	}

	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	cl, err := ai.NewClient(ctx, ai.Gemini3FlashPreview)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := ai.Generate[Menu](ctx, cl,
		ai.NewTextWithImage("Extract all the items from the menu. Include the weight in the 'portion' property.", b, mimeType), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", resp)
}

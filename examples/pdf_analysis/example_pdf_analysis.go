package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/phomola/ai-go/gemini/ai"
)

type (
	// CV ...
	CV struct {
		Name     string `json:"name"`
		Email    string `json:"email,omitempty"`
		Phone    string `json:"phone,omitempty"`
		Location string `json:"location,omitempty"`
		Summary  string `json:"summary,omitempty"`
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

	mimeType2 := http.DetectContentType(b)
	fmt.Println(mimeType, mimeType2)

	ctx := context.Background()

	cl, err := ai.NewClient(ctx, ai.Gemini3ProPreview)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := ai.Generate[CV](ctx, cl,
		ai.NewTextWithBytes("Extract relevant information for the CV file.", b, mimeType), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", resp)
}

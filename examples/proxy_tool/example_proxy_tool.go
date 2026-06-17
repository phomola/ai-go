package main

/*
#cgo LDFLAGS: -L. -lembedding

#include <stdlib.h>
#include <string.h>

typedef struct {
    double* data;
    int size;
} nl_vector;

void* nl_sentence_embedding_for_english();
void nl_embedding_release(void* embedding);
nl_vector nl_get_vector(void* embedding, char* text);
*/
import "C"

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime"
	"unsafe"

	"github.com/phomola/ai-go/gemini/ai"
	"github.com/phomola/ai-go/gemini/proxy"
	"github.com/phomola/ai-go/infer"
	"github.com/phomola/ai-go/nlp"
)

const (
	sizeofFloat64 = C.size_t(unsafe.Sizeof(float64(0.0)))
)

var (
	_ nlp.Embedding = (*Embedding)(nil)
)

// Embedding ...
type Embedding struct {
	natptr unsafe.Pointer
}

// Vector ...
func (e *Embedding) Vector(text string) (nlp.Vector, error) {
	t := C.CString(text)
	defer C.free(unsafe.Pointer(t))
	v := C.nl_get_vector(e.natptr, t)
	if v.data == nil {
		return nil, fmt.Errorf("no embedding for '%s'", text)
	}
	defer C.free(unsafe.Pointer(v.data))
	vec := make([]float64, v.size)
	C.memcpy(unsafe.Pointer(unsafe.SliceData(vec)), unsafe.Pointer(v.data), C.size_t(len(vec))*sizeofFloat64)
	return vec, nil
}

// NewSentenceEmbeddingForEnglish ...
func NewSentenceEmbeddingForEnglish() (*Embedding, error) {
	p := C.nl_sentence_embedding_for_english()
	if p == nil {
		return nil, errors.New("no sentence embedding for English")
	}
	emb := &Embedding{natptr: p}
	runtime.SetFinalizer(emb, func(emb *Embedding) {
		C.nl_embedding_release(emb.natptr)
	})
	return emb, nil
}

type weatherService struct{}

type forecast struct {
	Text string `json:"text"`
}

func (s *weatherService) GetForecast(ctx context.Context, in *struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}, _ *struct {
	Info any `guide:"Provides weather forecasts for the US."`
}) (*forecast, error) {
	return &forecast{
		Text: "It will be sunny with some wind.",
	}, nil
}

func main() {
	emb, err := NewSentenceEmbeddingForEnglish()
	if err != nil {
		log.Fatal(err)
	}

	var functions []*infer.Function
	funcs, err := infer.Functions(new(weatherService))
	if err != nil {
		log.Fatal(err)
	}
	functions = append(functions, funcs...)

	ctx := context.Background()

	cl, err := ai.NewClient(ctx, ai.Gemini31ProPreview)
	if err != nil {
		log.Fatal(err)
	}

	proxyTool, err := proxy.Tool(functions, emb, cl)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := cl.GenerateText(ctx, ai.NewText("What's the weather forecast for Seattle?"), []*ai.Tool{proxyTool})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)
}

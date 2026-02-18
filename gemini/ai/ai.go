package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"

	"github.com/fealsamh/go-utils/nocopy"
	"google.golang.org/genai"
)

// Client is an LLM client.
type Client struct {
	cl    *genai.Client
	model Model
}

// Model specifies an LLM model.
type Model string

const (
	// Gemini3FlashPreview represents the Gemini 3 Flash Preview model.
	Gemini3FlashPreview Model = "gemini-3-flash-preview"
	// Gemini3ProPreview represents the Gemini 3 Pro Preview model.
	Gemini3ProPreview Model = "gemini-3-pro-preview"
)

// NewClient creates a new client.
func NewClient(ctx context.Context, model Model) (*Client, error) {
	cl, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &Client{cl: cl, model: model}, nil
}

// GenerateText generates a text response.
func (cl *Client) GenerateText(ctx context.Context, in []*genai.Content, tools []*Tool) (*Response, error) {
	genaiTools := make([]*genai.Tool, 0, len(tools))
	for _, t := range tools {
		genaiTools = append(genaiTools, t.tool())
	}
	var config *genai.GenerateContentConfig
	if len(genaiTools) > 0 {
		config = &genai.GenerateContentConfig{Tools: genaiTools}
	}
	resp, err := cl.generate(ctx, in, config, tools)
	if err != nil {
		return nil, err
	}
	return &Response{resp: resp}, nil
}

// Generate generates a structured response.
func Generate[T any](ctx context.Context, cl *Client, in []*genai.Content, tools []*Tool) (*T, error) {
	schema, err := schemaFor[T]()
	if err != nil {
		return nil, err
	}
	genaiTools := make([]*genai.Tool, 0, len(tools))
	for _, t := range tools {
		genaiTools = append(genaiTools, t.tool())
	}
	config := &genai.GenerateContentConfig{
		ResponseMIMEType:   "application/json",
		ResponseJsonSchema: schema,
	}
	if len(genaiTools) > 0 {
		config.Tools = genaiTools
	}
	resp, err := cl.generate(ctx, in, config, tools)
	if err != nil {
		return nil, err
	}
	var obj T
	if err := json.Unmarshal(nocopy.Bytes(resp.Text()), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

func (cl *Client) generate(ctx context.Context, in []*genai.Content, config *genai.GenerateContentConfig, tools []*Tool) (*genai.GenerateContentResponse, error) {
	resp, err := cl.cl.Models.GenerateContent(ctx, string(cl.model), in, config)
	if err != nil {
		return nil, err
	}
	if len(resp.FunctionCalls()) > 0 {
		functions := make(map[string]func(map[string]any) (map[string]any, error))
		for _, t := range tools {
			maps.Copy(functions, t.functions)
		}
		in = append(in, resp.Candidates[0].Content)
		for _, call := range resp.FunctionCalls() {
			f, ok := functions[call.Name]
			if !ok {
				return nil, fmt.Errorf("tool function '%s' unknown", call.Name)
			}
			out, err := f(call.Args)
			if err != nil {
				return nil, err
			}
			in = append(in, genai.NewContentFromFunctionResponse(call.Name, map[string]any{"output": out}, ""))
		}
		resp, err = cl.cl.Models.GenerateContent(ctx, string(cl.model), in, config)
		if err != nil {
			return nil, err
		}
	}
	return resp, err
}

// NewText creates a new text content.
func NewText(text string) []*genai.Content {
	return genai.Text(text)
}

const (
	// MimeTypeImageJPEG is the image/jpeg type.
	MimeTypeImageJPEG = "image/jpeg"
	// MimeTypeImagePNG is the image/png type.
	MimeTypeImagePNG = "image/png"
	// MimeTypePDF is the application/pdf type.
	MimeTypePDF = "application/pdf"
)

// NewTextWithBytes creates a new text content with bytes.
func NewTextWithBytes(text string, data []byte, mimeType string) []*genai.Content {
	parts := []*genai.Part{
		genai.NewPartFromBytes(data, mimeType),
		genai.NewPartFromText(text),
	}
	return []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}
}

// Response is an LLM response.
type Response struct {
	resp *genai.GenerateContentResponse
}

func (resp *Response) String() string {
	return resp.resp.Text()
}

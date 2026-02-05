package ai

import (
	"context"
	"encoding/json"

	"github.com/google/jsonschema-go/jsonschema"
	"google.golang.org/genai"
)

// Client ...
type Client struct {
	cl    *genai.Client
	model Model
}

// Model ...
type Model string

const (
	// Gemini3FlashPreview ...
	Gemini3FlashPreview Model = "gemini-3-flash-preview"
	// Gemini3ProPreview ...
	Gemini3ProPreview Model = "gemini-3-pro-preview"
)

// NewClient ...
func NewClient(ctx context.Context, model Model) (*Client, error) {
	cl, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &Client{cl: cl, model: model}, nil
}

// GenerateText ...
func (cl *Client) GenerateText(ctx context.Context, in []*genai.Content) (*Response, error) {
	resp, err := cl.cl.Models.GenerateContent(ctx, string(cl.model), in, nil)
	if err != nil {
		return nil, err
	}
	return &Response{resp: resp}, nil
}

// Generate ...
func Generate[T any](ctx context.Context, cl *Client, in []*genai.Content) (*T, error) {
	schema, err := jsonschema.For[T](nil)
	if err != nil {
		return nil, err
	}
	config := &genai.GenerateContentConfig{
		ResponseMIMEType:   "application/json",
		ResponseJsonSchema: schema,
	}
	resp, err := cl.cl.Models.GenerateContent(ctx, string(cl.model), in, config)
	if err != nil {
		return nil, err
	}
	var obj T
	if err := json.Unmarshal([]byte(resp.Text()), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

// NewText ...
func NewText(text string) []*genai.Content {
	return genai.Text(text)
}

// Response ...
type Response struct {
	resp *genai.GenerateContentResponse
}

func (resp *Response) String() string {
	return resp.resp.Text()
}

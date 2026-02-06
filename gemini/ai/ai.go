package ai

import (
	"context"
	"encoding/json"

	"github.com/fealsamh/go-utils/nocopy"
	"github.com/google/jsonschema-go/jsonschema"
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
func (cl *Client) GenerateText(ctx context.Context, in []*genai.Content) (*Response, error) {
	resp, err := cl.cl.Models.GenerateContent(ctx, string(cl.model), in, nil)
	if err != nil {
		return nil, err
	}
	return &Response{resp: resp}, nil
}

// Generate generates a structured response.
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
	if err := json.Unmarshal(nocopy.Bytes(resp.Text()), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
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
)

// NewTextWithImage creates a new text content.
func NewTextWithImage(text string, data []byte, mimeType string) []*genai.Content {
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

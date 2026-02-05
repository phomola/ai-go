package ai

import (
	"context"

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

// Generate ...
func (cl *Client) Generate(ctx context.Context, in []*genai.Content) (*Response, error) {
	resp, err := cl.cl.Models.GenerateContent(ctx, string(cl.model), in, nil)
	if err != nil {
		return nil, err
	}
	return &Response{resp: resp}, nil
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

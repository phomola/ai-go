package infer

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type input struct {
	Name string
	Age  int
}

type output struct {
	Data string
	Num  int
}

type tool struct{}

func (t *tool) Func1(_ context.Context, in *input, _ *struct {
	Info any `guide:"GUIDE"`
}) (*output, error) {
	if in.Name == "" {
		return nil, errors.New("no name provided")
	}
	return &output{
		Data: in.Name + in.Name,
		Num:  in.Age * 2,
	}, nil
}

func TestFunctionInference(t *testing.T) {
	req := require.New(t)

	funcs, err := Functions(new(tool))
	req.Nil(err)
	req.Equal(1, len(funcs))
	req.Equal("Func1", funcs[0].Name)
	req.Equal("GUIDE", funcs[0].Description)

	out, err := funcs[0].Fn(context.Background(), map[string]any{
		"Name": "John",
		"Age":  20,
	})
	req.Nil(err)
	req.Equal(map[string]any{"Data": "JohnJohn", "Num": 40}, out)

	out, err = funcs[0].Fn(context.Background(), map[string]any{})
	req.NotNil(err)
	req.Equal("no name provided", err.Error())
}

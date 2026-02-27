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

func (t *tool) Func1(_ context.Context, in *input, _ struct {
	Info any `guide:"GUIDE"`
}) (*output, error) {
	return nil, errors.ErrUnsupported
}

func TestFunctionInference(t *testing.T) {
	req := require.New(t)

	funcs, err := Functions(new(tool))
	req.Nil(err)
	req.Equal(1, len(funcs))
	req.Equal("Func1", funcs[0].Name)
	req.Equal("GUIDE", funcs[0].Description)
}

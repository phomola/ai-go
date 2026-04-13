package rag

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTokeniser(t *testing.T) {
	req := require.New(t)

	{
		tokens := Tokenise("abcd 1234 efgh")
		req.Equal(3, len(tokens))
		req.Equal(Token{"abcd", true}, tokens[0])
		req.Equal(Token{"1234", true}, tokens[1])
		req.Equal(Token{"efgh", true}, tokens[2])
	}
}

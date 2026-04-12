package bm25

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDocument(t *testing.T) {
	req := require.New(t)

	doc := NewDocument("abcd,efgh ijkl:abcd 1234")
	req.Equal(5, doc.Length)
	req.Equal(4, len(doc.TF))
	req.Equal(2, doc.TF["abcd"])
	req.Equal(1, doc.TF["efgh"])
	req.Equal(1, doc.TF["ijkl"])
	req.Equal(1, doc.TF["1234"])
}

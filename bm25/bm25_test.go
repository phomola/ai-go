package bm25

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDocument(t *testing.T) {
	req := require.New(t)

	doc := NewDocument("", "abcd,efgh ijkl:abcd 1234")
	req.Equal(5, doc.Length)
	req.Equal(4, len(doc.TF))
	req.Equal(2, doc.TF["abcd"])
	req.Equal(1, doc.TF["efgh"])
	req.Equal(1, doc.TF["ijkl"])
	req.Equal(1, doc.TF["1234"])
}

func TestCorpus(t *testing.T) {
	req := require.New(t)

	var c Corpus
	c.AddDocument(NewDocument("1", "abcd efgh"))
	c.AddDocument(NewDocument("2", "abcd efgh ijkl 1234"))
	req.Equal(6, c.FullLength)
	req.Equal(3.0, c.AvgDocLength)

	{
		docs := c.SearchOne("abcd", 1.2, 0.75)
		req.Equal("1", docs[0].Document.ID)
	}
	{
		docs := c.SearchOne("ijkl", 1.2, 0.75)
		req.Equal("2", docs[0].Document.ID)
	}
	{
		docs := c.SearchMore([]string{"abcd"}, 1.2, 0.75)
		req.Equal("1", docs[0].Document.ID)
	}
	{
		docs := c.SearchMore([]string{"ijkl"}, 1.2, 0.75)
		req.Equal("2", docs[0].Document.ID)
	}
}

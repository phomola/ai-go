package bm25

import (
	"math"
	"sort"
)

// Document ...
type Document struct {
	ID     string
	Text   string
	Length int
	TF     map[string]int
}

// NewDocument ...
func NewDocument(id, text string) *Document {
	l, tf := 0, make(map[string]int)
	for _, t := range Tokenise(text) {
		if !t.IsAlphanum {
			continue
		}
		f := t.Form
		tf[f] = tf[f] + 1
		l++
	}
	return &Document{id, text, l, tf}
}

// GetTF ...
func (d *Document) GetTF(q string) (float64, bool) {
	if f, ok := d.TF[q]; ok {
		return float64(f), true
	}
	return 0, false
}

// Score ...
func (d *Document) Score(q string, idf, k1, b, avgDocLen float64) float64 {
	f, ok := d.GetTF(q)
	if !ok {
		return 0
	}
	return idf * f * (k1 + 1) / (f + k1*(1-b+b*float64(d.Length)/avgDocLen))
}

// Corpus ...
type Corpus struct {
	Documents    []*Document
	FullLength   int
	AvgDocLength float64
	nq           map[string]int
}

// AddDocument ...
func (c *Corpus) AddDocument(doc *Document) {
	c.Documents = append(c.Documents, doc)
	c.FullLength += doc.Length
	c.AvgDocLength = float64(c.FullLength) / float64(len(c.Documents))
	if c.nq == nil {
		c.nq = make(map[string]int)
	}
	for q := range doc.TF {
		c.nq[q] = c.nq[q] + 1
	}
}

// GetIDF ...
func (c *Corpus) GetIDF(q string) float64 {
	nq := float64(c.nq[q])
	return math.Log((float64(len(c.Documents))-nq+0.5)/(nq+0.5) + 1)
}

// Score ...
func (c *Corpus) Score(doc *Document, q string, k1, b float64) float64 {
	return doc.Score(q, c.GetIDF(q), k1, b, c.AvgDocLength)
}

// ScoredDocument ...
type ScoredDocument struct {
	Score    float64
	Document *Document
}

// SearchOne ...
func (c *Corpus) SearchOne(q string, k1, b float64) []ScoredDocument {
	scoredDocs := make([]ScoredDocument, 0, len(c.Documents))
	for _, doc := range c.Documents {
		score := c.Score(doc, q, k1, b)
		// fmt.Println(doc.ID, score)
		scoredDocs = append(scoredDocs, ScoredDocument{
			Score:    score,
			Document: doc,
		})
	}
	sort.Slice(scoredDocs, func(i, j int) bool {
		return scoredDocs[i].Score > scoredDocs[j].Score
	})
	return scoredDocs
}

// SearchMore ...
func (c *Corpus) SearchMore(qs []string, k1, b float64) []ScoredDocument {
	scoredDocs := make([]ScoredDocument, 0, len(c.Documents))
	for _, doc := range c.Documents {
		var score float64
		for _, q := range qs {
			score += c.Score(doc, q, k1, b)
		}
		// fmt.Println(doc.ID, score)
		scoredDocs = append(scoredDocs, ScoredDocument{
			Score:    score,
			Document: doc,
		})
	}
	sort.Slice(scoredDocs, func(i, j int) bool {
		return scoredDocs[i].Score > scoredDocs[j].Score
	})
	return scoredDocs
}

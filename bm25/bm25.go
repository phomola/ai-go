package bm25

// Document ...
type Document struct {
	Text   string
	Length int
	TF     map[string]int
}

// NewDocument ...
func NewDocument(text string) *Document {
	l, tf := 0, make(map[string]int)
	for _, t := range Tokenise(text) {
		if !t.IsAlphanum {
			continue
		}
		f := t.Form
		tf[f] = tf[f] + 1
		l++
	}
	return &Document{text, l, tf}
}

// Corpus ...
type Corpus struct {
	Documents []*Document
}

package nlp

// Embedding is a natural language embedding.
type Embedding interface {
	Vector(string) ([]float64, error)
}

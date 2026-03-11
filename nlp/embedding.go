package nlp

// Embedding is a natural language embedding.
type Embedding interface {
	Vector(string) (Vector, error)
}

package nlp

type Embedding interface {
	Vector(string) ([]float64, error)
}

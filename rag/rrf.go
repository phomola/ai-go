package rag

// RRF ...
func RRF(ranks []int, k float64) float64 {
	var score float64
	for _, r := range ranks {
		score += 1 / (k + float64(r))
	}
	return score
}

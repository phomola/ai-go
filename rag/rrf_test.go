package rag

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRRF(t *testing.T) {
	req := require.New(t)

	s1 := RRF([]int{1, 1}, 60)
	s2 := RRF([]int{1, 2}, 60)
	fmt.Println(s1, s2)
	req.True(s1 > s2)
}

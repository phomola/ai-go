package copier

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToMap(t *testing.T) {
	req := require.New(t)

	type (
		B struct {
			X string
		}
		A struct {
			X int
			Y string  `json:"y"`
			Z float64 `json:"-"`
			U B
			V *B
		}
	)

	x := A{
		X: 1234,
		Y: "abcd",
		Z: 12.34,
		U: B{X: "AB"},
		V: &B{X: "CD"},
	}

	m, err := ToMap(x)
	req.Nil(err)
	req.Equal(4, len(m))
	req.Equal(1234, m["X"])
	req.Equal("abcd", m["y"])
	req.Equal("AB", m["U"].(map[string]interface{})["X"])
	req.Equal("CD", m["V"].(map[string]interface{})["X"])
}

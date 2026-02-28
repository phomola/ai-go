package nlp

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"strconv"
	"strings"

	"github.com/fealsamh/go-utils/nocopy"
)

var (
	_ sql.Scanner   = (*Vector)(nil)
	_ driver.Valuer = Vector(nil)
)

// Vector is a numerical vector with driver value support.
type Vector []float64

// Value returns the driver value.
func (v Vector) Value() (driver.Value, error) {
	b := []byte{'['}
	for i, x := range v {
		if i > 0 {
			b = append(b, ',')
		}
		b = strconv.AppendFloat(b, x, 'f', -1, 64)
	}
	return append(b, ']'), nil
}

// Scan scans a value from the input.
func (v *Vector) Scan(src any) error {
	var s string
	switch x := src.(type) {
	case []byte:
		s = nocopy.String(x)
	case string:
		s = x
	default:
		return errors.New("vector from database neither string nor slice of bytes")
	}
	if len(s) < 2 {
		return errors.New("vector from database invalid")
	}
	if s[0] != '[' || s[len(s)-1] != ']' {
		return errors.New("vector from database ill-formed")
	}
	var xs []float64
	for s := range strings.SplitSeq(s[1:len(s)-1], ",") {
		x, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		xs = append(xs, x)
	}
	*v = xs
	return nil
}

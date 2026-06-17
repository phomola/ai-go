package copier

import (
	"encoding/json"
	"testing"
)

var (
	testJSON = []byte(`{"name":"Maeve","age":25,"addresses":[{"street":"Street 1","locality":"Locality 1"},{"street":"Street 2","locality":"Locality 2"}]}`)
	gr       any
)

func BenchmarkJSONObjectUnmarshal(b *testing.B) {
	type person struct {
		Name      string `json:"name"`
		Age       int    `json:"age"`
		Addresses []struct {
			Street   string `json:"street"`
			Locality string `json:"locality"`
		} `json:"addresses"`
	}
	var lr any
	for b.Loop() {
		var p person
		if err := json.Unmarshal(testJSON, &p); err != nil {
			b.Fatal(err)
		}
	}
	gr = lr
}

func BenchmarkJSONMapUnmarshal(b *testing.B) {
	var lr any
	for b.Loop() {
		var m map[string]any
		if err := json.Unmarshal(testJSON, &m); err != nil {
			b.Fatal(err)
		}
	}
	gr = lr
}

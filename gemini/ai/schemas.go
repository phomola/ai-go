package ai

import "github.com/google/jsonschema-go/jsonschema"

func schemaFor[T any]() (*jsonschema.Schema, error) {
	return jsonschema.For[T](nil)
}

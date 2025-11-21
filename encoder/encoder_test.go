package encoder

import (
	"testing"

	"github.com/Palaciodiego008/toonify/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodePrimitives(t *testing.T) {
	enc := New(types.DefaultEncodeOptions())

	tests := []struct {
		name  string
		input interface{}
	}{
		{"string", "hello"},
		{"int", 42},
		{"float", 3.14},
		{"bool_true", true},
		{"bool_false", false},
		{"nil", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := enc.Encode(tt.input)
			require.NoError(t, err)
			assert.NotEmpty(t, result)
		})
	}
}

func TestEncodeObject(t *testing.T) {
	enc := New(types.DefaultEncodeOptions())

	input := map[string]interface{}{
		"name": "Alice",
		"age":  30,
	}

	result, err := enc.Encode(input)
	require.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestEncodeArray(t *testing.T) {
	enc := New(types.DefaultEncodeOptions())

	input := []interface{}{1, 2, 3}
	result, err := enc.Encode(input)
	require.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestEncodeTabularArray(t *testing.T) {
	enc := New(types.DefaultEncodeOptions())

	input := []interface{}{
		map[string]interface{}{"id": 1, "name": "Alice"},
		map[string]interface{}{"id": 2, "name": "Bob"},
	}

	result, err := enc.Encode(input)
	require.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, string(result), "Alice")
	assert.Contains(t, string(result), "Bob")
}

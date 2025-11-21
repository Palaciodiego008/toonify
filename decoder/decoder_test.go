package decoder

import (
	"testing"

	"github.com/Palaciodiego008/toonify/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodePrimitives(t *testing.T) {
	dec := New(types.DefaultDecodeOptions())

	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{"string", "hello", "hello"},
		{"int", "42", int64(42)},
		{"float", "3.14", 3.14},
		{"bool_true", "true", true},
		{"bool_false", "false", false},
		{"null", "null", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result interface{}
			err := dec.Decode([]byte(tt.input), &result)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDecodeObject(t *testing.T) {
	dec := New(types.DefaultDecodeOptions())

	input := `name: Alice
age: 30`

	var result map[string]interface{}
	err := dec.Decode([]byte(input), &result)
	require.NoError(t, err)

	assert.Equal(t, "Alice", result["name"])
	assert.Equal(t, int64(30), result["age"])
}

func TestDecodeNestedObject(t *testing.T) {
	dec := New(types.DefaultDecodeOptions())

	input := `user:
  name: Alice
  age: 30`

	var result map[string]interface{}
	err := dec.Decode([]byte(input), &result)
	require.NoError(t, err)

	user := result["user"].(map[string]interface{})
	assert.Equal(t, "Alice", user["name"])
	assert.Equal(t, int64(30), user["age"])
}

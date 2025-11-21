package toonify

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeDecodeRoundtrip(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
	}{
		{
			name: "simple_object",
			data: map[string]interface{}{
				"name": "Alice",
				"age":  30,
			},
		},
		{
			name: "complex_structure",
			data: map[string]interface{}{
				"users": []interface{}{
					map[string]interface{}{"id": 1, "name": "Alice", "role": "admin"},
					map[string]interface{}{"id": 2, "name": "Bob", "role": "user"},
				},
				"active": true,
				"count":  2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode
			encoded, err := Encode(tt.data)
			require.NoError(t, err)
			assert.NotEmpty(t, encoded)

			// Decode
			var decoded interface{}
			err = Decode(encoded, &decoded)
			require.NoError(t, err)

			// Basic structure validation
			assert.NotNil(t, decoded)
		})
	}
}

func TestEncodeDecodeBytes(t *testing.T) {
	data := map[string]interface{}{
		"test": "value",
		"num":  42,
	}

	encoded, err := EncodeBytes(data)
	require.NoError(t, err)

	var decoded map[string]interface{}
	err = DecodeBytes(encoded, &decoded)
	require.NoError(t, err)

	assert.Equal(t, "value", decoded["test"])
	assert.Equal(t, int64(42), decoded["num"])
}

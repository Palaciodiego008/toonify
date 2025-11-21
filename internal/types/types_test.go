package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultEncodeOptions(t *testing.T) {
	opts := DefaultEncodeOptions()
	
	assert.Equal(t, 2, opts.Indent)
	assert.Equal(t, DelimiterComma, opts.Delimiter)
}

func TestDefaultDecodeOptions(t *testing.T) {
	opts := DefaultDecodeOptions()
	assert.NotNil(t, opts)
}

func TestToonError(t *testing.T) {
	err := NewToonError("test error", 5, 10)
	
	assert.Equal(t, "test error", err.Message)
	assert.Equal(t, 5, err.Line)
	assert.Equal(t, 10, err.Column)
	assert.Contains(t, err.Error(), "test error")
}

func TestValueTypes(t *testing.T) {
	// Test that Value interface can hold different types
	var values []Value
	
	values = append(values, "string")
	values = append(values, int64(42))
	values = append(values, 3.14)
	values = append(values, true)
	values = append(values, nil)
	
	obj := make(Object)
	obj["key"] = "value"
	values = append(values, obj)
	
	arr := make(Array, 0)
	arr = append(arr, 1, 2, 3)
	values = append(values, arr)
	
	assert.Len(t, values, 7)
}

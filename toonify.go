// Package toonify provides simple encoding and decoding for the TOON format.
package toonify

import (
	"github.com/Palaciodiego008/toonify/encoder"
	"github.com/Palaciodiego008/toonify/decoder"
	"github.com/Palaciodiego008/toonify/internal/types"
)

// Encode converts Go data to TOON format.
func Encode(v interface{}) (string, error) {
	opts := types.DefaultEncodeOptions()
	enc := encoder.New(opts)
	data, err := enc.Encode(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Decode converts TOON format to Go data.
func Decode(data string, v interface{}) error {
	opts := types.DefaultDecodeOptions()
	dec := decoder.New(opts)
	return dec.Decode([]byte(data), v)
}

// EncodeBytes converts Go data to TOON format as bytes.
func EncodeBytes(v interface{}) ([]byte, error) {
	opts := types.DefaultEncodeOptions()
	enc := encoder.New(opts)
	return enc.Encode(v)
}

// DecodeBytes converts TOON format bytes to Go data.
func DecodeBytes(data []byte, v interface{}) error {
	opts := types.DefaultDecodeOptions()
	dec := decoder.New(opts)
	return dec.Decode(data, v)
}

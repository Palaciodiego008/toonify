package types

import "fmt"

// Value represents any valid TOON value
type Value interface{}

// Object represents a TOON object (map)
type Object map[string]Value

// Array represents a TOON array
type Array []Value

// Delimiter represents the delimiter used for tabular arrays
type Delimiter string

const (
	DelimiterComma Delimiter = ","
	DelimiterTab   Delimiter = "\t"
	DelimiterPipe  Delimiter = "|"
)

// EncodeOptions configures TOON encoding behavior
type EncodeOptions struct {
	Indent        int       `json:"indent"`
	Delimiter     Delimiter `json:"delimiter"`
	KeyFolding    string    `json:"keyFolding"`
	FlattenDepth  int       `json:"flattenDepth"`
}

// DecodeOptions configures TOON decoding behavior
type DecodeOptions struct {
	Indent      int    `json:"indent"`
	Strict      bool   `json:"strict"`
	ExpandPaths string `json:"expandPaths"`
}

// DefaultEncodeOptions returns default encoding options
func DefaultEncodeOptions() *EncodeOptions {
	return &EncodeOptions{
		Indent:       2,
		Delimiter:    DelimiterComma,
		KeyFolding:   "off",
		FlattenDepth: 1000, // Equivalent to Number.POSITIVE_INFINITY
	}
}

// DefaultDecodeOptions returns default decoding options
func DefaultDecodeOptions() *DecodeOptions {
	return &DecodeOptions{
		Indent:      2,
		Strict:      true,
		ExpandPaths: "off",
	}
}

// ToonError represents an error during TOON processing
type ToonError struct {
	Message string
	Line    int
	Column  int
}

func (e *ToonError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("TOON error at line %d, column %d: %s", e.Line, e.Column, e.Message)
	}
	return fmt.Sprintf("TOON error: %s", e.Message)
}

// NewToonError creates a new TOON error
func NewToonError(message string, line, column int) *ToonError {
	return &ToonError{
		Message: message,
		Line:    line,
		Column:  column,
	}
}

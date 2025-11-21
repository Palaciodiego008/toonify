package utils

import (
	"reflect"
	"strings"
	"unicode"
)

// CountIndent counts the number of leading spaces in a line
func CountIndent(line string) int {
	count := 0
	for _, char := range line {
		if char == ' ' {
			count++
		} else {
			break
		}
	}
	return count
}

// NeedsQuoting determines if a string needs to be quoted in TOON format
func NeedsQuoting(s string) bool {
	if s == "" {
		return true
	}

	// Check for special values
	if s == "null" || s == "true" || s == "false" {
		return true
	}

	// Check if it looks like a number
	if isNumeric(s) {
		return true
	}

	// Check for special characters
	for _, char := range s {
		if char == ':' || char == ',' || char == '"' || char == '\n' || char == '\r' || char == '\t' {
			return true
		}
		if unicode.IsSpace(char) && (char == ' ' && (strings.HasPrefix(s, " ") || strings.HasSuffix(s, " "))) {
			return true
		}
	}

	return false
}

// isNumeric checks if a string represents a number
func isNumeric(s string) bool {
	if s == "" {
		return false
	}

	// Simple check for numeric patterns
	hasDigit := false
	hasDot := false
	
	for i, char := range s {
		if unicode.IsDigit(char) {
			hasDigit = true
		} else if char == '.' {
			if hasDot {
				return false // Multiple dots
			}
			hasDot = true
		} else if char == '-' || char == '+' {
			if i != 0 {
				return false // Sign not at beginning
			}
		} else if char == 'e' || char == 'E' {
			// Scientific notation - simplified check
			continue
		} else {
			return false
		}
	}

	return hasDigit
}

// SliceEqual compares two string slices for equality
func SliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

// FindStructField finds a struct field by name (case-insensitive)
func FindStructField(structType reflect.Type, fieldName string) (reflect.StructField, bool) {
	fieldName = strings.ToLower(fieldName)

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		
		// Check field name
		if strings.ToLower(field.Name) == fieldName {
			return field, true
		}

		// Check json tag
		if tag := field.Tag.Get("json"); tag != "" {
			tagName := strings.Split(tag, ",")[0]
			if strings.ToLower(tagName) == fieldName {
				return field, true
			}
		}

		// Check toon tag
		if tag := field.Tag.Get("toon"); tag != "" {
			tagName := strings.Split(tag, ",")[0]
			if strings.ToLower(tagName) == fieldName {
				return field, true
			}
		}
	}

	return reflect.StructField{}, false
}

// IsEmptyValue checks if a value is considered empty
func IsEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

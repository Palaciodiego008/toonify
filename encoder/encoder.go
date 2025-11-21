package encoder

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/Palaciodiego008/toonify/internal/types"
	"github.com/Palaciodiego008/toonify/internal/utils"
)

// Encoder handles TOON encoding
type Encoder struct {
	opts *types.EncodeOptions
}

// New creates a new TOON encoder
func New(opts *types.EncodeOptions) *Encoder {
	if opts == nil {
		opts = types.DefaultEncodeOptions()
	}
	return &Encoder{opts: opts}
}

// Encode encodes a value to TOON format
func (e *Encoder) Encode(v interface{}) ([]byte, error) {
	lines, err := e.EncodeLines(v)
	if err != nil {
		return nil, err
	}
	return []byte(strings.Join(lines, "\n")), nil
}

// EncodeLines encodes a value to TOON format as lines
func (e *Encoder) EncodeLines(v interface{}) ([]string, error) {
	normalized := e.normalizeValue(v)
	return e.encodeValue(normalized, 0)
}

func (e *Encoder) normalizeValue(v interface{}) interface{} {
	if v == nil {
		return nil
	}

	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Ptr:
		if val.IsNil() {
			return nil
		}
		return e.normalizeValue(val.Elem().Interface())
	case reflect.Interface:
		return e.normalizeValue(val.Elem().Interface())
	case reflect.Struct:
		result := make(map[string]interface{})
		typ := val.Type()
		for i := 0; i < val.NumField(); i++ {
			field := typ.Field(i)
			fieldValue := val.Field(i)
			
			// Skip unexported fields
			if !fieldValue.CanInterface() {
				continue
			}
			
			// Get field name from json tag or field name
			fieldName := field.Name
			if jsonTag := field.Tag.Get("json"); jsonTag != "" {
				if tagName := strings.Split(jsonTag, ",")[0]; tagName != "" && tagName != "-" {
					fieldName = tagName
				}
			} else if toonTag := field.Tag.Get("toon"); toonTag != "" {
				if tagName := strings.Split(toonTag, ",")[0]; tagName != "" && tagName != "-" {
					fieldName = tagName
				}
			}
			
			result[fieldName] = e.normalizeValue(fieldValue.Interface())
		}
		return result
	case reflect.Map:
		result := make(map[string]interface{})
		for _, key := range val.MapKeys() {
			keyStr := fmt.Sprintf("%v", key.Interface())
			result[keyStr] = e.normalizeValue(val.MapIndex(key).Interface())
		}
		return result
	case reflect.Slice, reflect.Array:
		result := make([]interface{}, val.Len())
		for i := 0; i < val.Len(); i++ {
			result[i] = e.normalizeValue(val.Index(i).Interface())
		}
		return result
	default:
		return v
	}
}

func (e *Encoder) encodeValue(v interface{}, depth int) ([]string, error) {
	if v == nil {
		return []string{"null"}, nil
	}

	switch val := v.(type) {
	case bool:
		return []string{strconv.FormatBool(val)}, nil
	case int, int8, int16, int32, int64:
		return []string{fmt.Sprintf("%d", val)}, nil
	case uint, uint8, uint16, uint32, uint64:
		return []string{fmt.Sprintf("%d", val)}, nil
	case float32, float64:
		return []string{fmt.Sprintf("%g", val)}, nil
	case string:
		return []string{e.encodeString(val)}, nil
	case map[string]interface{}:
		return e.encodeObject(val, depth)
	case []interface{}:
		return e.encodeArray(val, depth)
	default:
		return nil, types.NewToonError(fmt.Sprintf("unsupported type: %T", v), 0, 0)
	}
}

func (e *Encoder) encodeString(s string) string {
	if utils.NeedsQuoting(s) {
		return fmt.Sprintf(`"%s"`, strings.ReplaceAll(s, `"`, `\"`))
	}
	return s
}

func (e *Encoder) encodeObject(obj map[string]interface{}, depth int) ([]string, error) {
	if len(obj) == 0 {
		return []string{"{}"}, nil
	}

	var lines []string
	indent := strings.Repeat(" ", depth*e.opts.Indent)

	for key, value := range obj {
		valueLines, err := e.encodeValue(value, depth+1)
		if err != nil {
			return nil, err
		}

		if len(valueLines) == 1 && !e.isComplexValue(value) {
			lines = append(lines, fmt.Sprintf("%s%s: %s", indent, key, valueLines[0]))
		} else {
			lines = append(lines, fmt.Sprintf("%s%s:", indent, key))
			for _, line := range valueLines {
				lines = append(lines, fmt.Sprintf("%s%s", strings.Repeat(" ", e.opts.Indent), line))
			}
		}
	}

	return lines, nil
}

func (e *Encoder) encodeArray(arr []interface{}, depth int) ([]string, error) {
	if len(arr) == 0 {
		return []string{"[]"}, nil
	}

	// Check if array is tabular (uniform objects)
	if e.isTabularArray(arr) {
		return e.encodeTabularArray(arr, depth)
	}

	// Encode as regular array
	var lines []string
	indent := strings.Repeat(" ", depth*e.opts.Indent)

	for _, item := range arr {
		itemLines, err := e.encodeValue(item, depth+1)
		if err != nil {
			return nil, err
		}

		if len(itemLines) == 1 && !e.isComplexValue(item) {
			lines = append(lines, fmt.Sprintf("%s- %s", indent, itemLines[0]))
		} else {
			lines = append(lines, fmt.Sprintf("%s-", indent))
			for _, line := range itemLines {
				lines = append(lines, fmt.Sprintf("%s%s", strings.Repeat(" ", e.opts.Indent), line))
			}
		}
	}

	return lines, nil
}

func (e *Encoder) isTabularArray(arr []interface{}) bool {
	if len(arr) == 0 {
		return false
	}

	// Check if all items are objects with the same keys
	var firstKeys []string
	for i, item := range arr {
		obj, ok := item.(map[string]interface{})
		if !ok {
			return false
		}

		var keys []string
		for k := range obj {
			keys = append(keys, k)
		}

		if i == 0 {
			firstKeys = keys
		} else {
			// Check if keys match (order doesn't matter)
			if len(keys) != len(firstKeys) {
				return false
			}
			keyMap := make(map[string]bool)
			for _, k := range firstKeys {
				keyMap[k] = true
			}
			for _, k := range keys {
				if !keyMap[k] {
					return false
				}
			}
		}

		// Check if all values are primitives
		for _, v := range obj {
			if e.isComplexValue(v) {
				return false
			}
		}
	}

	return len(firstKeys) > 0
}

func (e *Encoder) encodeTabularArray(arr []interface{}, depth int) ([]string, error) {
	if len(arr) == 0 {
		return []string{"[]"}, nil
	}

	indent := strings.Repeat(" ", depth*e.opts.Indent)

	// Get field names from first object
	firstObj := arr[0].(map[string]interface{})
	var fields []string
	for k := range firstObj {
		fields = append(fields, k)
	}

	// Create header line
	header := fmt.Sprintf("[%d]{%s}:", len(arr), strings.Join(fields, ","))
	lines := []string{header}

	// Create data rows
	for _, item := range arr {
		obj := item.(map[string]interface{})
		var values []string

		for _, field := range fields {
			value := obj[field]
			valueStr, err := e.primitiveToString(value)
			if err != nil {
				return nil, err
			}
			values = append(values, valueStr)
		}

		row := fmt.Sprintf("%s%s", indent, strings.Join(values, string(e.opts.Delimiter)))
		lines = append(lines, row)
	}

	return lines, nil
}

func (e *Encoder) primitiveToString(v interface{}) (string, error) {
	switch val := v.(type) {
	case nil:
		return "", nil
	case bool:
		return strconv.FormatBool(val), nil
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", val), nil
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", val), nil
	case float32, float64:
		return fmt.Sprintf("%g", val), nil
	case string:
		if strings.Contains(val, string(e.opts.Delimiter)) || utils.NeedsQuoting(val) {
			return fmt.Sprintf(`"%s"`, strings.ReplaceAll(val, `"`, `\"`)), nil
		}
		return val, nil
	default:
		return "", types.NewToonError(fmt.Sprintf("non-primitive value in tabular array: %T", v), 0, 0)
	}
}

func (e *Encoder) isComplexValue(v interface{}) bool {
	switch v.(type) {
	case map[string]interface{}, []interface{}:
		return true
	default:
		return false
	}
}

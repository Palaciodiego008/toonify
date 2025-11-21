package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Palaciodiego008/toonify/internal/types"
	"github.com/Palaciodiego008/toonify/internal/utils"
)

var (
	tabularHeaderRegex = regexp.MustCompile(`^\[(\d+)\]\{([^}]+)\}:$`)
	arrayHeaderRegex   = regexp.MustCompile(`^(.+)\[\]:$`)
)

// Parse parses TOON format string into a Go value
func Parse(input string, opts *types.DecodeOptions) (interface{}, error) {
	if opts == nil {
		opts = types.DefaultDecodeOptions()
	}

	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) == 0 {
		return map[string]interface{}{}, nil
	}

	parser := &parser{
		lines: lines,
		opts:  opts,
		pos:   0,
	}

	return parser.parseValue(0)
}

type parser struct {
	lines []string
	opts  *types.DecodeOptions
	pos   int
}

func (p *parser) parseValue(indent int) (interface{}, error) {
	if p.pos >= len(p.lines) {
		return nil, nil
	}

	result := make(map[string]interface{})
	
	for p.pos < len(p.lines) {
		line := p.lines[p.pos]
		lineIndent := utils.CountIndent(line)

		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			p.pos++
			continue
		}

		// If indentation is less than expected, we're done with this level
		if lineIndent < indent {
			break
		}

		// If indentation is more than expected, it's an error
		if lineIndent > indent {
			return nil, types.NewToonError(fmt.Sprintf("unexpected indentation at line %d", p.pos+1), p.pos+1, lineIndent)
		}

		trimmed := strings.TrimSpace(line)

		// Check for object key-value pair
		if colonIndex := strings.Index(trimmed, ":"); colonIndex != -1 {
			key := strings.TrimSpace(trimmed[:colonIndex])
			value := strings.TrimSpace(trimmed[colonIndex+1:])

			if value == "" {
				// Multi-line value
				p.pos++
				parsedValue, err := p.parseMultiLineValue(key, indent)
				if err != nil {
					return nil, err
				}
				result[key] = parsedValue
			} else {
				// Single-line value
				parsedValue, err := p.parsePrimitive(value)
				if err != nil {
					return nil, err
				}
				result[key] = parsedValue
				p.pos++
			}
			continue
		}

		// Check for array item
		if strings.HasPrefix(trimmed, "- ") {
			arrayResult, err := p.parseArrayItems(indent)
			if err != nil {
				return nil, err
			}
			return arrayResult, nil
		}

		// Single primitive value
		return p.parsePrimitive(trimmed)
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result, nil
}

func (p *parser) parseMultiLineValue(key string, indent int) (interface{}, error) {
	if p.pos >= len(p.lines) {
		return nil, nil
	}

	nextLine := strings.TrimSpace(p.lines[p.pos])
	
	// Check if next line is a tabular array header
	if matches := tabularHeaderRegex.FindStringSubmatch(nextLine); matches != nil {
		return p.parseTabularArray(matches, indent+p.opts.Indent)
	}
	
	// Check if next line starts with "- " (regular array)
	nextLineIndent := utils.CountIndent(p.lines[p.pos])
	if nextLineIndent == indent+p.opts.Indent && strings.HasPrefix(nextLine, "- ") {
		return p.parseArrayItems(indent + p.opts.Indent)
	}
	
	// Regular multi-line value (nested object)
	return p.parseValue(indent + p.opts.Indent)
}

func (p *parser) parseTabularArray(matches []string, indent int) (interface{}, error) {
	count, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, types.NewToonError(fmt.Sprintf("invalid array count: %s", matches[1]), p.pos+1, 0)
	}

	fieldsStr := matches[2]
	fields := strings.Split(fieldsStr, ",")
	for i, field := range fields {
		fields[i] = strings.TrimSpace(field)
	}

	p.pos++ // Move past header

	var result []interface{}
	for i := 0; i < count && p.pos < len(p.lines); i++ {
		if p.pos >= len(p.lines) {
			break
		}

		line := p.lines[p.pos]
		lineIndent := utils.CountIndent(line)

		if lineIndent != indent+p.opts.Indent {
			break
		}

		trimmed := strings.TrimSpace(line)
		values := p.parseDelimitedValues(trimmed)

		if len(values) != len(fields) {
			return nil, types.NewToonError(fmt.Sprintf("field count mismatch at line %d: expected %d, got %d", p.pos+1, len(fields), len(values)), p.pos+1, 0)
		}

		obj := make(map[string]interface{})
		for j, field := range fields {
			parsedValue, err := p.parsePrimitive(values[j])
			if err != nil {
				return nil, err
			}
			obj[field] = parsedValue
		}

		result = append(result, obj)
		p.pos++
	}

	return result, nil
}

func (p *parser) parseRegularArray(key string, indent int) (interface{}, error) {
	p.pos++ // Move past header

	var items []interface{}
	for p.pos < len(p.lines) {
		line := p.lines[p.pos]
		lineIndent := utils.CountIndent(line)

		if lineIndent < indent+p.opts.Indent {
			break
		}

		if lineIndent == indent+p.opts.Indent {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "- ") {
				itemValue := strings.TrimSpace(trimmed[2:])
				if itemValue == "" {
					// Multi-line item
					p.pos++
					parsedItem, err := p.parseValue(indent + p.opts.Indent*2)
					if err != nil {
						return nil, err
					}
					items = append(items, parsedItem)
				} else {
					// Single-line item
					parsedItem, err := p.parsePrimitive(itemValue)
					if err != nil {
						return nil, err
					}
					items = append(items, parsedItem)
					p.pos++
				}
			} else {
				break
			}
		} else {
			break
		}
	}

	return map[string]interface{}{key: items}, nil
}

func (p *parser) parseArrayItems(indent int) (interface{}, error) {
	var items []interface{}

	for p.pos < len(p.lines) {
		line := p.lines[p.pos]
		lineIndent := utils.CountIndent(line)

		if lineIndent < indent {
			break
		}

		if lineIndent == indent {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "- ") {
				itemValue := strings.TrimSpace(trimmed[2:])
				if itemValue == "" {
					// Multi-line item
					p.pos++
					parsedItem, err := p.parseValue(indent + p.opts.Indent)
					if err != nil {
						return nil, err
					}
					items = append(items, parsedItem)
				} else {
					// Single-line item
					parsedItem, err := p.parsePrimitive(itemValue)
					if err != nil {
						return nil, err
					}
					items = append(items, parsedItem)
					p.pos++
				}
			} else {
				break
			}
		} else {
			break
		}
	}

	return items, nil
}

func (p *parser) parseDelimitedValues(line string) []string {
	// Simple CSV-like parsing - can be enhanced for proper CSV parsing
	var values []string
	var current strings.Builder
	inQuotes := false
	
	runes := []rune(line)
	for i := 0; i < len(runes); i++ {
		char := runes[i]
		switch char {
		case '"':
			if inQuotes && i+1 < len(runes) && runes[i+1] == '"' {
				// Escaped quote
				current.WriteRune('"')
				i++ // Skip next quote
			} else {
				inQuotes = !inQuotes
			}
		case ',':
			if !inQuotes {
				values = append(values, strings.TrimSpace(current.String()))
				current.Reset()
			} else {
				current.WriteRune(char)
			}
		default:
			current.WriteRune(char)
		}
	}
	
	values = append(values, strings.TrimSpace(current.String()))
	return values
}

func (p *parser) parsePrimitive(value string) (interface{}, error) {
	value = strings.TrimSpace(value)

	// Handle quoted strings
	if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
		unquoted := value[1 : len(value)-1]
		return strings.ReplaceAll(unquoted, `\"`, `"`), nil
	}

	// Handle null
	if value == "null" {
		return nil, nil
	}

	// Handle boolean
	if value == "true" {
		return true, nil
	}
	if value == "false" {
		return false, nil
	}

	// Try to parse as number
	if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
		return intVal, nil
	}

	if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
		return floatVal, nil
	}

	// Return as string
	return value, nil
}

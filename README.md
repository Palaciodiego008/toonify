# Toonify 🎨

[![Go Reference](https://pkg.go.dev/badge/github.com/Palaciodiego008/toonify.svg)](https://pkg.go.dev/github.com/Palaciodiego008/toonify)
[![Go Report Card](https://goreportcard.com/badge/github.com/Palaciodiego008/toonify)](https://goreportcard.com/report/github.com/Palaciodiego008/toonify)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Toonify** is a Go library for encoding and decoding [TOON (Token-Oriented Object Notation)](https://github.com/toon-format/toon) format. TOON is a compact, human-readable encoding of the JSON data model that minimizes tokens and makes structure easy for LLMs to follow.

## Features

- 🚀 **Drop-in replacement** for `encoding/json` with TOON format
- 📊 **Tabular arrays** - Uniform arrays of objects collapse into CSV-like tables
- 🎯 **LLM-optimized** - Reduces token count by ~40% compared to JSON
- 🔄 **Lossless** - Perfect round-trip compatibility with JSON data model
- ⚡ **Fast** - Efficient encoding and decoding
- 🛡️ **Type-safe** - Full Go struct support with tags

## Installation

```bash
go get github.com/Palaciodiego008/toonify
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/Palaciodiego008/toonify/pkg/toon"
)

func main() {
    // Encoding (similar to json.Marshal)
    data := map[string]interface{}{
        "name": "Alice",
        "age":  30,
        "users": []interface{}{
            map[string]interface{}{"id": 1, "name": "Alice", "role": "admin"},
            map[string]interface{}{"id": 2, "name": "Bob", "role": "user"},
        },
    }

    toonData, err := toon.Marshal(data)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(string(toonData))
    // Output:
    // name: Alice
    // age: 30
    // users:
    //   [2]{id,name,role}:
    //     1,Alice,admin
    //     2,Bob,user

    // Decoding (similar to json.Unmarshal)
    var result map[string]interface{}
    err = toon.Unmarshal(toonData, &result)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%+v\n", result)
}
```

### Struct Support

```go
type User struct {
    ID   int    `json:"id" toon:"id"`
    Name string `json:"name" toon:"name"`
    Role string `json:"role" toon:"role"`
}

type Response struct {
    Message string `json:"message" toon:"message"`
    Users   []User `json:"users" toon:"users"`
    Count   int    `json:"count" toon:"count"`
}

func main() {
    resp := Response{
        Message: "User list",
        Count:   2,
        Users: []User{
            {ID: 1, Name: "Alice", Role: "admin"},
            {ID: 2, Name: "Bob", Role: "user"},
        },
    }

    // Marshal struct to TOON
    toonData, err := toon.Marshal(resp)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(string(toonData))
    // Output:
    // message: User list
    // count: 2
    // users:
    //   [2]{id,name,role}:
    //     1,Alice,admin
    //     2,Bob,user

    // Unmarshal back to struct
    var decoded Response
    err = toon.Unmarshal(toonData, &decoded)
    if err != nil {
        log.Fatal(err)
    }
}
```

### Custom Options

```go
// Custom encoding options
opts := &toon.EncodeOptions{
    Indent:    4,                    // Use 4 spaces for indentation
    Delimiter: toon.DelimiterTab,    // Use tabs for tabular arrays
}

toonData, err := toon.MarshalWithOptions(data, opts)

// Custom decoding options
decodeOpts := &toon.DecodeOptions{
    Strict: false,  // Allow unknown fields
}

err = toon.UnmarshalWithOptions(toonData, &result, decodeOpts)
```

### Streaming Large Data

```go
// For large datasets, use EncodeLines for memory efficiency
lines, err := toon.EncodeLines(largeData, nil)
if err != nil {
    log.Fatal(err)
}

for _, line := range lines {
    fmt.Println(line)
}
```

## TOON Format Examples

### Simple Object
```go
// Go
map[string]interface{}{
    "name": "Alice",
    "age":  30,
}

// TOON
name: Alice
age: 30
```

### Tabular Array (Uniform Objects)
```go
// Go
map[string]interface{}{
    "users": []interface{}{
        map[string]interface{}{"id": 1, "name": "Alice", "active": true},
        map[string]interface{}{"id": 2, "name": "Bob", "active": false},
    },
}

// TOON
users:
  [2]{id,name,active}:
    1,Alice,true
    2,Bob,false
```

### Regular Array (Mixed Types)
```go
// Go
map[string]interface{}{
    "items": []interface{}{"apple", 42, true},
}

// TOON
items:
  - apple
  - 42
  - true
```

### Nested Objects
```go
// Go
map[string]interface{}{
    "user": map[string]interface{}{
        "profile": map[string]interface{}{
            "name": "Alice",
            "age":  30,
        },
    },
}

// TOON
user:
  profile:
    name: Alice
    age: 30
```

## API Reference

### Core Functions

- `toon.Marshal(v interface{}) ([]byte, error)` - Encode to TOON format
- `toon.Unmarshal(data []byte, v interface{}) error` - Decode from TOON format
- `toon.MarshalWithOptions(v interface{}, opts *EncodeOptions) ([]byte, error)` - Encode with options
- `toon.UnmarshalWithOptions(data []byte, v interface{}, opts *DecodeOptions) error` - Decode with options

### Convenience Functions

- `toon.EncodeToString(v interface{}) (string, error)` - Encode to string
- `toon.DecodeFromString(s string, v interface{}) error` - Decode from string
- `toon.EncodeLines(v interface{}, opts *EncodeOptions) ([]string, error)` - Encode to lines

### Options

#### EncodeOptions
```go
type EncodeOptions struct {
    Indent       int       // Indentation spaces (default: 2)
    Delimiter    Delimiter // Delimiter for tabular arrays (default: comma)
    KeyFolding   string    // Key folding strategy (default: "off")
    FlattenDepth int       // Maximum depth for flattening (default: 1000)
}
```

#### DecodeOptions
```go
type DecodeOptions struct {
    Indent      int    // Expected indentation (default: 2)
    Strict      bool   // Strict mode for unknown fields (default: true)
    ExpandPaths string // Path expansion strategy (default: "off")
}
```

## Performance

TOON typically achieves:
- **~40% fewer tokens** compared to JSON for uniform data structures
- **Better LLM comprehension** with 74% accuracy vs JSON's 70%
- **Efficient encoding/decoding** with minimal memory overhead

## Comparison with JSON

| Feature | JSON | TOON |
|---------|------|------|
| Token efficiency | Baseline | ~40% fewer tokens |
| LLM readability | Good | Better (74% vs 70% accuracy) |
| Tabular data | Verbose | Compact CSV-like format |
| Nested objects | Good | YAML-like indentation |
| Type safety | Full | Full |
| Ecosystem | Mature | Growing |

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Related Projects

- [TOON Specification](https://github.com/toon-format/spec) - Official TOON format specification
- [TOON TypeScript](https://github.com/toon-format/toon) - Reference implementation in TypeScript
- [Other TOON implementations](https://toonformat.dev/ecosystem/implementations) - Implementations in various languages

## Acknowledgments

- Thanks to the [TOON format team](https://github.com/toon-format) for creating this innovative format
- Inspired by the efficiency needs of LLM applications and token optimization

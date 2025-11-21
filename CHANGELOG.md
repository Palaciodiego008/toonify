# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-11-21

### Added
- Initial stable release of toonify
- Core TOON format encoder and decoder
- Simple public API (`Encode`, `Decode`, `EncodeBytes`, `DecodeBytes`)
- Support for tabular array encoding (compact CSV-like format)
- CLI tool for JSON/TOON conversion
- Comprehensive test suite
- Usage examples and documentation
- Build automation with Makefile

### Features
- **Simple Objects**: `key: value` format with proper indentation
- **Arrays**: Both simple (`- item`) and tabular (`[count]{columns}: data`) formats
- **Nested Structures**: Full support for complex nested data
- **Type Detection**: Automatic detection of numbers, booleans, null values
- **Human Readable**: Clean, compact format optimized for readability
- **CLI Tool**: Convert between JSON and TOON formats from command line

### API
- `toonify.Encode(data)` - Convert Go data to TOON string
- `toonify.Decode(toonData, &result)` - Convert TOON string to Go data
- `toonify.EncodeBytes(data)` - Convert Go data to TOON bytes
- `toonify.DecodeBytes(toonData, &result)` - Convert TOON bytes to Go data

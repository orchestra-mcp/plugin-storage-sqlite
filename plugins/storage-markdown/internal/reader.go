package internal

import (
	"bytes"
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/types/known/structpb"
)

// metaPrefix is the opening marker for the META comment block.
var metaPrefix = []byte("<!-- META ")

// metaSuffix is the closing marker for the META comment block.
var metaSuffix = []byte(" META -->")

// ParseMarkdownFile parses a markdown file that begins with an optional META
// comment block. The format is:
//
//	<!-- META {"key":"value"} META -->
//
//	# Markdown body here...
//
// If the file does not start with a META block, the entire content is treated
// as the body with nil metadata.
func ParseMarkdownFile(data []byte) (metadata *structpb.Struct, body []byte, err error) {
	if !bytes.HasPrefix(data, metaPrefix) {
		// No META block; the entire file is the body.
		return nil, data, nil
	}

	// Find the end of the first line (the META line).
	newlineIdx := bytes.IndexByte(data, '\n')
	if newlineIdx == -1 {
		// The entire file is a single META line with no body.
		newlineIdx = len(data)
	}

	metaLine := data[:newlineIdx]

	// Extract JSON between the prefix and suffix markers.
	if !bytes.HasSuffix(bytes.TrimRight(metaLine, "\r"), metaSuffix) {
		return nil, nil, fmt.Errorf("malformed META block: missing closing marker")
	}

	trimmed := bytes.TrimRight(metaLine, "\r")
	jsonData := trimmed[len(metaPrefix) : len(trimmed)-len(metaSuffix)]

	// Parse the JSON into a map.
	var m map[string]any
	if err := json.Unmarshal(jsonData, &m); err != nil {
		return nil, nil, fmt.Errorf("parse META JSON: %w", err)
	}

	metadata, err = structpb.NewStruct(m)
	if err != nil {
		return nil, nil, fmt.Errorf("convert META to structpb: %w", err)
	}

	// The body starts after the META line. Skip the newline and an optional
	// blank line separator.
	rest := data[newlineIdx:]
	if len(rest) > 0 && rest[0] == '\n' {
		rest = rest[1:]
	}
	// Skip one additional blank line if present (the separator between META
	// and body).
	if len(rest) > 0 && rest[0] == '\n' {
		rest = rest[1:]
	} else if len(rest) > 1 && rest[0] == '\r' && rest[1] == '\n' {
		rest = rest[2:]
	}

	return metadata, rest, nil
}

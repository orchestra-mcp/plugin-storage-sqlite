package internal

import (
	"bytes"
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/types/known/structpb"
)

// FormatMarkdownFile serializes metadata and body into the on-disk markdown
// format with a META comment header:
//
//	<!-- META {"key":"value"} META -->
//
//	# Markdown body here...
//
// If metadata is nil or has no fields, only the body is returned.
func FormatMarkdownFile(metadata *structpb.Struct, body []byte) ([]byte, error) {
	var buf bytes.Buffer

	if metadata != nil && len(metadata.Fields) > 0 {
		// Convert structpb to a plain map for compact JSON serialization.
		m := metadata.AsMap()
		jsonBytes, err := json.Marshal(m)
		if err != nil {
			return nil, fmt.Errorf("marshal META JSON: %w", err)
		}

		buf.WriteString("<!-- META ")
		buf.Write(jsonBytes)
		buf.WriteString(" META -->\n")
		buf.WriteString("\n")
	}

	buf.Write(body)

	return buf.Bytes(), nil
}

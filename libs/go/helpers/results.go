package helpers

import (
	"encoding/json"
	"fmt"

	pluginv1 "github.com/orchestrated-mcp/framework/gen/go/orchestra/plugin/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

// TextResult creates a successful ToolResponse containing a text result.
func TextResult(text string) *pluginv1.ToolResponse {
	s, _ := structpb.NewStruct(map[string]any{
		"text": text,
	})
	return &pluginv1.ToolResponse{
		Success: true,
		Result:  s,
	}
}

// JSONResult creates a successful ToolResponse containing the given data
// marshaled as a structpb.Struct. The data must be JSON-serializable.
func JSONResult(data any) (*pluginv1.ToolResponse, error) {
	raw, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal data: %w", err)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		// If the data is not a map, wrap it in one.
		m = map[string]any{"data": data}
		raw2, _ := json.Marshal(m)
		if err2 := json.Unmarshal(raw2, &m); err2 != nil {
			return nil, fmt.Errorf("convert to struct: %w", err2)
		}
	}
	s, err := structpb.NewStruct(m)
	if err != nil {
		return nil, fmt.Errorf("new struct: %w", err)
	}
	return &pluginv1.ToolResponse{
		Success: true,
		Result:  s,
	}, nil
}

// ErrorResult creates a failed ToolResponse with the given error code and message.
func ErrorResult(code string, message string) *pluginv1.ToolResponse {
	return &pluginv1.ToolResponse{
		Success:      false,
		ErrorCode:    code,
		ErrorMessage: message,
	}
}

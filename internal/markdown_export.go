package internal

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"google.golang.org/protobuf/types/known/structpb"
	"gopkg.in/yaml.v3"
)

const projectsDir = ".projects"

// exportMarkdown writes the YAML frontmatter + markdown body to .projects/
// for git visibility. This is fire-and-forget — SQLite is source of truth.
func (s *StoragePlugin) exportMarkdown(storagePath string, metadata *structpb.Struct, content []byte) {
	filePath, err := resolveMarkdownPath(s.workspace, storagePath)
	if err != nil {
		log.Printf("[storage-sqlite] export skip: %v", err)
		return
	}

	data, err := formatMarkdownFile(metadata, content)
	if err != nil {
		log.Printf("[storage-sqlite] export format error: %v", err)
		return
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("[storage-sqlite] export mkdir error: %v", err)
		return
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		log.Printf("[storage-sqlite] export write error: %v", err)
		return
	}

	// Write version sidecar for compatibility with markdown plugin.
	// (Version is already in SQLite, but sidecar keeps .projects/ self-consistent.)
}

// deleteMarkdown removes the markdown file and version sidecar from .projects/.
func (s *StoragePlugin) deleteMarkdown(storagePath string) {
	filePath, err := resolveMarkdownPath(s.workspace, storagePath)
	if err != nil {
		return
	}
	os.Remove(filePath)
	os.Remove(filePath + ".version")
}

func resolveMarkdownPath(workspace, storagePath string) (string, error) {
	if storagePath == "" {
		return "", fmt.Errorf("empty path")
	}
	cleaned := filepath.Clean(storagePath)
	if strings.Contains(cleaned, "..") {
		return "", fmt.Errorf("path traversal: %q", storagePath)
	}
	return filepath.Join(workspace, projectsDir, cleaned), nil
}

func formatMarkdownFile(metadata *structpb.Struct, body []byte) ([]byte, error) {
	var buf bytes.Buffer

	if metadata != nil && len(metadata.Fields) > 0 {
		m := metadata.AsMap()

		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		sortedMap := &yamlOrderedMap{keys: keys, m: m}
		yamlBytes, err := yaml.Marshal(sortedMap)
		if err != nil {
			return nil, fmt.Errorf("marshal YAML: %w", err)
		}

		buf.WriteString("---\n")
		buf.Write(yamlBytes)
		buf.WriteString("---\n")
		buf.WriteString("\n")
	}

	buf.Write(body)
	return buf.Bytes(), nil
}

type yamlOrderedMap struct {
	keys []string
	m    map[string]any
}

func (o *yamlOrderedMap) MarshalYAML() (any, error) {
	node := &yaml.Node{Kind: yaml.MappingNode}
	for _, k := range o.keys {
		keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: k}
		valNode := &yaml.Node{}
		if err := valNode.Encode(o.m[k]); err != nil {
			return nil, err
		}
		node.Content = append(node.Content, keyNode, valNode)
	}
	return node, nil
}

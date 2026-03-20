package internal

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const projectsDir = ".projects"

// ExportToMarkdown exports all SQLite data to .projects/ as markdown files.
// This is an on-demand operation triggered by explicit user request.
func ExportToMarkdown(workspace string) (int, error) {
	dbPath := DBPath(workspace)
	db, err := OpenDB(dbPath)
	if err != nil {
		return 0, fmt.Errorf("open database: %w", err)
	}

	exported := 0

	// Export each entity type.
	type exportSpec struct {
		table string
		dir   func(projectID string) string
		query string
	}

	specs := []exportSpec{
		{"features", func(pid string) string { return filepath.Join(pid, "features") },
			`SELECT id, project_id, title, body, status, priority, kind, assignee, estimate FROM features`},
		{"persons", func(pid string) string { return filepath.Join(pid, "persons") },
			`SELECT id, project_id, name, body, role, status, '', '', '' FROM persons`},
		{"plans", func(pid string) string { return filepath.Join(pid, "plans") },
			`SELECT id, project_id, title, body, status, '', '', '', '' FROM plans`},
		{"requests", func(pid string) string { return filepath.Join(pid, "requests") },
			`SELECT id, project_id, title, body, status, priority, kind, '', '' FROM requests`},
		{"delegations", func(pid string) string { return filepath.Join(pid, "delegations") },
			`SELECT id, project_id, question, body, status, feature_id, to_person, from_person, responded_at FROM delegations`},
		{"notes", func(_ string) string { return ".notes" },
			`SELECT id, '', title, body, '', '', '', '', '' FROM notes`},
		{"docs", func(pid string) string { return filepath.Join(pid, "docs") },
			`SELECT id, project_id, title, body, category, '', '', '', '' FROM docs`},
		{"skills", func(_ string) string { return ".skills" },
			`SELECT slug, '', name, '', description, scope, '', '', '' FROM skills`},
		{"agents", func(_ string) string { return ".agents" },
			`SELECT slug, '', name, '', description, scope, '', '', '' FROM agents`},
	}

	for _, spec := range specs {
		n, err := exportTable(db, workspace, spec.table, spec.dir, spec.query)
		if err != nil {
			log.Printf("[export] skip %s: %v", spec.table, err)
			continue
		}
		exported += n
	}

	// Export projects as project.json.
	rows, err := db.Query(`SELECT slug, name FROM projects`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var slug, name string
			if rows.Scan(&slug, &name) != nil {
				continue
			}
			projDir := filepath.Join(workspace, projectsDir, slug)
			os.MkdirAll(projDir, 0755)
			projJSON := fmt.Sprintf(`{"slug":%q,"name":%q}`, slug, name)
			if os.WriteFile(filepath.Join(projDir, "project.json"), []byte(projJSON), 0644) == nil {
				exported++
			}
		}
	}

	return exported, nil
}

func exportTable(db *sql.DB, workspace, table string, dirFunc func(string) string, query string) (int, error) {
	rows, err := db.Query(query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	exported := 0
	for rows.Next() {
		var id, projectID, name, body, f1, f2, f3, f4, f5 string
		if err := rows.Scan(&id, &projectID, &name, &body, &f1, &f2, &f3, &f4, &f5); err != nil {
			continue
		}

		dir := dirFunc(projectID)
		outDir := filepath.Join(workspace, projectsDir, dir)
		os.MkdirAll(outDir, 0755)

		// Build frontmatter.
		meta := make(map[string]string)
		meta["id"] = id
		entityType := strings.TrimSuffix(table, "s")
		meta["type"] = entityType

		switch table {
		case "features":
			if name != "" {
				meta["title"] = name
			}
			if projectID != "" {
				meta["project_slug"] = projectID
			}
			if f1 != "" {
				meta["status"] = f1
			}
			if f2 != "" {
				meta["priority"] = f2
			}
			if f3 != "" {
				meta["kind"] = f3
			}
			if f4 != "" {
				meta["assignee"] = f4
			}
			if f5 != "" {
				meta["estimate"] = f5
			}
		case "persons":
			if name != "" {
				meta["name"] = name
			}
			if projectID != "" {
				meta["project_slug"] = projectID
			}
			if f1 != "" {
				meta["role"] = f1
			}
			if f2 != "" {
				meta["status"] = f2
			}
		case "plans":
			if name != "" {
				meta["title"] = name
			}
			if projectID != "" {
				meta["project_slug"] = projectID
			}
			if f1 != "" {
				meta["status"] = f1
			}
		case "requests":
			if name != "" {
				meta["title"] = name
			}
			if projectID != "" {
				meta["project_slug"] = projectID
			}
			if f1 != "" {
				meta["status"] = f1
			}
			if f2 != "" {
				meta["priority"] = f2
			}
			if f3 != "" {
				meta["kind"] = f3
			}
		case "delegations":
			if name != "" {
				meta["question"] = name
			}
			if projectID != "" {
				meta["project_slug"] = projectID
			}
			if f1 != "" {
				meta["status"] = f1
			}
			if f2 != "" {
				meta["feature_id"] = f2
			}
			if f3 != "" {
				meta["to_person"] = f3
			}
			if f4 != "" {
				meta["from_person"] = f4
			}
			if f5 != "" {
				meta["responded_at"] = f5
			}
		case "notes":
			if name != "" {
				meta["title"] = name
			}
		case "docs":
			if name != "" {
				meta["title"] = name
			}
			if projectID != "" {
				meta["project_slug"] = projectID
			}
			if f1 != "" {
				meta["category"] = f1
			}
		case "skills", "agents":
			if name != "" {
				meta["name"] = name
			}
			if f1 != "" {
				meta["description"] = f1
			}
			if f2 != "" {
				meta["scope"] = f2
			}
		}

		var buf bytes.Buffer
		buf.WriteString("---\n")
		keys := make([]string, 0, len(meta))
		for k := range meta {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			buf.WriteString(k)
			buf.WriteString(": ")
			buf.WriteString(meta[k])
			buf.WriteByte('\n')
		}
		buf.WriteString("---\n\n")
		buf.WriteString(body)

		filename := id + ".md"
		if err := os.WriteFile(filepath.Join(outDir, filename), buf.Bytes(), 0644); err != nil {
			continue
		}
		exported++
	}
	return exported, nil
}


package internal

import "database/sql"

// workspaceSchema is the SQL schema for workspace-scoped data.
// Each workspace gets its own database at ~/.orchestra/db/<hash>.db.
const workspaceSchema = `
CREATE TABLE IF NOT EXISTS projects (
    slug TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT DEFAULT '',
    metadata TEXT DEFAULT '{}',
    body TEXT DEFAULT '',
    version INTEGER DEFAULT 1,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS features (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL REFERENCES projects(slug) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT DEFAULT '',
    status TEXT DEFAULT 'backlog',
    priority TEXT DEFAULT 'P2',
    kind TEXT DEFAULT 'feature',
    assignee TEXT DEFAULT '',
    estimate TEXT DEFAULT '',
    labels TEXT DEFAULT '[]',
    depends_on TEXT DEFAULT '[]',
    body TEXT DEFAULT '',
    version INTEGER DEFAULT 1,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);
CREATE INDEX IF NOT EXISTS idx_features_project ON features(project_id);
CREATE INDEX IF NOT EXISTS idx_features_status ON features(project_id, status);
CREATE INDEX IF NOT EXISTS idx_features_assignee ON features(project_id, assignee);

CREATE TABLE IF NOT EXISTS persons (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL REFERENCES projects(slug) ON DELETE CASCADE,
    name TEXT NOT NULL,
    email TEXT DEFAULT '',
    role TEXT DEFAULT 'developer',
    status TEXT DEFAULT 'active',
    bio TEXT DEFAULT '',
    github_email TEXT DEFAULT '',
    integrations TEXT DEFAULT '{}',
    labels TEXT DEFAULT '[]',
    body TEXT DEFAULT '',
    version INTEGER DEFAULT 1,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS plans (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT DEFAULT '',
    status TEXT DEFAULT 'draft',
    features TEXT DEFAULT '[]',
    body TEXT DEFAULT '',
    version INTEGER DEFAULT 1,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS requests (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT DEFAULT '',
    kind TEXT NOT NULL DEFAULT 'feature',
    status TEXT DEFAULT 'pending',
    priority TEXT DEFAULT 'P2',
    body TEXT DEFAULT '',
    version INTEGER DEFAULT 1,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS assignment_rules (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    kind TEXT NOT NULL,
    person_id TEXT NOT NULL,
    body TEXT DEFAULT '',
    version INTEGER DEFAULT 1,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS notes (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    title TEXT NOT NULL,
    body TEXT DEFAULT '',
    pinned INTEGER DEFAULT 0,
    deleted INTEGER DEFAULT 0,
    tags TEXT DEFAULT '[]',
    icon TEXT DEFAULT '',
    color TEXT DEFAULT '',
    version INTEGER DEFAULT 1,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS wip_config (
    project_id TEXT PRIMARY KEY,
    max_in_progress INTEGER DEFAULT 0,
    version INTEGER DEFAULT 1,
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    account_id TEXT NOT NULL,
    name TEXT DEFAULT '',
    workspace TEXT DEFAULT '',
    model TEXT DEFAULT '',
    permission_mode TEXT DEFAULT 'plan',
    allowed_tools TEXT DEFAULT '[]',
    max_budget REAL DEFAULT 0,
    system_prompt TEXT DEFAULT '',
    status TEXT DEFAULT 'active',
    message_count INTEGER DEFAULT 0,
    total_tokens_in INTEGER DEFAULT 0,
    total_tokens_out INTEGER DEFAULT 0,
    total_cost_usd REAL DEFAULT 0,
    claude_session_id TEXT DEFAULT '',
    last_message_at TEXT DEFAULT '',
    body TEXT DEFAULT '',
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS session_turns (
    session_id TEXT NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    turn_number INTEGER NOT NULL,
    user_prompt TEXT DEFAULT '',
    response TEXT DEFAULT '',
    tokens_in INTEGER DEFAULT 0,
    tokens_out INTEGER DEFAULT 0,
    cost_usd REAL DEFAULT 0,
    model TEXT DEFAULT '',
    duration_ms INTEGER DEFAULT 0,
    timestamp TEXT NOT NULL DEFAULT (datetime('now')),
    PRIMARY KEY (session_id, turn_number)
);

CREATE TABLE IF NOT EXISTS packs (
    name TEXT PRIMARY KEY,
    version TEXT NOT NULL,
    repo TEXT NOT NULL,
    installed_at TEXT NOT NULL DEFAULT (datetime('now')),
    metadata TEXT DEFAULT '{}'
);

CREATE TABLE IF NOT EXISTS stacks (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    stacks TEXT DEFAULT '[]',
    version INTEGER DEFAULT 1
);

CREATE TABLE IF NOT EXISTS hook_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_type TEXT NOT NULL,
    payload TEXT DEFAULT '{}',
    created_at TEXT DEFAULT (datetime('now'))
);

-- Generic key-value store for data that doesn't fit a dedicated table
-- (e.g., .events/hook-events.toon or any custom storage paths).
CREATE TABLE IF NOT EXISTS kv_store (
    path TEXT PRIMARY KEY,
    content BLOB,
    metadata TEXT DEFAULT '{}',
    version INTEGER DEFAULT 1,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);
`

// InitSchema runs the workspace schema on the given database.
func InitSchema(db *sql.DB) error {
	_, err := db.Exec(workspaceSchema)
	return err
}

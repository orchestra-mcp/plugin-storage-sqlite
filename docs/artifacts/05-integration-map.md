# Integration Map — Orchestra Reference

> Every external service integration: auth flow, APIs, data exchange. Extracted from `orch-ref/`.

## Overview

| Integration | Auth Method | Credential Storage | Bidirectional Sync | MCP Tools |
|------------|------------|-------------------|-------------------|-----------|
| GitHub | OAuth 2.0 / PAT | AES-256-GCM encrypted SQLite | Yes (issues, PRs, CI) | Yes (17 tools) |
| Jira | OAuth 2.0 (3LO) | AES-256-GCM encrypted SQLite | Yes (issues) | No |
| Linear | OAuth 2.0 | AES-256-GCM encrypted SQLite | Yes (issues, teams, projects, cycles) | No |
| Notion | OAuth 2.0 | AES-256-GCM encrypted SQLite | Push only (pages) | No |
| Figma | OAuth 2.0 PKCE / PAT | AES-256-GCM encrypted SQLite | Pull only (files, nodes, components) | Yes (via MCP bridge) |
| Discord | Bot token | Settings store (JSON) | Push only (notifications) | No (but MCP tools accessible via bot) |
| Slack | Bot token + App token | Settings store (JSON) | Push only (notifications) | No (but MCP tools accessible via bot) |
| Firebase | Service account JSON | File path or env var | Push only (notifications, analytics) | No |
| Orchestra Web | Laravel Sanctum token | AES-256-GCM encrypted SQLite | Yes (sync protocol) | No |
| Apple Notes | AppleScript (macOS only) | N/A (no credentials) | Push only (notes) | No |

## Shared Credential Store Pattern

All OAuth integrations (GitHub, Jira, Linear, Notion, Figma) follow the same encrypted credential storage pattern:

- **Database**: SQLite file at `<workspace>/storage/<provider>/credentials.db` (or `tokens.db` for Figma)
- **Encryption**: AES-256-GCM with a key derived from `SHA-256(hostname + "orchestra-<provider>-cred-v1")`
- **Schema**: Single-row upsert on `(provider, account_id)` unique constraint
- **Fields**: Encrypted token blob, nonce, scopes, timestamps, provider-specific metadata
- **Driver**: `modernc.org/sqlite` (pure Go, no CGo dependency)

---

## 1. GitHub Integration

### Authentication

Two authentication methods are supported:

**OAuth 2.0 Authorization Code Flow:**
1. Client requests auth URL via `GET /api/github/auth/start`
2. User is redirected to `https://github.com/login/oauth/authorize` with parameters:
   - `client_id` from env `GH_OAUTH_CLIENT_ID`
   - `redirect_uri` (default `http://localhost:19191/api/github/auth/callback`)
   - `scope`: `repo,user,read:org`
3. GitHub redirects back with `code` to callback endpoint
4. Server exchanges code for access token at `https://github.com/login/oauth/access_token`
5. Server fetches user info via `GET /user`
6. Token is encrypted and stored in SQLite

**Personal Access Token (PAT):**
1. Client sends PAT via `POST /api/github/auth/pat` with `{"token": "..."}`
2. Server validates by calling `GET /user` with the PAT
3. Token is stored with scope marker `"pat"`

### Credential Storage

- **Path**: `<workspace>/storage/github/credentials.db`
- **Table**: `credentials` with columns: `provider`, `account_id`, `token_enc`, `nonce`, `scopes`, `expires_at`, `created_at`, `updated_at`
- **Key derivation**: `SHA-256(hostname + "orchestra-github-cred-v1")`
- **Encryption**: AES-256-GCM, random nonce per save

### API Client

- **Base URL**: `https://api.github.com`
- **Auth header**: `Authorization: Bearer <token>`
- **API version header**: `X-GitHub-Api-Version: 2022-11-28`
- **Accept header**: `application/vnd.github+json`
- **Timeout**: 15 seconds
- **Methods**: GET, POST, PUT, PATCH, DELETE

### Data Exchange

**Inbound (read from GitHub):**
- Issues: list (with state/assignee/labels/sort filter), get single, paginated (30 per page), filters out PRs from `/issues` endpoint
- Pull Requests: list (with state/sort/direction/head/base filter), get single, list files, list reviews
- CI/CD: check runs status aggregation for a commit ref (pass/fail/pending counts + overall status)
- User: profile information for auth verification

**Outbound (write to GitHub):**
- Issues: create (with title, body, assignees, labels), close, add comment
- Pull Requests: create (with title, head, base, body, draft), merge (merge/squash/rebase), submit review (APPROVE/REQUEST_CHANGES/COMMENT), add comment

**Repo Tracking:**
- Local markdown file at `<project>/.github/tracked-repos.md`
- Track/untrack repositories with PR and issue watching flags

### MCP Tools (17 tools)

**Auth tools (4):**
- `github_auth_start` — Get OAuth authorization URL
- `github_auth_status` — Check authentication state (connected/disconnected + user info)
- `github_auth_pat` — Sign in with Personal Access Token
- `github_sign_out` — Disconnect and remove stored credentials

**Issue tools (6):**
- `github_list_issues` — List issues for owner/repo with state/assignee/labels filters
- `github_get_issue` — Get single issue by number
- `github_create_issue` — Create issue with title, body, assignees, labels
- `github_close_issue` — Close an issue by number
- `github_issue_comment` — Add comment to an issue
- `github_ci_status` — Get CI/CD status for a commit ref

**PR tools (7):**
- `github_list_prs` — List pull requests with state/sort/direction filters
- `github_get_pr` — Get PR details by number
- `github_pr_files` — List files changed in a PR
- `github_create_pr` — Create PR with title, head, base, body, draft flag
- `github_merge_pr` — Merge PR with method (merge/squash/rebase)
- `github_review_pr` — Submit review (APPROVE/REQUEST_CHANGES/COMMENT)
- `github_pr_comment` — Add comment to a PR

**Repo tools (3):**
- `github_list_repos` — List tracked repositories for a project
- `github_track_repo` — Add a repository to tracking
- `github_untrack_repo` — Remove a repository from tracking

### Settings API Routes

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/github/auth/start` | Returns OAuth authorization URL |
| GET | `/api/github/auth/callback` | OAuth redirect handler (exchanges code for token) |
| GET | `/api/github/auth/status` | Check authentication state |
| POST | `/api/github/auth/pat` | Sign in with Personal Access Token |
| POST | `/api/github/disconnect` | Sign out and remove credentials |

### Source Files

| File | Description |
|------|-------------|
| `app/github/hub.go` | Hub facade: Auth, PRs, Issues, CI sub-services, mutex-protected refresh |
| `app/github/auth.go` | OAuth code flow, PAT auth, status check, sign out |
| `app/github/client.go` | HTTP client with Bearer auth, JSON request/response handling |
| `app/github/config.go` | Default config from env vars (GH_OAUTH_CLIENT_ID, etc.) |
| `app/github/credentials.go` | AES-256-GCM encrypted SQLite credential store |
| `app/github/issue_service.go` | Issue CRUD, comment, PR filtering from /issues endpoint |
| `app/github/pr_service.go` | PR CRUD, merge, review, files, comments |
| `app/github/ci_service.go` | CI check-runs aggregation (pass/fail/pending/overall) |
| `app/settings/github.go` | Settings server HTTP handlers for auth flow |
| `app/tools/github_auth.go` | MCP tools for GitHub auth (4 tools) |
| `app/tools/github_issues.go` | MCP tools for issues and CI (6 tools) |
| `app/tools/github_pr.go` | MCP tools for pull requests (7 tools) |
| `app/tools/github_repos.go` | MCP tools for repo tracking (3 tools) |

---

## 2. Jira Integration

### Authentication

**OAuth 2.0 (3LO — Three-Legged OAuth):**
1. Client requests auth URL via `GET /api/jira/auth/start`
2. User is redirected to `https://auth.atlassian.com/authorize` with parameters:
   - `audience`: `api.atlassian.com`
   - `client_id` from env `JIRA_CLIENT_ID`
   - `redirect_uri` (default `http://127.0.0.1:19191/api/jira/auth/callback`)
   - `scope`: `read:jira-work write:jira-work read:jira-user offline_access`
   - `response_type`: `code`
   - `prompt`: `consent`
3. Atlassian redirects back with `code` to callback endpoint
4. Server exchanges code for tokens at `https://auth.atlassian.com/oauth/token`
5. Server fetches accessible resources at `https://api.atlassian.com/oauth/token/accessible-resources` to get `cloud_id` and site URL
6. Server verifies identity via `GET /ex/jira/{cloudId}/rest/api/3/myself`
7. Access token, refresh token, cloud ID, and site URL are encrypted and stored

### Credential Storage

- **Path**: `<workspace>/storage/jira/credentials.db`
- **Table**: `credentials` with columns: `provider`, `access_token_enc`, `nonce`, `refresh_token`, `cloud_id`, `site_url`, `created_at`, `updated_at`
- **Key derivation**: `SHA-256(hostname + "orchestra-jira-cred-v1")`
- **Unique constraint**: `UNIQUE(provider)` (single credential per provider)
- **Additional data**: Stores cloud_id and site_url alongside the encrypted token

### API Client

- **Base URL**: Dynamic per site (e.g., `https://mycompany.atlassian.net`)
- **API path prefix**: `/rest/api/3`
- **Auth**: Two modes — Basic (email:token base64) or Bearer (OAuth access token). Hub uses Bearer auth.
- **Headers**: `Accept: application/json`, `Content-Type: application/json`
- **Timeout**: 30 seconds
- **Methods**: GET, POST, PUT, DELETE

### Data Exchange

**Inbound (read from Jira):**
- Issues: JQL search (with maxResults), get single by key (e.g., `PROJ-123`)
- Transitions: list available status transitions for an issue

**Outbound (write to Jira):**
- Issues: create, update fields, transition (change status via transition ID)
- Comments: add plain-text comment (converted to Atlassian Document Format)

### Settings API Routes

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/jira/auth/start` | Returns Atlassian OAuth authorization URL |
| GET | `/api/jira/auth/callback` | OAuth redirect handler (multi-step token exchange) |
| GET | `/api/jira/auth/status` | Check authentication state (validates via /myself) |
| POST | `/api/jira/disconnect` | Sign out and remove credentials |

### Source Files

| File | Description |
|------|-------------|
| `app/jira/hub.go` | Hub facade: Auth, Issues sub-services, Bearer auth only |
| `app/jira/auth.go` | OAuth 2.0 3LO flow, resource discovery, /myself verification, sign out |
| `app/jira/client.go` | HTTP client with Basic or Bearer auth, REST API v3 prefix |
| `app/jira/config.go` | Default config from env vars (JIRA_CLIENT_ID, JIRA_CLIENT_SECRET) |
| `app/jira/credentials.go` | AES-256-GCM encrypted SQLite store with cloud_id and site_url |
| `app/jira/issue_service.go` | JQL search, issue CRUD, transitions, ADF-formatted comments |
| `app/jira/auth_test.go` | Auth status and sign-out tests |
| `app/jira/client_test.go` | Client tests: Basic auth, Bearer auth, URL paths, POST body, errors |
| `app/settings/jira.go` | Settings server HTTP handlers for auth flow |

---

## 3. Linear Integration

### Authentication

**OAuth 2.0 Authorization Code Flow:**
1. Client requests auth URL via `GET /api/linear/auth/start`
2. User is redirected to `https://linear.app/oauth/authorize` with parameters:
   - `client_id` from env `LINEAR_CLIENT_ID`
   - `redirect_uri` (default `http://127.0.0.1:19191/api/linear/auth/callback`)
   - `response_type`: `code`
   - `scope`: `read,write`
3. Linear redirects back with `code` to callback endpoint
4. Server exchanges code for access token at `https://api.linear.app/oauth/token`
5. Server verifies identity via GraphQL `viewer` query
6. Access token is encrypted and stored

### Credential Storage

- **Path**: `<workspace>/storage/linear/credentials.db`
- **Table**: `credentials` with columns: `provider`, `access_token_enc`, `nonce`, `created_at`, `updated_at`
- **Key derivation**: `SHA-256(hostname + "orchestra-linear-cred-v1")`
- **Simpler schema**: No refresh token, no scopes — just encrypted access token

### API Client

- **Endpoint**: `https://api.linear.app/graphql`
- **Protocol**: GraphQL (POST only)
- **Auth header**: `Authorization: <token>` (raw token, no "Bearer" prefix)
- **Content-Type**: `application/json`
- **Timeout**: 30 seconds
- **Methods**: Query and Mutate (both are POST to GraphQL endpoint)

### Data Exchange

**Inbound (read from Linear):**
- Issues: list by team (with fields: id, identifier, title, description, priority, state, assignee, team, project, labels, parent), get single by ID
- Teams: list all teams (id, name, key)
- Projects: list all projects (id, name, state, slugId)
- Cycles: list cycles for a team (id, number, name, startsAt, endsAt, progress)
- Viewer: identity verification (id, name, displayName, email, avatarUrl, active)

**Outbound (write to Linear):**
- Issues: create (via `issueCreate` mutation), update (via `issueUpdate` mutation)
- Comments: add comment to issue (via `commentCreate` mutation)

### GraphQL Queries

```graphql
# Issue fields fragment (shared across queries)
id identifier title description priority createdAt updatedAt
state { id name type color }
assignee { id name displayName email }
team { id name key }
project { id name state slugId }
labels { nodes { id name color } }
parent { id identifier }

# Key operations
query issues($teamId, $first)       # List by team
query issue($id)                     # Get single
mutation issueCreate($input)         # Create
mutation issueUpdate($id, $input)    # Update
mutation commentCreate($issueId, $body) # Comment
query teams                          # List teams
query projects                       # List projects
query team($teamId).cycles           # List cycles
query viewer                         # Auth verification
```

### Settings API Routes

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/linear/auth/start` | Returns Linear OAuth authorization URL |
| GET | `/api/linear/auth/callback` | OAuth redirect handler |
| GET | `/api/linear/auth/status` | Check authentication state |
| POST | `/api/linear/disconnect` | Sign out and remove credentials |

### Source Files

| File | Description |
|------|-------------|
| `app/linear/hub.go` | Hub facade: Auth, Issues sub-services |
| `app/linear/auth.go` | OAuth 2.0 flow, viewer verification, sign out |
| `app/linear/client.go` | GraphQL client with raw token auth, error handling |
| `app/linear/config.go` | Default config from env vars (LINEAR_CLIENT_ID, LINEAR_CLIENT_SECRET) |
| `app/linear/credentials.go` | AES-256-GCM encrypted SQLite store (simplified schema) |
| `app/linear/issue_service.go` | GraphQL queries/mutations for issues, teams, projects, cycles |
| `app/linear/auth_test.go` | Auth status, sign-out, credential save/load/delete tests |
| `app/linear/client_test.go` | Client tests: body, auth header, content-type, GraphQL errors, HTTP errors |
| `app/settings/linear.go` | Settings server HTTP handlers for auth flow |

---

## 4. Notion Integration

### Authentication

**OAuth 2.0 Authorization Code Flow:**
1. Client requests auth URL via `GET /api/notion/auth/start`
2. User is redirected to `https://api.notion.com/v1/oauth/authorize` with parameters:
   - `client_id` from env `NOTION_CLIENT_ID`
   - `redirect_uri` (default `http://localhost:19191/api/notion/auth/callback`)
   - `response_type`: `code`
   - `owner`: `user`
3. Notion redirects back with `code` to callback endpoint
4. Server exchanges code for access token at `https://api.notion.com/v1/oauth/token`
   - Uses **Basic auth** with `base64(client_id:client_secret)` (Notion-specific requirement)
   - Request body is JSON, not form-encoded
5. Response includes: access_token, bot_id, workspace_id, workspace_name, workspace_icon
6. Token is encrypted and stored with workspace_id as account_id

### Credential Storage

- **Path**: `<workspace>/storage/notion/credentials.db`
- **Table**: `credentials` with columns: `provider`, `account_id`, `token_enc`, `nonce`, `scopes`, `expires_at`, `created_at`, `updated_at`
- **Key derivation**: `SHA-256(hostname + "orchestra-notion-cred-v1")`

### API Client

- **Base URL**: `https://api.notion.com/v1`
- **Auth header**: `Authorization: Bearer <token>`
- **Version header**: `Notion-Version: 2022-06-28`
- **Content-Type**: `application/json`
- **Timeout**: 30 seconds

### Data Exchange

**Outbound (write to Notion):**
- Pages: create under a database or under another page, with title and content blocks
- Markdown conversion: headings (h1-h3), bullet lists, code blocks, dividers, paragraphs converted to Notion block objects

**Inbound (read from Notion):**
- Databases: search for accessible databases (via POST `/search` with filter)
- User validation: GET `/users/me` for status check

### Markdown to Blocks Converter

The `MarkdownToBlocks()` function converts markdown to Notion API block objects:
- `# Heading` to `heading_1` block
- `## Heading` to `heading_2` block
- `### Heading` to `heading_3` block
- `- Item` or `* Item` to `bulleted_list_item` block
- `` ```code``` `` to `code` block (with language detection)
- `---` / `***` / `___` to `divider` block
- Plain text to `paragraph` block

### Settings API Routes

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/notion/config` | Check if OAuth credentials are configured |
| GET | `/api/notion/auth/start` | Returns OAuth authorization URL |
| GET | `/api/notion/auth/callback` | OAuth redirect handler |
| GET | `/api/notion/auth/status` | Check authentication state |
| POST | `/api/notion/disconnect` | Sign out and remove credentials |
| GET | `/api/notion/databases` | List accessible databases |
| POST | `/api/notion/send` | Create a page (title, content, database_id) |

### Source Files

| File | Description |
|------|-------------|
| `app/notion/hub.go` | Hub facade: Auth, SendNote, ListDatabases |
| `app/notion/auth.go` | OAuth 2.0 flow with Basic auth token exchange, status, sign out |
| `app/notion/client.go` | REST client, CreatePage (under DB or page), SearchDatabases |
| `app/notion/config.go` | Default config from env vars (NOTION_CLIENT_ID, NOTION_CLIENT_SECRET) |
| `app/notion/credentials.go` | AES-256-GCM encrypted SQLite credential store |
| `app/notion/markdown.go` | Markdown to Notion blocks converter |
| `app/settings/notion.go` | Settings server HTTP handlers |

---

## 5. Figma Integration

### Authentication

Two authentication methods are supported:

**OAuth 2.0 PKCE Authorization Code Flow:**
1. Client requests auth URL via `GET /api/figma/auth/start`
2. Server generates PKCE code verifier (32 random bytes, URL-safe base64) and challenge (SHA-256 of verifier)
3. User is redirected to `https://www.figma.com/oauth` with parameters:
   - `client_id` from env `FIGMA_CLIENT_ID`
   - `redirect_uri` from env (default not hardcoded)
   - `scope`: `current_user:read file_content:read`
   - `state`: random 16-byte value
   - `response_type`: `code`
   - `code_challenge`: SHA-256 of verifier (S256 method)
   - `code_challenge_method`: `S256`
4. Figma redirects back with `code` to callback endpoint
5. Server exchanges code + verifier at `https://www.figma.com/api/oauth/token`
6. Server fetches user identity via `GET https://api.figma.com/v1/me`
7. Token, user info, and expiry are stored (Figma does not issue refresh tokens)

**Personal Access Token (PAT):**
1. Client sends PAT via `POST /api/figma/auth/pat` with `{"token": "..."}`
2. Server validates by calling `GET https://api.figma.com/v1/me` with `X-Figma-Token` header
3. Token is stored with scope marker `"pat"` and far-future expiry (year 9999)

### Credential Storage

- **Path**: `<workspace>/storage/figma/tokens.db`
- **Table**: `figma_tokens` with columns: `user_id`, `user_name`, `user_email`, `access_token_enc`, `nonce`, `refresh_token`, `scopes`, `expires_at`, `created_at`, `updated_at`
- **Key derivation**: `SHA-256(hostname + "orchestra-figma-token-v1")`
- **Single-row**: Always stored as `id=1`, upserted on save
- **Expiry tracking**: `NeedsRefresh()` returns true if token expires within 5 minutes

### API Clients

**Auth header differentiation:**
- OAuth tokens: `Authorization: Bearer <token>`
- PAT tokens: `X-Figma-Token: <token>`
- Both checked via `token.Scopes == "pat"` flag

**FilesClient** (with in-memory cache):
- **Endpoint**: `GET https://api.figma.com/v1/me/files`
- **Cache TTL**: 5 minutes (in-memory)
- **Method**: `ListFiles()` returns key, name, thumbnail_url, last_modified
- **Method**: `InvalidateCache()` forces fresh API call

**ProxyClient** (with retry logic):
- **Endpoints**:
  - `GET https://api.figma.com/v1/files/{key}` — Full file data
  - `GET https://api.figma.com/v1/files/{key}/nodes?ids={nodeId}` — Specific nodes
  - `GET https://api.figma.com/v1/files/{key}/components` — Component list
- **Retry**: 3 attempts with exponential backoff (100ms, 400ms, 1600ms)
- **Rate limiting**: Respects `Retry-After` header on 429 responses

### Data Exchange

**Inbound (pull from Figma):**
- Files: list user's files with metadata
- File content: full file JSON tree
- Nodes: specific node subtrees by ID
- Components: component list for a file
- User identity: id, handle, email

### MCP Bridge

The Figma MCP integration allows installing the `figma-developer-mcp` npm package into `.mcp.json`:
- **Status check**: Detects if `figma` entry exists in `.mcp.json` and if `figma-developer-mcp` process is running
- **Install**: Adds `{"command": "npx", "args": ["-y", "figma-developer-mcp", "--figma-api-key=${FIGMA_API_KEY}"]}` to `.mcp.json`

### Settings API Routes

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/figma/auth/start` | Returns OAuth authorization URL (PKCE) |
| GET | `/api/figma/auth/callback` | OAuth redirect handler (PKCE exchange) |
| GET | `/api/figma/auth/status` | Check authentication state |
| POST | `/api/figma/auth/pat` | Save Personal Access Token (validates first) |
| POST | `/api/figma/auth/disconnect` | Sign out and remove credentials |
| GET | `/api/figma/files` | List user's Figma files (cached) |
| GET | `/api/figma/files/{key}` | Get full file data (proxy) |
| GET | `/api/figma/files/{key}/nodes` | Get specific nodes (proxy, `?ids=` param) |
| GET | `/api/figma/files/{key}/components` | Get components (proxy) |
| GET | `/api/figma/mcp/status` | Check if figma-developer-mcp is installed/running |
| POST | `/api/figma/mcp/install` | Install figma-developer-mcp into .mcp.json |

### Source Files

| File | Description |
|------|-------------|
| `app/figma/auth.go` | OAuth 2.0 PKCE flow, PAT auth, PKCE verifier/challenge generation |
| `app/figma/token.go` | AES-256-GCM encrypted SQLite token store with NeedsRefresh |
| `app/figma/files.go` | FilesClient with 5-minute in-memory TTL cache |
| `app/figma/proxy.go` | ProxyClient with 3-attempt retry + rate-limit backoff |
| `app/figma/auth_test.go` | Auth flow tests |
| `app/figma/files_test.go` | Files client tests |
| `app/figma/proxy_test.go` | Proxy client tests |
| `app/figma/token_test.go` | Token store tests |
| `app/settings/figma_auth.go` | Auth settings handlers (OAuth, PAT, disconnect) |
| `app/settings/figma_files.go` | File listing handler |
| `app/settings/figma_mcp.go` | MCP bridge status and install handlers |
| `app/settings/figma_proxy.go` | File/nodes/components proxy handlers |

---

## 6. Discord Integration

### Authentication

- **Method**: Bot token (from settings store, not OAuth)
- **No user-facing auth flow**: Bot token is configured via the settings system
- **Config fields**: `enabled`, `webhook_url`, `bot_token`, `application_id`, `channel_id`, `guild_id`, `command_prefix`

### Gateway Connection

- **Protocol**: WebSocket to Discord Gateway v10
- **URL discovery**: `GET https://discord.com/api/v10/gateway`
- **Connection sequence**: Connect -> Hello (op 10, heartbeat interval) -> Identify (op 2, with intents and presence)
- **Intents**: `GUILDS` (1<<0) | `GUILD_MESSAGES` (1<<9) | `MESSAGE_CONTENT` (1<<15)
- **Presence**: Online status with activity "Playing Orchestra MCP"
- **Heartbeat**: Periodic op 1 messages at the interval specified by Hello
- **Events handled**: `MESSAGE_CREATE`, `INTERACTION_CREATE`

### REST Client

- **Base URL**: `https://discord.com/api/v10`
- **Auth header**: `Authorization: Bot <token>`
- **Timeout**: 15 seconds
- **Operations**: SendMessage, EditMessage, RespondInteraction, RegisterSlashCommands

### Data Exchange

**Outbound (push to Discord):**
- Workflow transition notifications: rich embeds with status emoji, color coding, task/status fields, project footer
- Interactive messages: buttons, action rows
- Slash command responses
- Thread-based AI conversations

**Notification triggers:**
- Implements `workflow.TransitionListener` interface
- Fires on every task state transition (e.g., todo -> in-progress)
- Sends via webhook URL (preferred) or bot API to channel

### Bot Architecture

The Discord bot has a full command system:

**Command routing:**
- Prefix commands: `!<command>` (customizable prefix, default `!`)
- Slash commands: registered via Discord API (guild-specific for instant update, global as fallback)
- Interaction handlers: button clicks and component interactions

**Registered handlers (10):**
| Handler | Description |
|---------|-------------|
| `PingHandler` | Basic ping/pong |
| `StatusHandler` | Bot and workspace status |
| `ChatHandler` | AI-powered chat (default for unmatched commands) |
| `PermissionHandler` | MCP tool permission requests (approve/deny buttons) |
| `StopHandler` | Stop running AI sessions |
| `McpHandler` | Execute MCP tools from Discord |
| `ToolsHandler` | List available MCP tools |
| `PromptsHandler` | List available MCP prompts |
| `ActionsHandler` | Interactive action handling |
| `ProgressHandler` | Live progress updates for running tasks |

### Source Files

| File | Description |
|------|-------------|
| `app/discord/bot.go` | Bot lifecycle, handler registration, config from settings, tool/prompt maps |
| `app/discord/gateway.go` | WebSocket gateway: connect, identify, heartbeat, read loop |
| `app/discord/rest.go` | REST API client: messages, interactions, slash command registration |
| `app/discord/service.go` | Workflow notification service with embed builder |
| `app/discord/handler.go` | Handler, HandlerAPI, InteractionHandler interfaces |
| `app/discord/router.go` | Command router: prefix matching, slash dispatch, component interactions |
| `app/discord/types.go` | Discord types: MessageCreate, InteractionCreate, SlashCommandDef, Component |
| `app/discord/embed_helpers.go` | Embed builders: success, error, info, warning, tool, action, permission |
| `app/discord/handlers/*.go` | Individual command handlers (9 files) |
| `app/bot/manager.go` | Bot manager: starts Discord and Slack, registers all handlers and tools |

---

## 7. Slack Integration

### Authentication

- **Method**: Bot token (`xoxb-*`) + App token (`xapp-*`) for Socket Mode
- **No user-facing auth flow**: Tokens are configured via the settings system
- **Config fields**: `enabled`, `bot_token`, `app_token`, `channel_id`, `webhook_url`, `command_prefix`

### Socket Mode Connection

- **URL discovery**: `POST https://slack.com/api/apps.connections.open` with Bearer app_token
- **Protocol**: WebSocket with Socket Mode envelope format
- **Connection sequence**: Connect -> Hello event -> Read loop
- **ACK**: Every event with an `envelope_id` is acknowledged immediately
- **Reconnection**: Automatic with exponential backoff (1s to 30s)

### REST Client

- **Base URL**: `https://slack.com/api`
- **Auth header**: `Authorization: Bearer <bot_token>`
- **Timeout**: 15 seconds
- **Rate limiting**: Respects `Retry-After` header on 429 responses
- **Operations**: `chat.postMessage`, `chat.update`, response URL posting

### Data Exchange

**Outbound (push to Slack):**
- Workflow transition notifications: Block Kit messages with colored attachments
- Interactive messages: buttons via action blocks
- Thread replies
- Slash command responses via response_url

**Notification triggers:**
- Implements `workflow.TransitionListener` interface
- Fires on every task state transition
- Sends via webhook URL (preferred) or bot API to channel

### Bot Architecture

Mirror of Discord bot with Slack-specific types:

**Event handling:**
- Socket Mode events: `events_api` (messages), `slash_commands`, `interactive` (block_actions)
- Prefix commands: `!<command>` (customizable)
- Slash commands: native Slack slash commands

**Block Kit types:**
- `Block` (header, section, context, actions, fields)
- `TextObj` (mrkdwn, plain_text)
- `ButtonElement` with action_id and styles
- `Attachment` with colored sidebar wrapping blocks

**Registered handlers (10):** Same as Discord — Ping, Status, Chat, Permission, Stop, MCP, Tools, Prompts, Actions, Progress

### Source Files

| File | Description |
|------|-------------|
| `app/slack/bot.go` | Bot lifecycle, handler registration, config from settings |
| `app/slack/socket.go` | Socket Mode: connect, ACK, read loop, auto-reconnect |
| `app/slack/rest.go` | REST client: chat.postMessage, chat.update, response_url |
| `app/slack/service.go` | Workflow notification service with Block Kit builder |
| `app/slack/handler.go` | Handler, HandlerAPI, InteractionHandler interfaces |
| `app/slack/router.go` | Command router for messages, slash commands, interactions |
| `app/slack/types.go` | Slack types: SlackEvent, MessageEvent, SlashCommand, BlockAction, Block |
| `app/slack/block_helpers.go` | Block Kit builder utilities |
| `app/slack/handlers/*.go` | Individual command handlers (10 files) |

---

## 8. Firebase Integration

### Authentication

- **Method**: Google Service Account credentials
- **Sources (in priority order)**:
  1. `CredentialsJSON` — JSON string (for CI/CD environments)
  2. `CredentialsPath` — Path to `firebase.json` service account file
- **Fallback**: If neither is provided, Firebase is disabled gracefully (no error)
- **SDK**: Firebase Admin SDK for Go (`firebase.google.com/go/v4`)

### Services

**Cloud Messaging (FCM):**
- `SendPushNotification(token, title, body, imageURL, data)` — Single device
- `SendMulticastNotification(tokens, title, body, imageURL, data)` — Multiple devices
- `SendTopicNotification(topic, title, body, imageURL, data)` — Topic subscribers
- `SubscribeToTopic(tokens, topic)` — Subscribe devices to a topic
- `UnsubscribeFromTopic(tokens, topic)` — Unsubscribe devices from a topic

**Analytics (GA4 Measurement Protocol):**
- `LogEvent(event)` — Log custom analytics event (local logging, extensible)
- `LogEventBatch(events)` — Batch event logging
- `SendToMeasurementProtocol(measurementID, apiSecret, event)` — Server-side GA4 events via `https://www.google-analytics.com/mp/collect`

**Crashlytics (Server-Side):**
- `ReportCrash(report)` — Log crash with message, stack trace, device info, custom keys
- `ReportError(err, customKeys)` — Log non-fatal error
- Note: Firebase Admin SDK does not have direct Crashlytics API; mobile/web SDKs handle this automatically

### Data Exchange

**Outbound only:**
- Push notifications to device tokens and topics
- Analytics events to GA4
- Server-side crash/error logging

### Source Files

| File | Description |
|------|-------------|
| `app/firebase/client.go` | Firebase Admin SDK init, FCM messaging (push, multicast, topic) |
| `app/firebase/analytics.go` | Analytics event logging, GA4 Measurement Protocol |
| `app/firebase/crashlytics.go` | Server-side crash/error reporting |
| `app/firebase/types.go` | Request types: PushRequest, MulticastRequest, TopicRequest, AnalyticsEvent, CrashReport |

---

## 9. Orchestra Web (Auth)

### Authentication

**Laravel Sanctum Token Flow:**
1. Client sends email + password via `POST /api/auth/login`
2. Server forwards to Laravel web app at `POST <baseURL>/api/auth/login`
3. Laravel returns a Sanctum bearer token + user info (id, name, email, status, avatar, subscription, roles)
4. Token is encrypted and stored locally
5. Subsequent requests use `Authorization: Bearer <token>` header

**Base URL**: Configurable via env `ORCHESTRA_WEB_URL`, default `https://orchestra-mcp.dev`

### Credential Storage

- **Path**: `<configDir>/storage/auth/credentials.db`
- **Table**: `credentials` with columns: `provider`, `account_id`, `token_enc`, `nonce`, `user_json`, `created_at`, `updated_at`
- **Key derivation**: `SHA-256(hostname + "orchestra-auth-cred-v1")`
- **Provider name**: `orchestra-web`
- **Additional data**: Stores serialized user JSON alongside the encrypted token

### API Client

- **Base URL**: `https://orchestra-mcp.dev` (or env override)
- **Endpoints**:
  - `POST /api/auth/login` — Authenticate with email/password
  - `GET /api/auth/user` — Fetch authenticated user profile
  - `POST /api/auth/logout` — Revoke Sanctum token
- **Auth header**: `Authorization: Bearer <token>`
- **Timeout**: 30 seconds

### Data Exchange

**Bidirectional sync:**
- Integration status sync: pushes connected/disconnected state of all integrations to web dashboard
- Uses `syncclient.PushRecord` with entity_type "integration"
- Each integration gets a deterministic entity ID from `SHA-256("integration:" + deviceID + ":" + provider)`
- Supports upsert (connected) and delete (disconnected) actions

**User data:**
- `UserInfo`: id, name, email, status, avatar_url, subscription (plan, status), roles

### Settings API Routes

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/auth/login` | Authenticate with email/password |
| GET | `/api/auth/status` | Check authentication state |
| POST | `/api/auth/store` | Save token + user from direct API login |
| GET | `/api/auth/token` | Get stored bearer token |
| POST | `/api/auth/logout` | Sign out and clear credentials |
| POST | `/api/integrations/sync` | Push all integration states to web dashboard |

### Source Files

| File | Description |
|------|-------------|
| `app/auth/hub.go` | Hub: login, logout, status, credential management |
| `app/auth/client.go` | WebAuthClient: Laravel Sanctum HTTP client |
| `app/auth/config.go` | Web URL config (env ORCHESTRA_WEB_URL) |
| `app/auth/store.go` | AES-256-GCM encrypted SQLite store with user_json |
| `app/settings/auth.go` | Settings server HTTP handlers for auth |
| `app/settings/integration_sync.go` | Integration sync: collects all provider states, pushes to web |

---

## 10. Apple Notes Integration

### Authentication

- **Method**: No authentication required — uses macOS AppleScript via `osascript`
- **Platform**: macOS only (build tag `//go:build darwin`)
- **Non-macOS**: Returns `ErrNotSupported` for all operations

### AppleScript Bridge

**Operations:**
- `SendNote(title, htmlBody, folder)` — Creates a new note via AppleScript `tell application "Notes"`
- `ListFolders()` — Returns all folder names via AppleScript (pipe-delimited)
- `Available()` — Returns `true` on macOS, `false` on all other platforms

**Security:**
- String escaping: backslashes and double quotes are escaped for safe AppleScript embedding
- No credential storage needed — relies on macOS user session permissions

### Data Exchange

**Outbound only (push to Apple Notes):**
- Markdown content is converted to HTML (h1-h3, bold, italic, code, lists, links, hr, paragraphs)
- HTML is passed to Apple Notes via AppleScript `make new note with properties {name:..., body:...}`
- Notes are created in a specified folder (default: "Notes")

### Markdown to HTML Converter

The `MarkdownToHTML()` function converts markdown to Apple Notes-compatible HTML:
- `# Heading` to `<h1>...</h1>`
- `**bold**` / `__bold__` to `<b>...</b>`
- `*italic*` / `_italic_` to `<i>...</i>`
- `` `code` `` to `<code>...</code>`
- `` ```code block``` `` to `<pre>...</pre>`
- `- Item` / `* Item` to `<ul><li>...</li></ul>`
- `---` to `<hr>`
- Empty lines to `<br>`
- Plain text to `<p>...</p>`

### Settings API Routes

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/applenotes/status` | Check if Apple Notes is available on this platform |
| GET | `/api/applenotes/folders` | List Apple Notes folders |
| POST | `/api/applenotes/send` | Create a note (title, content as markdown, folder) |

### Source Files

| File | Description |
|------|-------------|
| `app/applenotes/send_darwin.go` | macOS implementation: SendNote, ListFolders via osascript |
| `app/applenotes/send_other.go` | Non-macOS stub: returns ErrNotSupported |
| `app/applenotes/markdown.go` | Markdown to HTML converter with inline formatting |
| `app/settings/applenotes.go` | Settings server HTTP handlers |

---

## Cross-Cutting Concerns

### Integration Sync Protocol

All OAuth integrations participate in the integration sync protocol (`app/settings/integration_sync.go`):

1. On connect/disconnect, `pushIntegrationSync(provider)` fires asynchronously
2. Collects auth state from each hub (GitHub, Jira, Linear, Notion, Figma)
3. Builds `PushRecord` with entity_type `"integration"`:
   - Connected: `upsert` action with provider, account_id, user_name, user_email
   - Disconnected: `delete` action
4. Pushes to Orchestra Web via sync client
5. Web dashboard "Connected Accounts" page stays up to date

### Bot Manager

The `app/bot/manager.go` orchestrates both Discord and Slack bots:

1. Creates bot instances with workspace path and AI bridge
2. Collects all MCP tools from the tools package (Project, Epic, Story, Task, Workflow, Lifecycle, PRD, Bugfix, Usage, Readme, Artifacts)
3. Registers 10 command handlers per bot
4. Starts permission event listeners for interactive approve/deny flows
5. Provides unified status reporting for both bots

### Environment Variables Reference

| Variable | Integration | Description |
|----------|------------|-------------|
| `GH_OAUTH_CLIENT_ID` | GitHub | OAuth client ID |
| `GH_OAUTH_CLIENT_SECRET` | GitHub | OAuth client secret |
| `GH_OAUTH_REDIRECT_URL` / `GH_OAUTH_REDIRECT_URI` | GitHub | OAuth redirect URI |
| `GH_OAUTH_SCOPES` | GitHub | OAuth scopes (default: `repo,user,read:org`) |
| `JIRA_CLIENT_ID` | Jira | OAuth client ID |
| `JIRA_CLIENT_SECRET` | Jira | OAuth client secret |
| `JIRA_REDIRECT_URL` | Jira | OAuth redirect URI |
| `LINEAR_CLIENT_ID` | Linear | OAuth client ID |
| `LINEAR_CLIENT_SECRET` | Linear | OAuth client secret |
| `LINEAR_REDIRECT_URI` / `LINEAR_REDIRECT_URL` | Linear | OAuth redirect URI |
| `NOTION_CLIENT_ID` | Notion | OAuth client ID |
| `NOTION_CLIENT_SECRET` | Notion | OAuth client secret |
| `NOTION_REDIRECT_URL` / `NOTION_REDIRECT_URI` | Notion | OAuth redirect URI |
| `FIGMA_CLIENT_ID` | Figma | OAuth client ID |
| `FIGMA_CLIENT_SECRET` | Figma | OAuth client secret |
| `FIGMA_API_KEY` | Figma | Personal Access Token for MCP bridge |
| `ORCHESTRA_WEB_URL` | Orchestra Web | Laravel web app URL (default: `https://orchestra-mcp.dev`) |

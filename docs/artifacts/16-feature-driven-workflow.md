# Feature-Driven Workflow — Orchestra

> AI agents don't do Scrum. They work best with one self-contained feature document.
> The doc IS the spec. The doc IS the memory. The doc IS the deliverable.
> Doc-first. TDD-driven. Human-reviewed. **Cyclical until done.**

---

## 0. The Cycle (Core Principle)

The workflow is fundamentally a **continuous cycle**, not a linear pipeline. The feature is never "done" until the human explicitly says so.

```
    ┌──────────────────────────────────────────────────────────┐
    │                   THE DELIVERY CYCLE                     │
    │                                                          │
    │   ┌─────┐    ┌───────────┐    ┌──────┐    ┌──────────┐  │
    │   │ DOC │───→│ IMPLEMENT │───→│ TEST │───→│ REVIEW   │  │
    │   └──┬──┘    └───────────┘    └──────┘    └────┬─────┘  │
    │      ▲                                         │         │
    │      │         ┌──────────┐    ┌────────┐      │         │
    │      └─────────│ UPDATE   │←───│ HUMAN  │←─────┘         │
    │                │ DOC      │    │ REVIEW │                │
    │                └──────────┘    └───┬────┘                │
    │                                    │                     │
    │                              ┌─────▼─────┐              │
    │                              │  APPROVED? │              │
    │                              └─────┬─────┘              │
    │                               no   │   yes              │
    │                            ┌───────┴──────┐             │
    │                            ▼              ▼              │
    │                     back to DOC        DONE              │
    │                     (new cycle)     (feature closed)     │
    └──────────────────────────────────────────────────────────┘
```

**Every cycle iteration**:
1. **Doc** — Write or update the feature document (context, plan, DOD, tests)
2. **Implement** — Agent writes code based on the doc
3. **Test** — Agent runs test cases defined in the doc
4. **Review** — Agent self-reviews against DOD
5. **Human** — Human reviews the feature doc + code
6. **Update** — If human requests edits, update the doc with new version → back to step 2

The doc accumulates evidence across cycles. Each cycle adds a new version to the Implementation Log. The human sees the full history and decides when to close.

**"Done" means the human said "done."** Not the agent. Not the tests. Not the gates. The human.

---

## 1. Why Not Scrum

Scrum's epic → story → task hierarchy was designed for human teams:
- Humans need small, bite-sized tasks (2-4 hour chunks)
- Humans need standups to sync context across team members
- Humans need sprint boundaries to create urgency

AI agents are different:
- They lose context when work is fragmented across 50 small task files
- They work best when they can see the ENTIRE feature in one document
- They don't need standups or sprint boundaries — they execute continuously
- They need a clear spec, clear tests, and clear acceptance criteria in ONE place

**The fix**: Replace the epic/story/task hierarchy with a single **Feature Document** that contains everything the agent needs to deliver a feature end-to-end.

---

## 2. The Feature Document

One Markdown file per feature. Contains: context, plan, DOD, test cases, implementation status, evidence, version history. Sized to fit the agent's context window.

### File: `.projects/{project}/features/{feature-id}.md`

```markdown
<!-- META
{
  "id": "FEAT-001",
  "title": "User authentication with JWT",
  "status": "in-progress",
  "priority": "high",
  "labels": ["backend", "auth"],
  "created_at": "2026-02-26T10:00:00Z",
  "updated_at": "2026-02-26T14:30:00Z",
  "version": 3,

  // --- Assignment ---
  "assignee": "go-architect",
  "domain": "backend",
  "required_skills": ["go", "auth", "jwt"],
  "model": "claude-opus-4-6",
  "context_budget": 180000,

  // --- Dependencies ---
  "depends_on": ["FEAT-000"],
  "blocked_by": [],
  "parent": null,
  "children": ["FEAT-001a", "FEAT-001b"],

  // --- Time Management ---
  "estimate": {
    "complexity": "medium",
    "cycles_estimated": 2,
    "cycles_actual": 3
  },
  "time": {
    "started_at": "2026-02-26T10:05:00Z",
    "completed_at": null,
    "cycle_log": [
      {"cycle": 1, "started": "2026-02-26T10:05:00Z", "ended": "2026-02-26T11:20:00Z", "duration_min": 75},
      {"cycle": 2, "started": "2026-02-26T11:25:00Z", "ended": "2026-02-26T12:10:00Z", "duration_min": 45},
      {"cycle": 3, "started": "2026-02-26T13:00:00Z", "ended": null, "duration_min": null}
    ]
  },

  // --- AgentOps (Cost & Token Tracking) ---
  "agentops": {
    "total_input_tokens": 245000,
    "total_output_tokens": 38000,
    "total_cost_usd": 4.82,
    "cycle_costs": [
      {"cycle": 1, "input_tokens": 120000, "output_tokens": 18000, "cost_usd": 2.34, "api_calls": 12},
      {"cycle": 2, "input_tokens": 85000, "output_tokens": 12000, "cost_usd": 1.63, "api_calls": 8},
      {"cycle": 3, "input_tokens": 40000, "output_tokens": 8000, "cost_usd": 0.85, "api_calls": 5}
    ],
    "model_usage": {
      "claude-opus-4-6": {"calls": 20, "cost_usd": 4.32},
      "claude-haiku-4-5": {"calls": 5, "cost_usd": 0.50}
    }
  },

  // --- Git ---
  "git": {
    "branch": "feat/FEAT-001-auth",
    "base": "main",
    "merge_strategy": "squash",
    "commits": [
      {"cycle": 1, "sha": "abc1234", "message": "feat(FEAT-001): v1 — initial auth"},
      {"cycle": 2, "sha": "def5678", "message": "feat(FEAT-001): v2 — fix rotation"}
    ]
  },

  // --- File Locks ---
  "files_locked": ["handlers/auth_handler.go", "middleware/auth.go"],

  // --- Gates & Reviews ---
  "gate_evidence": {
    "testing": "All 12 tests passing. See ## Test Results.",
    "review": "Human approved v2. Edits applied in v3."
  },
  "review_history": [
    {"version": 2, "reviewer": "fady", "verdict": "needs-edits", "note": "Add refresh token rotation"},
    {"version": 3, "reviewer": "fady", "verdict": "approved"}
  ]
}
META -->

# FEAT-001: User authentication with JWT

## Context

The API server needs authentication. Users log in with email/password and receive
a JWT access token (15min) and refresh token (7 days). All protected endpoints
validate the access token. Refresh tokens support rotation.

### Relevant Files
- `services/orch-server/internal/handlers/auth_handler.go` (create)
- `services/orch-server/internal/services/auth_service.go` (create)
- `services/orch-server/internal/middleware/auth.go` (create)
- `proto/orchestra/server/v1/auth.proto` (modify)

## Plan

1. Define auth protobuf messages (LoginRequest, TokenPair, RefreshRequest)
2. Create auth_service.go with bcrypt password hashing + JWT signing
3. Create auth_handler.go with POST /login, POST /refresh, POST /logout
4. Create auth middleware that validates Bearer token on protected routes
5. Add refresh token rotation (new refresh token on each use, old one invalidated)
6. Write tests for all endpoints + edge cases

## Definition of Done

- [ ] POST /login returns access + refresh token for valid credentials
- [ ] POST /login returns 401 for invalid credentials
- [ ] Protected routes return 401 without token
- [ ] Protected routes return 401 with expired token
- [ ] POST /refresh returns new token pair and invalidates old refresh token
- [ ] POST /logout invalidates refresh token
- [ ] All passwords stored as bcrypt hashes, never plaintext
- [ ] Rate limiting on login endpoint (5 attempts per minute)

## Test Cases

```go
// auth_handler_test.go
func TestLogin_ValidCredentials(t *testing.T) { /* ... */ }
func TestLogin_InvalidPassword(t *testing.T) { /* ... */ }
func TestLogin_NonexistentUser(t *testing.T) { /* ... */ }
func TestLogin_RateLimited(t *testing.T) { /* ... */ }
func TestRefresh_ValidToken(t *testing.T) { /* ... */ }
func TestRefresh_ExpiredToken(t *testing.T) { /* ... */ }
func TestRefresh_RotatedToken(t *testing.T) { /* ... */ }
func TestRefresh_ReusedToken(t *testing.T) { /* ... */ }
func TestProtectedRoute_NoToken(t *testing.T) { /* ... */ }
func TestProtectedRoute_ValidToken(t *testing.T) { /* ... */ }
func TestProtectedRoute_ExpiredToken(t *testing.T) { /* ... */ }
func TestLogout_InvalidatesRefresh(t *testing.T) { /* ... */ }
```

## Implementation Log

### v1 (initial)
- Created auth_service.go with bcrypt + JWT
- Created auth_handler.go with login/refresh/logout
- Created auth middleware
- 10/12 tests passing

### v2 (after testing gate)
- Fixed refresh token rotation — was not invalidating old token
- 12/12 tests passing
- Submitted for human review

### v3 (after human review)
> **Human review (fady)**: "Add refresh token rotation detection — if a rotated-out
> token is reused, invalidate ALL tokens for that user (security measure)."

- Added token family tracking
- Added reuse detection that invalidates entire family
- Added TestRefresh_ReusedToken test
- 12/12 tests passing

## Test Results

```
=== RUN   TestLogin_ValidCredentials          PASS (0.02s)
=== RUN   TestLogin_InvalidPassword           PASS (0.01s)
=== RUN   TestLogin_NonexistentUser           PASS (0.01s)
=== RUN   TestLogin_RateLimited               PASS (0.05s)
=== RUN   TestRefresh_ValidToken              PASS (0.02s)
=== RUN   TestRefresh_ExpiredToken            PASS (0.01s)
=== RUN   TestRefresh_RotatedToken            PASS (0.02s)
=== RUN   TestRefresh_ReusedToken             PASS (0.03s)
=== RUN   TestProtectedRoute_NoToken          PASS (0.01s)
=== RUN   TestProtectedRoute_ValidToken       PASS (0.01s)
=== RUN   TestProtectedRoute_ExpiredToken     PASS (0.01s)
=== RUN   TestLogout_InvalidatesRefresh       PASS (0.02s)
PASS
ok  	orch-server/internal/handlers	0.22s
```

## Files Changed
- `services/orch-server/internal/handlers/auth_handler.go` (+245 lines)
- `services/orch-server/internal/handlers/auth_handler_test.go` (+380 lines)
- `services/orch-server/internal/services/auth_service.go` (+178 lines)
- `services/orch-server/internal/middleware/auth.go` (+67 lines)
- `proto/orchestra/server/v1/auth.proto` (+42 lines)
```

---

## 3. Workflow (Feature Lifecycle)

Same gate philosophy as the original workflow, but applied to features and **explicitly cyclical**. The `needs-edits → in-progress → test → human-review` loop is the heartbeat of the system, not an edge case.

```
┌─────────┐    ┌──────┐
│ backlog │───→│ todo │
└─────────┘    └──┬───┘
                  │
                  ▼
           ┌─────────────┐  ◄─── CYCLE ENTRY POINT
           │ in-progress │  ◄──────────────────────────────────┐
           └──────┬──────┘                                      │
                  │                                             │
                  ▼                                             │
      ┌───────────────────┐                                     │
      │ ready-for-testing │                                     │
      └─────────┬─────────┘                                     │
                │                                               │
                ▼                                               │
         ┌──────────────┐                                       │
         │ in-testing   │                                       │
         └──────┬───────┘                                       │
                │                                               │
                ▼                                               │
           ┌────────┐                                           │
           │ tested │                                           │
           └───┬────┘                                           │
               │                                                │
               ▼                                                │
      ┌───────────────┐      ┌─────────────┐                   │
      │ human-review  │─────→│ needs-edits │───────────────────┘
      └───────┬───────┘      └─────────────┘
              │                  (human requests changes →
              │                   agent updates doc →
              ▼                   new cycle begins)
           ┌──────┐
           │ done │  ◄── ONLY when human says "approved"
           └──────┘

      ┌────────────┐
      │  blocked   │  (waiting on dependency — any state can enter)
      └────────────┘
      ┌────────────┐
      │  rejected  │  (feature cancelled — human only)
      └────────────┘
```

**The cycle**: `in-progress → ready-for-testing → in-testing → tested → human-review → needs-edits → in-progress` (repeat). Every pass through the cycle adds a new version to the feature doc. The doc grows with evidence until the human approves.

### States

| State | Who | What Happens |
|-------|-----|-------------|
| `backlog` | — | Feature exists but not prioritized |
| `todo` | — | Prioritized, ready to be picked up |
| `in-progress` | Agent | Agent reads feature doc, implements code |
| `ready-for-testing` | Gate | Agent declares implementation complete, provides evidence |
| `in-testing` | Agent | Agent runs test cases defined in the doc |
| `tested` | Gate | All tests pass, evidence added to doc |
| `human-review` | Human | Human reads the feature doc + code, decides |
| `needs-edits` | Human | Human writes edit requests in the doc |
| `done` | — | Feature delivered and approved |
| `blocked` | — | Waiting on dependency |
| `rejected` | Human | Feature cancelled |

### Gate transitions

| From → To | Evidence Required | Who Advances |
|-----------|-------------------|-------------|
| backlog → todo | — | Human or agent |
| todo → in-progress | — | Agent (picks up feature) |
| in-progress → ready-for-testing | "Implementation complete. Files changed: ..." | Agent |
| ready-for-testing → in-testing | — | Agent |
| in-testing → tested | "All N tests passing. See ## Test Results" | Agent |
| tested → human-review | — | Automatic |
| human-review → done | "Approved by {reviewer}" | Human |
| human-review → needs-edits | "Edit requests: ..." | Human |
| needs-edits → in-progress | — | Agent (reads edits, starts new version) |
| any → blocked | "Blocked by FEAT-XXX" | Agent or human |
| any → rejected | "Reason: ..." | Human |

---

## 4. Context-Aware Sizing

Every feature must fit within the agent's context window. This is a hard constraint.

### The Rule

```
feature_doc_size + relevant_code_size + agent_system_prompt ≤ model_context_window
```

### Metadata fields

```json
{
  "model": "claude-opus-4-6",
  "context_budget": 180000
}
```

The `context_budget` is in tokens. The feature doc + referenced code must fit.

### When a feature is too big

Split it into child features:

```json
{
  "id": "FEAT-001",
  "children": ["FEAT-001a", "FEAT-001b", "FEAT-001c"]
}
```

Each child is a self-contained feature doc that fits in the context window. Children can depend on each other via `depends_on`.

### The tool `create_feature` handles sizing

When a request comes in:
1. Estimate total complexity
2. If it fits in one context window → create single feature doc
3. If too big → split into child features, each sized to fit
4. Each child contains: its portion of the plan, its DOD items, its test cases

---

## 5. Version Control Inside the Doc

The feature doc tracks its own history in the `## Implementation Log` section:

```markdown
## Implementation Log

### v1 (initial)
- Created auth_service.go with bcrypt + JWT
- 10/12 tests passing

### v2 (testing gate)
- Fixed refresh token rotation
- 12/12 tests passing
- Submitted for human review

### v3 (human review edits)
> **Edit request (fady)**: Add token reuse detection
- Added token family tracking
- Added reuse detection
- 12/12 tests passing
```

The `version` field in metadata increments each time. The `review_history` array tracks all human review rounds.

This means the agent can:
1. See the full history of what was tried
2. See what the human asked for
3. See what was changed and why
4. Continue from where it left off

---

## 6. Tools (Replacing Scrum Tools)

The old 45 Scrum tools (project/epic/story/task CRUD, sprint management, scrum ceremonies) are replaced with feature-focused tools:

### Core Feature Tools

| Tool | What It Does |
|------|-------------|
| `create_project` | Create a project |
| `list_projects` | List all projects |
| `delete_project` | Delete a project |
| `get_project_status` | Project status summary (features by state) |
| `create_feature` | Create a feature doc (auto-sizes to context window) |
| `get_feature` | Read a feature document |
| `update_feature` | Update a feature document (adds new version) |
| `list_features` | List features with filters (status, priority, label, assignee) |
| `delete_feature` | Delete a feature |
| `search` | Full-text search across all features |

### Workflow Tools

| Tool | What It Does |
|------|-------------|
| `get_next_feature` | Get next feature to work on (respects priority, deps, WIP limits) |
| `set_current_feature` | Mark a feature as in-progress |
| `advance_feature` | Move feature to next gate (with evidence) |
| `reject_feature` | Reject/cancel a feature |
| `get_workflow_status` | Current feature pipeline (how many in each state) |
| `request_review` | Submit feature for human review |
| `submit_review` | Human submits review (approve or request edits) |
| `get_pending_reviews` | List features waiting for human review |

### Dependency & Planning Tools

| Tool | What It Does |
|------|-------------|
| `add_dependency` | FEAT-002 depends on FEAT-001 |
| `remove_dependency` | Remove a dependency |
| `get_dependency_graph` | Visualize dependency chain |
| `set_wip_limits` | Max in-progress features (global + per-assignee) |
| `get_wip_limits` | Current WIP limits |
| `check_wip_limit` | Check if assignee can take more work |
| `split_feature` | Split a feature into sized children |

### AgentOps Tools

| Tool | What It Does |
|------|-------------|
| `record_usage` | Record an API call (model, tokens, cost) for the current feature cycle |
| `get_usage` | Get token/cost summary for a feature (per cycle or total) |
| `get_project_cost` | Aggregate cost across all features in a project |
| `set_budget` | Set cost budget for a feature (with alert threshold) |

### Time Management Tools

| Tool | What It Does |
|------|-------------|
| `set_estimate` | Set complexity + estimated cycles for a feature |
| `get_estimate_accuracy` | Compare estimated vs actual cycles across completed features |
| `get_cycle_time` | Get average cycle duration by domain + complexity |

### Reporting Tools

| Tool | What It Does |
|------|-------------|
| `get_progress` | Overall delivery progress (features done / total, cost, time) |
| `get_blocked_features` | List all blocked features and why |
| `get_review_queue` | Features waiting for human review |

### Assignment & Metadata

| Tool | What It Does |
|------|-------------|
| `assign_feature` | Assign feature to agent + domain (auto-suggests from labels) |
| `unassign_feature` | Unassign |
| `add_labels` | Add labels to a feature |
| `remove_labels` | Remove labels |
| `save_note` | Save a project note |
| `list_notes` | List project notes |
| `search_notes` | Search notes |

### Git Tools

| Tool | What It Does |
|------|-------------|
| `create_branch` | Create feature branch from main |
| `get_branch_status` | Branch diff vs main, ahead/behind |
| `merge_feature` | Merge feature branch to main on approval |

### Notification Tools

| Tool | What It Does |
|------|-------------|
| `send_notification` | Send alert through configured channels |
| `get_notification_config` | Get notification settings |
| `set_notification_config` | Update notification settings |

### Conflict Resolution Tools

| Tool | What It Does |
|------|-------------|
| `lock_files` | Register files being modified by a feature |
| `unlock_files` | Release file locks |
| `check_conflicts` | Check file overlap with other in-progress features |
| `resolve_conflict` | Mark conflict as resolved with explanation |

### Quality Gate Tools

| Tool | What It Does |
|------|-------------|
| `run_quality_check` | Run all quality checks on a feature branch |
| `get_quality_config` | Get quality gate configuration |
| `set_quality_config` | Update quality gate configuration |
| `get_quality_report` | Get last quality check results |

**Total: ~58 tools** (feature CRUD 10 + workflow 8 + dependencies 7 + agentops 4 + time 3 + reporting 3 + assignment/metadata 7 + git 3 + notifications 3 + conflicts 4 + quality 4 + notes 3)

---

## 7. Directory Structure

```
.projects/
  {project-slug}/
    project.md                    # Project metadata + description
    features/
      FEAT-001.md                 # Feature document (self-contained)
      FEAT-001a.md                # Child feature (split from FEAT-001)
      FEAT-001b.md                # Child feature
      FEAT-002.md                 # Another feature
    notes/
      architecture-decisions.md   # Project notes
      meeting-2026-02-26.md
    config.md                     # Project config (WIP limits, labels, etc.)
```

**Key change**: No more `epics/*/stories/*/tasks/` hierarchy. Flat list of feature docs under `features/`. Relationships are via metadata (`parent`, `children`, `depends_on`), not directory nesting.

---

## 8. The Agent Loop (Cyclical)

The agent loop is a **while loop**, not a sequence. It runs until the human says "done."

```
┌──────────────────────────────────────────────────────────────┐
│ AGENT LOOP — runs until human approves                       │
│                                                              │
│ 1. PICK UP FEATURE                                           │
│    get_next_feature → FEAT-003 (todo, highest priority)      │
│    set_current_feature("FEAT-003") → in-progress             │
│                                                              │
│ 2. READ DOC                                                  │
│    get_feature("FEAT-003")                                   │
│    → Full doc: context, plan, DOD, test cases                │
│    → On cycle 2+: also sees human's edit requests            │
│                                                              │
│ 3. IMPLEMENT                                                 │
│    Agent reads doc, writes/updates code                      │
│    update_feature("FEAT-003", { log: "vN: ...", files: []})  │
│                                                              │
│ 4. TEST                                                      │
│    Agent runs test cases defined in the doc                  │
│    advance_feature(evidence="Implementation complete")       │
│    → in-testing                                              │
│    Agent captures test output                                │
│    update_feature("FEAT-003", { test_results: "..." })       │
│    advance_feature(evidence="All 12 tests passing")          │
│    → tested                                                  │
│                                                              │
│ 5. SUBMIT FOR REVIEW                                         │
│    → tested → human-review (automatic)                       │
│    request_review("FEAT-003")                                │
│    Agent STOPS. Waits for human.                             │
│                                                              │
│ 6. HUMAN REVIEWS                                             │
│    Human reads: feature doc + code + test results            │
│    Human decides:                                            │
│    ├─ submit_review(verdict="approved")                      │
│    │  → done. Feature closed. Go to step 1 (next feature).  │
│    └─ submit_review(verdict="needs-edits", note="Add X")     │
│       → needs-edits → in-progress                            │
│       → BACK TO STEP 2 (new cycle, new version)              │
│                                                              │
│ CYCLE REPEATS:                                               │
│   v1: initial implementation → human: "add refresh tokens"   │
│   v2: add refresh tokens → human: "add reuse detection"      │
│   v3: add reuse detection → human: "approved" → DONE         │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

### What each cycle produces

Each cycle through the loop adds to the feature doc:
- A new `### vN` entry in `## Implementation Log`
- Updated `## Test Results` with latest output
- Updated `## Files Changed` with cumulative changes
- A new entry in `review_history` metadata array
- Incremented `version` in metadata
- A new entry in `agentops.cycle_costs[]` (tokens, cost, API calls)
- A new entry in `time.cycle_log[]` (started, ended, duration)
- Updated `## Cost Summary` table

The doc is the **single source of truth** that grows across cycles. By the time the feature is "done," the doc contains a complete history of every attempt, every human feedback, every fix, how much it cost, and how long it took.

---

## 9. What Changes from Artifact 15

| Artifact 15 (Scrum) | Now (Feature-Driven) |
|---------------------|----------------------|
| `tools.scrum` plugin (45 tools) | `tools.features` plugin (~35 tools) |
| epic/story/task hierarchy | Flat feature docs with parent/children metadata |
| Sprint management (7 tools) | Removed — no sprints for AI agents |
| Scrum ceremonies (5 tools) | Replaced by progress/review tools |
| `epics/*/stories/*/tasks/*.md` | `features/FEAT-xxx.md` |
| Task = small work item | Feature = self-contained deliverable with DOD + tests |
| No human review gate | Human review gate after testing |
| Version in metadata only | Version control inside the doc (Implementation Log) |
| Context window not considered | Context budget is a hard constraint per feature |

### What stays the same

- Plugin host architecture (orchestrator + plugins)
- storage.markdown plugin (reads/writes feature docs)
- transport.stdio plugin (MCP JSON-RPC)
- QUIC + Protobuf between plugins
- `google.protobuf.Struct` for metadata
- JSON metadata block + Markdown body format
- libs/go/ shared packages
- Workflow gates with evidence

---

## 10. Updated Build Order

| Step | What | Change from Artifact 15 |
|------|------|------------------------|
| 1 | Proto contract | Same |
| 2 | Plugin SDK (`libs/go/`) | Same, but `types/` has Feature instead of IssueData |
| 3 | Orchestrator | Same |
| 4 | storage.markdown | Same (it's generic — reads/writes any .md) |
| 5 | **tools.features** | Replaces tools.scrum. ~35 feature-focused tools |
| 6 | transport.stdio | Same |
| 7 | Integration test | Same flow, but creates features not tasks |

---

## 11. AgentOps (Cost & Token Tracking)

Every API call the agent makes is tracked per feature and per cycle. This gives visibility into how much each feature costs to deliver.

### What's tracked

| Metric | Granularity | Where Stored |
|--------|-------------|-------------|
| Input tokens | Per API call, per cycle, per feature | `agentops.cycle_costs[]` in metadata |
| Output tokens | Per API call, per cycle, per feature | `agentops.cycle_costs[]` in metadata |
| Cost (USD) | Per cycle, per feature, per project | `agentops.total_cost_usd` in metadata |
| API calls count | Per cycle, per feature | `agentops.cycle_costs[].api_calls` |
| Model usage | Per model per feature | `agentops.model_usage` in metadata |
| Wall time | Per cycle, per feature | `time.cycle_log[]` in metadata |

### How it works

1. **Before each API call**: The agent records the model, timestamp
2. **After each API call**: The MCP records input_tokens, output_tokens, cost from the response
3. **At cycle end**: Totals are rolled up into `cycle_costs[]`
4. **At feature done**: `total_cost_usd`, `total_input_tokens`, `total_output_tokens` are finalized
5. **At project level**: `get_project_status` aggregates costs across all features

### Cost visibility in the doc

```markdown
## Cost Summary

| Cycle | Input Tokens | Output Tokens | Cost | API Calls | Duration |
|-------|-------------|---------------|------|-----------|----------|
| v1    | 120,000     | 18,000        | $2.34 | 12       | 75 min   |
| v2    | 85,000      | 12,000        | $1.63 | 8        | 45 min   |
| v3    | 40,000      | 8,000         | $0.85 | 5        | 30 min   |
| **Total** | **245,000** | **38,000** | **$4.82** | **25** | **150 min** |
```

### Budget limits

Features can have a cost budget. If exceeded, the agent pauses and asks the human:

```json
{
  "agentops": {
    "budget_usd": 10.00,
    "budget_alert_at": 0.8
  }
}
```

At 80% budget consumed, the agent notifies. At 100%, the agent stops and requests human approval to continue.

---

## 12. Time Management (Estimation)

Every feature has an estimation. The system tracks estimated vs actual to improve future estimates.

### Estimation fields

```json
{
  "estimate": {
    "complexity": "medium",
    "cycles_estimated": 2,
    "cycles_actual": 3
  },
  "time": {
    "started_at": "2026-02-26T10:05:00Z",
    "completed_at": "2026-02-26T14:30:00Z",
    "total_duration_min": 150,
    "cycle_log": [
      {"cycle": 1, "started": "...", "ended": "...", "duration_min": 75},
      {"cycle": 2, "started": "...", "ended": "...", "duration_min": 45},
      {"cycle": 3, "started": "...", "ended": "...", "duration_min": 30}
    ]
  }
}
```

### Complexity levels

| Level | Context Budget | Typical Cycles | Typical Duration |
|-------|---------------|----------------|-----------------|
| `trivial` | < 50K tokens | 1 cycle | < 30 min |
| `small` | 50-100K tokens | 1-2 cycles | 30-60 min |
| `medium` | 100-150K tokens | 2-3 cycles | 1-3 hours |
| `large` | 150-180K tokens | 3-5 cycles | 3-8 hours |
| `split-required` | > 180K tokens | N/A — must split into children | N/A |

### How estimation improves

The system learns from completed features:
1. After each feature is done, `cycles_estimated` vs `cycles_actual` is compared
2. For the same `domain` + `complexity`, historical averages are stored in project memory
3. `create_feature` uses historical data to suggest more accurate estimates
4. `get_project_status` shows estimation accuracy (e.g., "medium features average 2.5 cycles, estimated 2")

### Time tracking in the cycle

- **Cycle start**: When `set_current_feature` or `needs-edits → in-progress` fires, record `started`
- **Cycle end**: When `advance_feature` moves past `tested`, record `ended` and compute `duration_min`
- **Feature complete**: When `human-review → done`, set `completed_at` and `total_duration_min`

---

## 13. Dependencies (Feature Graph)

Every feature can depend on other features. The system enforces dependency order and surfaces blockers.

### Dependency types

| Type | Meaning | Example |
|------|---------|---------|
| `depends_on` | Must be done before this can start | FEAT-002 depends_on FEAT-001 (auth before dashboard) |
| `blocked_by` | Currently blocked, waiting | FEAT-003 blocked_by FEAT-002 (can't test without API) |
| `parent/children` | Feature was split into parts | FEAT-001 has children FEAT-001a, FEAT-001b |

### Enforcement rules

1. **`get_next_feature`** skips features whose `depends_on` features are not yet `done`
2. **`set_current_feature`** rejects if any `blocked_by` feature is not `done`
3. **`advance_feature`** at `tested` gate checks all runtime dependencies are met
4. **`split_feature`** auto-creates `depends_on` chains between children if ordered

### Dependency graph

```
FEAT-000 (project setup)
  ├──→ FEAT-001 (auth)
  │      ├──→ FEAT-001a (login/logout)
  │      └──→ FEAT-001b (refresh tokens) [depends_on: FEAT-001a]
  ├──→ FEAT-002 (database schema)
  └──→ FEAT-003 (API endpoints) [depends_on: FEAT-001, FEAT-002]
         └──→ FEAT-004 (dashboard UI) [depends_on: FEAT-003]
```

The `get_dependency_graph` tool renders this tree for the agent to understand build order.

---

## 14. Smart Assignment (Agent Routing)

Features are assigned to the right agent based on the type of work. Not every feature goes to the same agent.

### Assignment metadata

```json
{
  "assignee": "go-architect",
  "domain": "backend",
  "required_skills": ["go", "auth", "jwt"]
}
```

### Domain → Agent mapping

| Domain | Primary Agent | Skills |
|--------|--------------|--------|
| `backend` | `go-architect` | Go, Fiber, GORM, services, handlers |
| `engine` | `rust-engineer` | Rust, Tonic, Tree-sitter, Tantivy |
| `frontend` | `frontend-dev` | React, TypeScript, Zustand, Tailwind |
| `ui` | `ui-ux-designer` | shadcn/ui, Tailwind, accessibility |
| `database` | `dba` | PostgreSQL, SQLite, Redis, migrations |
| `mobile` | `mobile-dev` | React Native, WatermelonDB |
| `platform` | `platform-engineer` | macOS CGo, Spotlight, Keychain |
| `extension` | `extension-architect` | Extension API, Raycast, VS Code compat |
| `ai` | `ai-engineer` | LLM, RAG, embeddings, vector search |
| `infra` | `devops` | Docker, GCP, CI/CD, monitoring |
| `cross-cutting` | determined by labels | Multi-domain features — assigned by primary label |

### How assignment works

1. **`create_feature`**: Analyzes the feature description + labels to suggest `domain` and `assignee`
2. **`get_next_feature(assignee="go-architect")`**: Returns only features matching that agent's domain
3. **`assign_feature(feature_id, assignee, domain)`**: Manual override by human
4. **Parallel agents**: Multiple agents work in parallel, each pulling features from their domain queue:
   ```
   go-architect    → pulls backend features
   rust-engineer   → pulls engine features
   frontend-dev    → pulls frontend features
   ```
5. **Cross-domain features**: Features touching multiple domains get a primary assignee + `required_skills` that may trigger handoff

### Agent capacity

Each agent has a WIP limit (from `set_wip_limits`). The system tracks per-agent:

| Agent | WIP Limit | In-Progress | Available |
|-------|-----------|-------------|-----------|
| `go-architect` | 2 | 1 (FEAT-001) | 1 |
| `rust-engineer` | 1 | 0 | 1 |
| `frontend-dev` | 2 | 2 (FEAT-004, FEAT-005) | 0 |

`get_next_feature` respects these limits — won't assign to a maxed-out agent.

---

## 15. Git Workflow (Feature Branches)

Every feature maps to a git branch. The cycle happens on the branch. Merge happens only when the human approves.

### Branch strategy

```
main (protected — only merged via approved features)
  │
  ├── feat/FEAT-001-auth
  │     ├── cycle 1: commits for v1
  │     ├── cycle 2: commits for v2 (human edits applied)
  │     └── cycle 3: commits for v3 → human approves → merge to main
  │
  ├── feat/FEAT-002-database  (parallel, different agent)
  │     └── cycle 1: commits → human approves → merge to main
  │
  └── feat/FEAT-003-api  (blocked until FEAT-001 + FEAT-002 merged)
        └── created from main AFTER deps are merged
```

### Rules

1. **`set_current_feature`** → creates branch `feat/{feature-id}` from `main` (or from latest main if returning for a new cycle)
2. **Each cycle**: Agent commits on the feature branch. Commit messages reference the feature ID and cycle number:
   ```
   feat(FEAT-001): v1 — initial auth implementation
   feat(FEAT-001): v2 — fix refresh token rotation
   feat(FEAT-001): v3 — add reuse detection (human review edit)
   ```
3. **`advance_feature` → tested**: Agent pushes the branch
4. **`human-review → done`**: Feature branch is merged to `main` (squash or merge commit, configurable per project)
5. **`human-review → needs-edits`**: Agent rebases on latest `main`, applies edits, commits new cycle
6. **Dependency enforcement**: `set_current_feature` for FEAT-003 fails if FEAT-001's branch isn't merged to `main` yet

### Metadata

```json
{
  "git": {
    "branch": "feat/FEAT-001-auth",
    "base": "main",
    "merge_strategy": "squash",
    "commits": [
      {"cycle": 1, "sha": "abc1234", "message": "feat(FEAT-001): v1 — initial auth"},
      {"cycle": 2, "sha": "def5678", "message": "feat(FEAT-001): v2 — fix rotation"},
      {"cycle": 3, "sha": "ghi9012", "message": "feat(FEAT-001): v3 — reuse detection"}
    ],
    "merged_at": null,
    "merge_sha": null
  }
}
```

### Tools

| Tool | What It Does |
|------|-------------|
| `create_branch` | Create feature branch from main (called by `set_current_feature`) |
| `get_branch_status` | Show branch diff vs main, ahead/behind count |
| `merge_feature` | Merge feature branch to main (called on `done`) |

---

## 16. Notifications (Human Alerts)

The cycle stalls at `human-review` if the human doesn't know to come back. Notifications solve this.

### When notifications fire

| Event | Notification | Priority |
|-------|-------------|----------|
| Feature reaches `human-review` | "FEAT-001 is ready for your review" | High |
| Feature reaches `blocked` | "FEAT-003 is blocked by FEAT-001" | Medium |
| Budget 80% consumed | "FEAT-001 has used $8 of $10 budget" | Medium |
| Budget 100% consumed | "FEAT-001 budget exceeded — agent paused" | High |
| All features done | "Project milestone: all 5 features delivered" | Low |
| Agent error/crash | "go-architect stopped unexpectedly on FEAT-001" | Critical |
| Dependency resolved | "FEAT-001 is done — FEAT-003 is now unblocked" | Low |
| Review queue building up | "3 features waiting for review" | Medium |

### Notification channels

| Channel | When | Config |
|---------|------|--------|
| **MCP response** | Always — returned in tool results | Default, always on |
| **Desktop notification** | When agent reaches human-review | `notifications.desktop: true` |
| **Slack** | When configured | `notifications.slack.webhook: "https://..."` |
| **Email** | For daily digest | `notifications.email.to: "fady@..."` |
| **Sound** | Agent completes or needs attention | `notifications.sound: true` |

### Config in project

```json
{
  "notifications": {
    "desktop": true,
    "sound": true,
    "slack": {
      "webhook": "https://hooks.slack.com/...",
      "channel": "#orchestra-reviews"
    },
    "email": null,
    "digest": "daily"
  }
}
```

### Tools

| Tool | What It Does |
|------|-------------|
| `send_notification` | Send a notification through configured channels |
| `get_notification_config` | Get notification settings for a project |
| `set_notification_config` | Update notification settings |

---

## 17. Parallel Conflict Resolution

Multiple agents working simultaneously on different features can touch the same files. The system detects and resolves conflicts.

### How conflicts happen

```
go-architect working on FEAT-001 (auth):
  → edits middleware/auth.go
  → edits routes/api.go

frontend-dev working on FEAT-004 (dashboard):
  → edits routes/api.go  ← CONFLICT: same file
  → edits components/Dashboard.tsx
```

### Prevention (dependency-based)

The first line of defense is the dependency graph:
- If FEAT-004 `depends_on` FEAT-001, they never run in parallel — FEAT-004 waits
- Only independent features run in parallel
- `create_feature` should set dependencies to minimize overlap

### Detection (file-level)

When dependencies don't fully prevent overlap:

1. **`set_current_feature`** records which files the feature plans to touch (from `## Relevant Files` in the doc)
2. **File lock registry**: The orchestrator maintains a file → feature mapping:
   ```
   routes/api.go       → FEAT-001 (in-progress)
   middleware/auth.go   → FEAT-001 (in-progress)
   components/Dash.tsx  → FEAT-004 (in-progress)
   ```
3. **Conflict detection**: If FEAT-004 tries to touch `routes/api.go` while FEAT-001 holds it:
   - **Soft lock** (default): Agent gets a warning — "routes/api.go is being modified by FEAT-001 (go-architect). Proceed with caution or wait."
   - **Hard lock** (configurable): Agent is blocked until FEAT-001 releases the file

### Resolution (when conflicts happen anyway)

If two agents both modify the same file on different branches:

1. **At merge time** (`merge_feature`): Git detects the conflict
2. **Auto-resolve**: If changes are in different sections of the file → git auto-merges
3. **Manual resolve**: If changes overlap → feature goes to `blocked` state:
   ```
   FEAT-004 blocked: merge conflict with FEAT-001 in routes/api.go
   ```
4. **Agent resolves**: The agent reads both versions, understands both features' intent (from their docs), and resolves the conflict
5. **Human review**: If the conflict is complex, escalate to human

### Metadata

```json
{
  "files_locked": ["routes/api.go", "middleware/auth.go"],
  "conflicts": [
    {
      "file": "routes/api.go",
      "other_feature": "FEAT-004",
      "resolved": true,
      "resolved_by": "go-architect",
      "resolution": "merged both route registrations"
    }
  ]
}
```

### Tools

| Tool | What It Does |
|------|-------------|
| `lock_files` | Register files being modified by a feature |
| `unlock_files` | Release file locks when feature advances past in-progress |
| `check_conflicts` | Check if any files overlap with other in-progress features |
| `resolve_conflict` | Mark a conflict as resolved with explanation |

---

## 18. Automated Quality Gates

The agent self-reports test results, but the system independently verifies. Trust but verify.

### Quality checks at each gate

| Gate | Automated Check | What It Does |
|------|----------------|-------------|
| in-progress → ready-for-testing | **Lint** | Run linter on changed files (golangci-lint, eslint, clippy) |
| in-progress → ready-for-testing | **Type check** | Run type checker (go vet, tsc --noEmit, cargo check) |
| in-progress → ready-for-testing | **Build** | Verify the project compiles |
| ready-for-testing → in-testing | **Test runner** | Run tests independently (not agent-reported) |
| in-testing → tested | **Coverage** | Check coverage meets threshold (configurable) |
| in-testing → tested | **Security scan** | Run security scanner (gosec, npm audit, cargo-audit) |
| tested → human-review | **Diff review** | Generate diff summary for human |

### How it works

```
Agent calls advance_feature("FEAT-001", evidence="Implementation complete")

Before accepting the transition, the system:
  1. Identifies changed files from git diff on the feature branch
  2. Runs lint on changed files
  3. Runs type check on the project
  4. Runs build

If ANY check fails:
  → Transition is REJECTED
  → Agent gets the error output
  → Agent must fix and retry

If all pass:
  → Transition accepted → ready-for-testing
```

### Independent test runner

The agent reports "all 12 tests passing" — but the system runs them too:

```
Agent calls advance_feature("FEAT-001", evidence="All 12 tests passing")

System independently runs:
  $ go test ./... -v -count=1

Compares:
  - Agent reported: 12/12 passing
  - System verified: 12/12 passing ✓

If mismatch:
  → Transition REJECTED
  → Agent gets: "You reported 12/12 but system found 10/12 passing. Failures: ..."
```

### Quality config per project

```json
{
  "quality": {
    "lint": {
      "enabled": true,
      "go": "golangci-lint run",
      "ts": "eslint --ext .ts,.tsx",
      "rust": "cargo clippy -- -D warnings"
    },
    "type_check": {
      "enabled": true,
      "go": "go vet ./...",
      "ts": "tsc --noEmit",
      "rust": "cargo check"
    },
    "test": {
      "enabled": true,
      "go": "go test ./... -v -count=1 -race",
      "ts": "vitest run",
      "rust": "cargo test"
    },
    "coverage": {
      "enabled": true,
      "threshold": 80,
      "go": "go test ./... -coverprofile=coverage.out",
      "ts": "vitest run --coverage"
    },
    "security": {
      "enabled": true,
      "go": "gosec ./...",
      "ts": "npm audit --audit-level=high",
      "rust": "cargo audit"
    }
  }
}
```

### Gate failure flow

```
advance_feature("FEAT-001")
  → lint FAILS: "unused variable on line 42 of auth_handler.go"
  → REJECTED — feature stays in current state
  → Agent receives error → fixes → retries advance_feature
  → lint PASSES, type check PASSES, build PASSES
  → ACCEPTED — feature advances
```

This means the agent can't skip quality. The gates are enforced by the system, not by the agent's self-reporting.

### Tools

| Tool | What It Does |
|------|-------------|
| `run_quality_check` | Run all quality checks on a feature branch |
| `get_quality_config` | Get quality gate configuration |
| `set_quality_config` | Update quality gate configuration |
| `get_quality_report` | Get last quality check results for a feature |

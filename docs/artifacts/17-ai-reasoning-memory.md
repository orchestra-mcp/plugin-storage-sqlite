# AI Reasoning & Project Memory — Orchestra

> The agent doesn't just execute. It thinks.
> Every project starts with understanding. Every feature starts with questions.
> Memory is per-project. The agent gets smarter the longer it works on your project.

---

## 1. The Problem with Current AI Agents

Today's AI agents are **amnesiacs with a to-do list**:
- They lose context between sessions
- They execute blindly without questioning whether the approach is right
- They don't build understanding of YOUR project over time
- They treat every task as isolated, missing the bigger picture
- They don't help you think — they just do what you say

**Orchestra changes this.** The MCP server is the agent's brain:
- **Project memory**: Everything learned persists and is searchable
- **Guided reasoning**: The agent asks the right questions at the right time
- **Progressive understanding**: The more the agent works on your project, the better it gets
- **Context injection**: Before any work, relevant memory is injected into the agent's context

---

## 2. Project Lifecycle: Think → Plan → Build → Learn

```
┌──────────────────────────────────────────────────────────────┐
│                    PROJECT LIFECYCLE                          │
│                                                              │
│  1. INCEPTION          2. DISCOVERY        3. DELIVERY       │
│  ┌────────────┐       ┌────────────┐      ┌────────────┐    │
│  │ What are   │       │ Break into │      │ Build each │    │
│  │ you        │──────→│ features   │─────→│ feature    │    │
│  │ building?  │       │ with agent │      │ (artifact  │    │
│  │ Why?       │       │ reasoning  │      │  16 flow)  │    │
│  └────────────┘       └────────────┘      └────────────┘    │
│        │                    │                    │            │
│        ▼                    ▼                    ▼            │
│  ┌──────────────────────────────────────────────────────┐    │
│  │              PROJECT MEMORY (RAG)                     │    │
│  │                                                       │    │
│  │  Decisions, patterns, code knowledge, user prefs,     │    │
│  │  what worked, what didn't, domain context             │    │
│  └──────────────────────────────────────────────────────┘    │
└──────────────────────────────────────────────────────────────┘
```

---

## 3. Phase 1: Project Inception

When a user starts a new project, the agent doesn't just `create_project`. It runs an **inception interview** — a structured conversation that builds the project's foundational context.

### The Inception Flow

```
User: "I want to build a task management app"

Agent: Before we start building, let me understand what you need.
       I'll ask a few questions to make sure we build the right thing.

  Q1: What problem does this solve?
      → "My team uses spreadsheets to track work, things fall through the cracks"

  Q2: Who will use it? How many people?
      → "My team of 8 developers, maybe expanding to 20"

  Q3: What tools do they use today?
      → "Google Sheets, Slack, GitHub"

  Q4: What's the ONE thing it must do well?
      → "Show me what's blocked and why, instantly"

  Q5: What's the tech context? (existing stack, constraints)
      → "We're a Go shop, everything runs on GCP"

  Q6: What's the timeline/urgency?
      → "Need something working in 2 weeks for the team"

Agent: Based on what you told me, here's what I understand:
       [Summary of the project vision]
       [Proposed first 3 features based on the ONE thing]
       [Technical approach given constraints]

       Does this match what you're thinking? What would you change?
```

### What Gets Stored

The inception creates the **project context document** — the first file in the project:

```markdown
<!-- META
{
  "id": "PROJECT",
  "type": "project-context",
  "created_at": "2026-02-26T10:00:00Z",
  "version": 1
}
META -->

# Project: TaskFlow

## Problem
Team of 8 developers using Google Sheets for task tracking.
Work items fall through the cracks. No visibility into blockers.

## Users
- Primary: 8 developers (expanding to 20)
- Environment: Go shop, GCP infrastructure
- Current tools: Google Sheets, Slack, GitHub

## Core Value
Show what's blocked and why, instantly. If the app does nothing else,
it must do this well.

## Constraints
- Must work within 2 weeks for initial version
- Go backend (team expertise)
- GCP deployment
- Must integrate with GitHub (where code lives) and Slack (where team communicates)

## Technical Decisions
- Go + Fiber for API (team knows Go)
- PostgreSQL on Cloud SQL (GCP native)
- React frontend (team has some React experience)
- GitHub OAuth for login (team already uses GitHub)

## First Features (agent-proposed, user-approved)
1. Blocked items dashboard — show all blocked work with dependency chains
2. GitHub integration — auto-create tasks from issues, sync status
3. Slack notifications — alert when items get blocked/unblocked
```

### Inception Tool: `start_project`

Not just `create_project`. The tool `start_project` triggers the inception interview:

```
1. Agent calls start_project(name="TaskFlow")
   → Returns: first inception question

2. User answers → Agent calls continue_inception(answer="...")
   → Returns: next question (adapted based on previous answers)

3. Repeat until inception is complete (typically 5-8 questions)

4. Agent calls complete_inception(summary="...")
   → Creates project context doc
   → Proposes initial features
   → User approves/modifies
   → Features created as feature docs
```

The inception questions are **adaptive** — they change based on previous answers:
- If the user says "mobile app" → ask about platforms (iOS, Android, both?)
- If the user says "internal tool" → skip marketing/monetization questions
- If the user mentions "real-time" → ask about WebSocket vs polling
- If the user says "already have a backend" → ask about API format, auth

---

## 4. Phase 2: Feature Discovery (Reasoning Before Building)

Before creating a feature doc, the agent **reasons** about the right approach. It doesn't just take orders — it helps the user think.

### The Feature Discovery Flow

```
User: "Add user authentication"

Agent: Let me think about this in the context of your project.
       [Reads project context from memory]
       [Reads existing features from memory]

  Q1: You mentioned GitHub OAuth for login in your project context.
      Should authentication be GitHub OAuth only, or do you also
      need email/password?
      → "GitHub OAuth for now, email/password later"

  Q2: Your team uses Slack. Should users be able to log in with
      their Slack account too?
      → "Good idea, yes, add Slack OAuth as well"

  Q3: For the blocked items dashboard (FEAT-001), users need to see
      each other's work. Should there be any permission levels, or
      can everyone see everything?
      → "Everyone sees everything for now, it's a small team"

  Q4: Session management — how long should sessions last? Do users
      need to stay logged in across devices?
      → "Keep them logged in for 30 days, yes across devices"

Agent: Here's what I'll build:
       - GitHub OAuth login (primary)
       - Slack OAuth login (secondary)
       - No permission levels (everyone sees all)
       - 30-day sessions with refresh tokens, multi-device
       - JWT access tokens (15min) + refresh tokens (30 days)

       This fits in one feature doc (~3000 tokens context).
       Should I create it?
```

### Why This Matters

Without reasoning, the agent would:
1. Create a generic "auth" feature doc
2. Build email/password auth (most common, but wrong for this project)
3. Miss the Slack OAuth opportunity
4. Not connect it to the dashboard feature
5. Make assumptions about sessions that don't match the team's needs

With reasoning, the agent:
1. Reads the project context (knows it's a Go shop using GitHub + Slack)
2. Connects to existing features (dashboard needs user identity)
3. Asks targeted questions based on what it knows
4. Proposes an approach grounded in the project's reality
5. Creates a feature doc that's correct from the start

### Feature Discovery Tool: `plan_feature`

```
1. Agent calls plan_feature(project="taskflow", description="Add user authentication")
   → Returns: first reasoning question + context from project memory

2. User answers → Agent calls continue_planning(answer="...")
   → Returns: next question (informed by project context + previous answers)

3. Repeat until planning is complete (typically 3-6 questions)

4. Agent calls create_feature(project="taskflow", ...)
   → Creates the feature doc with full context, DOD, test cases
```

---

## 5. Project Memory (Per-Project RAG)

Every project has its own memory store. The agent accumulates knowledge over time and uses it to make better decisions.

### What Gets Stored in Memory

| Category | Examples | When Stored |
|----------|---------|-------------|
| **Project context** | Problem, users, constraints, tech stack | Inception |
| **Decisions** | "Chose JWT over sessions because..." | Feature planning |
| **Patterns** | "This codebase uses repository pattern" | Code analysis |
| **User preferences** | "User prefers short functions" | Human review feedback |
| **Domain knowledge** | "A 'sprint' in this project means 1 week" | Conversations |
| **What worked** | "The auth approach from FEAT-002 was approved first try" | Feature completion |
| **What didn't work** | "The caching strategy caused race conditions" | Review feedback |
| **Code structure** | "Auth middleware is at /middleware/auth.go" | Implementation |
| **Dependencies** | "FEAT-003 depends on FEAT-002's auth tokens" | Feature planning |
| **Conventions** | "Error responses use {error, message, details} format" | Code analysis |

### Memory Architecture

```
.projects/{project}/
  project.md                      # Project context (inception output)
  features/                       # Feature docs
  memory/
    decisions/                    # Technical decisions with rationale
      001-auth-approach.md
      002-database-choice.md
    patterns/                     # Code patterns discovered
      repository-pattern.md
      error-handling.md
    knowledge/                    # Domain-specific knowledge
      user-types.md
      workflow-states.md
    feedback/                     # Human review learnings
      prefer-short-functions.md
      always-add-tests.md
  config.md                       # Project config
```

### Memory Tools

| Tool | What It Does |
|------|-------------|
| `save_memory` | Store a piece of knowledge with category + tags |
| `search_memory` | Semantic search across all project memory |
| `get_context` | Get relevant memory for a given task/feature |
| `list_memories` | List all memories by category |
| `update_memory` | Update an existing memory (e.g., decision changed) |
| `delete_memory` | Remove outdated memory |

### How Memory Is Used

**Before every action**, the agent calls `get_context` to retrieve relevant memory:

```
Agent is about to work on FEAT-005 (API rate limiting)

1. Agent calls get_context(feature="FEAT-005", query="rate limiting API")
   → Returns:
     - Project context: "Go shop, GCP, Cloud Run"
     - Decision: "Chose Redis for caching (from FEAT-003)"
     - Pattern: "Middleware pattern used for auth (FEAT-002)"
     - Convention: "All middleware in /middleware/ directory"
     - Knowledge: "Team expects 100 req/s per user max"

2. Agent now has full context before writing any code
   → Uses Redis (already chosen)
   → Follows middleware pattern (already established)
   → Puts file in right directory (already known)
   → Sizes rate limits correctly (already discussed)
```

### RAG Pipeline

```
┌─────────────┐     ┌──────────────┐     ┌─────────────────┐
│ Agent asks   │────→│ Embed query  │────→│ Search LanceDB  │
│ get_context  │     │ (ai.claude)  │     │ (engine.vectors)│
│              │     └──────────────┘     └────────┬────────┘
│              │                                    │
│              │     ┌──────────────┐     ┌────────▼────────┐
│              │←────│ Rank + merge │←────│ Top-K results   │
│              │     │ with feature │     │ from memory     │
│              │     │ context      │     │ files           │
└─────────────┘     └──────────────┘     └─────────────────┘
```

Memory files are embedded and stored in LanceDB (the `engine.vectors` plugin). When the agent asks `get_context`, the query is embedded, similar memories are retrieved, and the most relevant ones are returned within the agent's context budget.

---

## 6. Reasoning at Every Gate

The workflow gates (artifact 16) aren't just mechanical checkpoints. The agent should reason at each one.

### Gate 1: in-progress → ready-for-testing

**Agent asks itself:**
- Did I implement everything in the DOD?
- Did I follow the patterns established in this project? (`search_memory(category="patterns")`)
- Did I handle edge cases the user has flagged before? (`search_memory(category="feedback")`)
- Does my implementation connect correctly with dependent features?

**Agent asks the user (if uncertain):**
- "The DOD says 'rate limiting on login'. Your project context mentions Redis. Should I use the same Redis instance from FEAT-003, or a separate one?"

### Gate 2: in-testing → tested

**Agent asks itself:**
- Do the test cases from the feature doc all pass?
- Did I test edge cases that failed in previous features? (`search_memory(query="test failures")`)
- Is test coverage adequate for the patterns this project uses?

### Gate 3: tested → human-review

**Agent prepares:**
- Summary of what was built and why (referencing project context)
- List of decisions made and trade-offs
- Any concerns or areas where human input would help

### Gate 4: human-review → done (or needs-edits)

**After human review, agent stores:**
- What the human approved → confirms patterns
- What the human changed → stores as preference/feedback
- Why changes were needed → updates project knowledge

```
Human review: "The error messages are too technical. Use user-friendly language."

Agent: save_memory(
  category="feedback",
  content="Use user-friendly error messages, not technical jargon",
  tags=["errors", "ux", "user-preference"]
)
```

Next time the agent writes error messages, `get_context` returns this preference.

---

## 7. Conversation Memory

Not just project-level memory. The agent remembers conversations within a project:

### Session Logs

Every interaction is logged:

```
.projects/{project}/
  memory/
    sessions/
      2026-02-26-inception.md       # The inception conversation
      2026-02-26-feat-001.md        # Working on FEAT-001
      2026-02-27-feat-002-review.md # Review discussion for FEAT-002
```

These are searchable. If the user says "remember when we discussed caching?" — the agent can find it.

### Cross-Session Continuity

When the agent starts a new session:

```
1. Load project context (project.md)
2. Load current feature status (what's in-progress, what's blocked)
3. Load recent session log (what we were doing last time)
4. Load relevant memories for current work
5. Inject all of this into agent context

Result: The agent picks up exactly where you left off.
```

---

## 8. Updated Tool List

### Inception Tools (new)

| Tool | What It Does |
|------|-------------|
| `start_project` | Begin project inception interview |
| `continue_inception` | Answer inception question, get next |
| `complete_inception` | Finalize project context, propose features |

### Discovery Tools (new)

| Tool | What It Does |
|------|-------------|
| `plan_feature` | Start feature discovery conversation |
| `continue_planning` | Answer planning question, get next |
| `reason` | Agent explicitly reasons about a decision (stores rationale) |

### Memory Tools (new)

| Tool | What It Does |
|------|-------------|
| `save_memory` | Store knowledge with category + tags |
| `search_memory` | Semantic search across project memory |
| `get_context` | Get relevant memory for current work |
| `list_memories` | List memories by category |
| `update_memory` | Update existing memory |
| `delete_memory` | Remove outdated memory |
| `save_session` | Save session log |
| `get_session` | Retrieve session log |

### Feature Tools (from artifact 16, unchanged)

35 tools for feature CRUD, workflow, review, dependencies, etc.

### Total Tool Count

| Category | Count |
|----------|-------|
| Inception | 3 |
| Discovery | 3 |
| Memory | 8 |
| Features (artifact 16) | 35 |
| **Total Phase 1** | **49** |

---

## 9. The Agent's Reasoning Loop

Every time the agent works, it follows this loop:

```
┌─────────────────────────────────────────────────┐
│ 1. UNDERSTAND                                    │
│    - Read feature doc                            │
│    - Call get_context() for relevant memory       │
│    - Identify what I know and what I don't       │
│                                                  │
│ 2. QUESTION (if uncertain)                       │
│    - Ask the user targeted questions             │
│    - Questions are informed by project context    │
│    - Never ask what memory already answers        │
│                                                  │
│ 3. PLAN                                          │
│    - Decide approach based on context + answers   │
│    - Call reason() to store the decision          │
│    - Update feature doc with plan                │
│                                                  │
│ 4. EXECUTE                                       │
│    - Implement following project patterns        │
│    - Use conventions from memory                 │
│    - Follow DOD from feature doc                 │
│                                                  │
│ 5. VERIFY                                        │
│    - Run test cases from feature doc             │
│    - Check against past failures from memory     │
│    - Advance through gates with evidence         │
│                                                  │
│ 6. LEARN                                         │
│    - Store what worked                           │
│    - Store what the human changed                │
│    - Update patterns if new ones emerged         │
│    - Save session log                            │
└─────────────────────────────────────────────────┘
```

---

## 10. What This Solves

| Problem | How It's Solved |
|---------|----------------|
| Agent forgets between sessions | Project memory persists, loaded at session start |
| Agent builds wrong thing | Inception interview + feature discovery questions |
| Agent ignores project context | `get_context` injects relevant memory before every action |
| Agent repeats mistakes | Review feedback stored, checked at every gate |
| Agent doesn't follow patterns | Code patterns stored in memory, retrieved during implementation |
| User must specify everything | Agent asks smart questions based on what it already knows |
| No continuity across features | Memory connects features (decisions, patterns, dependencies) |
| Agent can't reason about trade-offs | `reason` tool + project context = informed decisions |
| New session = cold start | Session logs + project context = warm start |

---

## 11. Plugin Architecture

### tools.features (updated from artifact 16)

Gains inception + discovery methods. Same plugin, more tools.

### tools.memory (new plugin)

Handles save/search/get/list/update/delete memory operations. Uses `engine.vectors` for semantic search.

### engine.vectors (Rust, Phase 3)

LanceDB vector store. Embeds memory chunks, enables semantic search. Falls back to keyword search if engine not running.

### Storage flow for memory

```
Agent calls save_memory("User prefers short functions")
  → orchestrator routes to tools.memory
  → tools.memory calls StorageWrite to save .md file in memory/
  → tools.memory calls engine.vectors to embed and index (if available)
  → Memory is now searchable by keyword AND semantically
```

---

## 12. Impact on Build Order

| Step | What | Change from Artifact 15 |
|------|------|------------------------|
| 1 | Proto contract | Add memory messages |
| 2 | Plugin SDK | Same |
| 3 | Orchestrator | Same |
| 4 | storage.markdown | Same |
| 5 | **tools.features** | Add inception + discovery tools (49 tools total) |
| 6 | transport.stdio | Same |
| 7 | Integration | Update e2e test to use `start_project` instead of `create_project` |
| — | **tools.memory** (Phase 2) | Separate plugin for memory CRUD + search |
| — | **engine.vectors** (Phase 3) | LanceDB for semantic search |

Phase 1 gets inception + discovery + basic memory (file-based search).
Phase 2 adds dedicated memory plugin.
Phase 3 adds vector search for semantic memory.

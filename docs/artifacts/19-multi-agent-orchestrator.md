# Multi-Agent Orchestrator & Agent Testing Kit

> Multi-LLM agent orchestration using Google ADK Go patterns.
> Each agent can use a different LLM provider (Claude, OpenAI, Gemini, Ollama/local).
> An orchestrator agent coordinates sub-agents in workflows.
> A testing kit validates agents through MCP tools.

---

## 1. Core Concept

```
═══════════════════════════════════════════════════════════════
              MULTI-AGENT ORCHESTRATION
═══════════════════════════════════════════════════════════════

Today: bridge-claude → Claude Code CLI only
Tomorrow: bridge-{any} → any LLM provider

Today: send_message → one agent, one provider
Tomorrow: run_workflow → N agents, N providers, orchestrated

The SAME plugin system. The SAME QUIC routing.
Just more bridges + an ADK-powered workflow engine.

═══════════════════════════════════════════════════════════════
```

### What This Enables

```
Project: "PR Review Pipeline"
┌─────────────────────────────────────────────────────┐
│  WORKFLOW: sequential                                │
│                                                      │
│  Step 1: Code Analyzer (Gemini 2.5 Pro)             │
│    → Parses code structure, identifies patterns      │
│    → Output: {analysis} in shared state              │
│                                                      │
│  Step 2: Security Scanner (GPT-4o)                  │
│    → Reads {analysis}, scans for vulnerabilities     │
│    → Output: {security_report} in shared state       │
│                                                      │
│  Step 3: Code Reviewer (Claude Sonnet)              │
│    → Reads {analysis} + {security_report}            │
│    → Writes final review with recommendations        │
│                                                      │
│  Step 4: PARALLEL                                    │
│    ├── Doc Generator (Ollama/Llama3 local)          │
│    └── Test Generator (Claude Haiku)                │
└─────────────────────────────────────────────────────┘
```

Each agent:
- Uses its own LLM provider + model
- Has its own agentops account (separate API key, budget)
- Can access Orchestra's MCP tools (search, parse, index)
- Shares state with other agents via ADK session.state
- Has budget enforcement per-agent and per-workflow

---

## 2. Architecture

### Current State (What We Have)

```
┌──────────────┐     ┌──────────────────┐     ┌──────────────┐
│ tools-agentops│     │  tools-sessions  │     │ bridge-claude│
│  8 tools      │◄────│  6 tools         │────►│  5 tools     │
│  accounts     │     │  send_message    │     │  claude CLI   │
│  budgets      │     │  turn history    │     │  spawn/kill   │
└──────────────┘     └──────────────────┘     └──────────────┘
                              │
                     Hardcoded to Claude
                     No multi-provider
                     No multi-agent
```

### Target State (What We Build)

```
┌──────────────┐     ┌──────────────────┐     ┌──────────────┐
│ tools-agentops│     │  tools-sessions  │     │ bridge-claude│
│  + provider   │     │  + provider      │     │  provides_ai:│
│  field on     │     │  dispatch        │     │  ["claude"]  │
│  accounts     │     │                  │     └──────────────┘
└──────┬───────┘     └────────┬─────────┘     ┌──────────────┐
       │                      │               │ bridge-openai │
       │                      │               │  provides_ai:│
       │                      ▼               │  ["openai"]  │
       │              ┌──────────────┐        └──────────────┘
       │              │  orchestrator │        ┌──────────────┐
       └──────────────│  AI routing  │────────│bridge-gemini │
                      │  provider →  │        │  provides_ai:│
                      │  bridge      │        │  ["gemini"]  │
                      └──────┬───────┘        └──────────────┘
                             │                ┌──────────────┐
                             │                │bridge-ollama │
                      ┌──────▼───────┐        │  provides_ai:│
                      │ agent-orch   │        │  ["ollama"]  │
                      │  ADK engine  │        └──────────────┘
                      │  workflows   │
                      │  testing kit │
                      │  20 tools    │
                      └──────────────┘
```

---

## 3. Phase 1: Multi-LLM Bridge Foundation

### 3.1 Account Provider Field

Add `Provider` to the Account struct in `libs/plugin-tools-agentops/`:

```go
type Account struct {
    ID           string            `json:"id"`
    Name         string            `json:"name"`
    Provider     string            `json:"provider"`     // NEW: claude, openai, gemini, ollama
    AuthMethod   string            `json:"auth_method"`
    Config       map[string]string `json:"config"`
    DefaultModel string            `json:"default_model"`
    // ... rest unchanged
}
```

Provider-aware env var generation:

| Provider | Auth Method | Env Vars |
|----------|-------------|----------|
| claude | claude_code | (none — uses native login) |
| claude | setup_token | `CLAUDE_CODE_TOKEN` |
| claude | api_key | `ANTHROPIC_API_KEY` |
| openai | api_key | `OPENAI_API_KEY`, `OPENAI_BASE_URL?` |
| gemini | api_key | `GOOGLE_API_KEY` |
| ollama | custom | `OLLAMA_HOST` (default `http://localhost:11434`) |

### 3.2 New Bridge Plugins

Three new plugins, each following the bridge-claude pattern:

| Plugin | Binary | Provider | SDK/API |
|--------|--------|----------|---------|
| `plugin-bridge-openai` | `bridge-openai` | openai | `github.com/openai/openai-go` |
| `plugin-bridge-gemini` | `bridge-gemini` | gemini | `github.com/google/generative-ai-go/genai` |
| `plugin-bridge-ollama` | `bridge-ollama` | ollama | HTTP REST (`POST /api/chat`) |

Each registers the same 5 tools: `ai_prompt`, `spawn_session`, `kill_session`, `session_status`, `list_active`.

Key difference: OpenAI/Gemini/Ollama bridges use SDK/HTTP directly (no CLI like Claude Code).

```go
// bridge-openai: internal/process.go
func Spawn(ctx context.Context, opts SpawnOptions) (*ChatResponse, error) {
    client := openai.NewClient(option.WithAPIKey(opts.Env["OPENAI_API_KEY"]))
    resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
        Model:    openai.ChatModel(opts.Model),
        Messages: buildMessages(opts),
    })
    return convertResponse(resp), err
}
```

### 3.3 Provider-Aware AI Routing

Add `provider` field to `ToolRequest` proto:

```protobuf
message ToolRequest {
    string tool_name = 1;
    google.protobuf.Struct arguments = 2;
    string caller_plugin = 3;
    string trace_parent = 4;
    string provider = 5;  // NEW: target AI provider
}
```

Router adds `aiRoutes` map: `provider → toolName → plugin`:

```go
// orchestrator/internal/router.go
func (r *Router) RouteAIToolCall(ctx context.Context, provider string, req *ToolRequest) (*ToolResponse, error) {
    routes := r.aiRoutes[provider]
    plugin := routes[req.ToolName]
    return plugin.Client.Send(ctx, req)
}
```

Bridge manifests declare `provides_ai: ["claude"]`, `provides_ai: ["openai"]`, etc.

### 3.4 Sessions Provider Dispatch

`send_message` in tools-sessions resolves provider from account, passes it in tool request:

```go
// tools-sessions: chat.go
env, provider, _ := getAccountEnvAndProvider(ctx, store, session.AccountID)
resp, _ := store.CallToolWithProvider(ctx, "spawn_session", args, provider)
```

---

## 4. Phase 2: ADK-Based Agent Orchestrator

### 4.1 Google ADK Go

[Google ADK](https://google.github.io/adk-docs/) is a Go/Python/TypeScript/Java framework for multi-agent systems:

- **LlmAgent**: Wraps any LLM with instructions + tools
- **SequentialAgent**: Runs sub-agents in order
- **ParallelAgent**: Runs sub-agents concurrently
- **LoopAgent**: Iterates until condition met
- **Custom agents**: Implement `agent.Agent` interface
- **Session state**: Shared key-value store between agents
- **Model-agnostic**: Each agent can use a different LLM

```go
import (
    "google.golang.org/adk/agent/llmagent"
    "google.golang.org/adk/agent/workflowagents/sequentialagent"
)

step1, _ := llmagent.New(llmagent.Config{
    Name: "analyzer", Model: geminiModel,
    Instruction: "Analyze the code...", OutputKey: "analysis",
})
step2, _ := llmagent.New(llmagent.Config{
    Name: "reviewer", Model: claudeModel,
    Instruction: "Review based on {analysis}...",
})
pipeline, _ := sequentialagent.New(sequentialagent.Config{
    AgentConfig: agent.Config{Name: "review-pipeline", SubAgents: []agent.Agent{step1, step2}},
})
```

### 4.2 Plugin: agent-orchestrator

**New directory**: `libs/plugin-agent-orchestrator/`

```
libs/plugin-agent-orchestrator/
  cmd/main.go
  internal/
    plugin.go                    # RegisterTools (20 tools)
    engine/
      engine.go                  # ADK runtime — build agents, run workflows
      models.go                  # Model factory: provider → ADK model.Model
      tools_adapter.go           # Orchestra MCP tools → ADK tool.Tool adapter
      state_bridge.go            # Orchestra storage ↔ ADK session.State
      evaluator.go               # Test case evaluation engine
    store/
      agent_store.go             # Agent definition CRUD
      workflow_store.go          # Workflow definition CRUD
      run_store.go               # Execution history
    tools/
      agent_crud.go              # define/get/list/delete agent (4 tools)
      workflow_crud.go           # define/get/list/delete workflow (4 tools)
      execution.go               # run_agent/run_workflow/get_run_status/list_runs/cancel_run (5 tools)
      oneshot.go                 # list_available_models (1 tool)
      testing.go                 # create_test_suite/run_test_suite/get_test_results/add_test_case/evaluate_response/compare_providers (6 tools)
  orchestra.json
```

### 4.3 Agent Definition Format (YAML)

Stored in `.projects/<project>/agents/<agent-id>.yaml`:

```yaml
id: "code-reviewer"
name: "Code Reviewer"
provider: "claude"
model: "claude-sonnet-4-20250514"
account_id: "ACC-XYZW"
instruction: |
  You are a senior code reviewer. Review for:
  1. Logic errors   2. Security vulnerabilities
  3. Performance    4. Code style
tools:
  - "search"
  - "parse_file"
  - "get_symbols"
max_budget: 5.00
output_key: "review_result"
```

### 4.4 Workflow Definition Format (YAML)

Stored in `.projects/<project>/workflows/<workflow-id>.yaml`:

```yaml
id: "code-review-pipeline"
name: "Code Review Pipeline"
type: "sequential"
agents:
  - ref: "code-analyzer"
    provider: "gemini"
    model: "gemini-2.5-pro"
    account_id: "ACC-GEMINI"
    output_key: "analysis"
    tools: ["parse_file", "get_symbols"]

  - ref: "security-scanner"
    provider: "openai"
    model: "gpt-4o"
    account_id: "ACC-OPENAI"
    instruction: "Scan for security issues based on {analysis}"
    output_key: "security_report"

  - ref: "code-reviewer"
    # Inherits from agent definition

  - type: "parallel"
    agents:
      - ref: "doc-generator"
        provider: "ollama"
        model: "llama3"
      - ref: "test-generator"
        provider: "claude"
        model: "claude-haiku-4-5-20251001"

state:
  target_directory: ""
  review_depth: "thorough"
```

### 4.5 Model Factory

Maps providers to ADK `model.Model` implementations:

```go
func (mf *ModelFactory) CreateModel(ctx context.Context, provider, modelName, accountID string) (model.Model, error) {
    env, _ := mf.getAccountEnv(ctx, accountID)
    switch provider {
    case "claude":
        return anthropic.NewModel(ctx, modelName, env["ANTHROPIC_API_KEY"])
    case "openai":
        return openai.NewModel(ctx, modelName, env["OPENAI_API_KEY"])
    case "gemini":
        return gemini.NewModel(ctx, modelName, &genai.ClientConfig{APIKey: env["GOOGLE_API_KEY"]})
    case "ollama":
        return ollama.NewModel(ctx, modelName, env["OLLAMA_HOST"])
    }
}
```

### 4.6 MCP Tool Adapter

Bridges Orchestra's MCP tools into ADK's tool interface:

```go
type MCPToolAdapter struct {
    name, description string
    schema            map[string]any
    client            StorageClient
}

func (t *MCPToolAdapter) Call(ctx context.Context, args map[string]any) (string, error) {
    resp, _ := callMCPTool(ctx, t.client, t.name, args)
    return extractToolResult(resp), nil
}
```

### 4.7 MCP Tools (20 total)

**Agent CRUD (4):** `define_agent`, `get_agent`, `list_agents`, `delete_agent`

**Workflow CRUD (4):** `define_workflow`, `get_workflow`, `list_workflows`, `delete_workflow`

**Execution (5):** `run_agent`, `run_workflow`, `get_run_status`, `list_runs`, `cancel_run`

**Discovery (1):** `list_available_models`

**Testing (6):** `create_test_suite`, `run_test_suite`, `get_test_results`, `add_test_case`, `evaluate_response`, `compare_providers`

---

## 5. Phase 3: Agent Testing Kit

### 5.1 Test Suite Format

```yaml
id: "code-reviewer-tests"
name: "Code Reviewer Agent Tests"
target:
  agent_id: "code-reviewer"
test_cases:
  - name: "detects SQL injection"
    input:
      prompt: "Review this code"
      state:
        code: 'query := "SELECT * FROM users WHERE id = " + userInput'
    expected:
      contains: ["SQL injection", "parameterized"]
      not_contains: ["looks good"]
      min_length: 100

  - name: "provides line references"
    input:
      prompt: "Review this function"
      state:
        code: "func add(a, b int) int { return a - b }"
    expected:
      regex: ["line \\d+"]
```

### 5.2 Evaluation Engine

Assertion types:
- `contains` — response must include these strings (case-insensitive)
- `not_contains` — response must NOT include these
- `regex` — response must match these patterns
- `min_length` / `max_length` — response length bounds
- `json_path` — for structured JSON output validation

### 5.3 Provider Comparison

`compare_providers` runs the same prompt across all specified providers in parallel and returns:

```
| Provider | Model     | Length | Tokens    | Cost    | Duration |
|----------|-----------|--------|-----------|---------|----------|
| claude   | sonnet    | 1234   | 200/800   | $0.0024 | 3200ms   |
| openai   | gpt-4o    | 987    | 180/650   | $0.0031 | 2800ms   |
| gemini   | 2.5-pro   | 1456   | 210/920   | $0.0018 | 2100ms   |
| ollama   | llama3    | 876    | 150/520   | $0.00   | 4500ms   |
```

### 5.4 Desktop App Integration

All 20 tools are standard MCP tools — accessible from any MCP client (Claude Code, Cursor, desktop app, web dashboard). Desktop-specific features:

- Agent builder UI (YAML editor with validation)
- Test runner with pass/fail indicators
- Provider comparison side-by-side diff
- Run history timeline with cost/token charts
- Live workflow visualization (which agent is running, state flow)

---

## 6. Build Order & Dependencies

```
Phase 1 — Multi-LLM Foundation
  1.1  Add Provider field to Account (agentops)
  1.2  Provider-aware buildEnvVars (agentops)
  1.3  Create bridge-openai plugin
  1.4  Create bridge-gemini plugin (parallel with 1.3)
  1.5  Create bridge-ollama plugin (parallel with 1.3)
  1.6  Add provider to ToolRequest proto + orchestrator AI routing
  1.7  Update tools-sessions for provider dispatch
  1.8  Wire into go.work, Makefile, serve.go

Phase 2 — ADK Agent Orchestrator
  2.1  Scaffold plugin-agent-orchestrator
  2.2  ADK engine: model factory + tool adapter + state bridge
  2.3  Agent/workflow YAML store
  2.4  14 MCP tools (CRUD + execution + discovery)
  2.5  Wire into go.work, Makefile, serve.go

Phase 3 — Testing Kit
  3.1  Test suite format + store
  3.2  Evaluation engine
  3.3  6 testing MCP tools
  3.4  Provider comparison
  3.5  Desktop app UI (frontend)

Phase 4 — Packs (follow-up)
  4.1  Pre-built agent packs (code-reviewer, doc-writer, etc.)
  4.2  Pre-built workflow packs (pr-review-pipeline, etc.)
  4.3  Marketplace integration
```

---

## 7. New Go Dependencies

| Plugin | Dependency | Purpose |
|--------|-----------|---------|
| bridge-openai | `github.com/openai/openai-go` | OpenAI chat completions |
| bridge-gemini | `github.com/google/generative-ai-go/genai` | Gemini API |
| bridge-ollama | (none — plain HTTP) | Ollama REST API |
| agent-orchestrator | `google.golang.org/adk` | ADK runtime, agents, workflows |
| agent-orchestrator | `google.golang.org/adk/agent/llmagent` | LLM agent |
| agent-orchestrator | `google.golang.org/adk/agent/workflowagents/*` | Sequential/Parallel/Loop |

---

## 8. Key Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Separate bridge plugins vs unified | **Separate** | Matches plugin arch, independently optional, each has its own SDK |
| ADK as library vs subprocess | **Library** (Go import) | Direct access to agent.Agent, session.State, workflow patterns |
| Agent definitions | **YAML manifest** (hybrid) | Simple agents in YAML, complex in Go code |
| State management | **ADK session.state** (runtime) + **Orchestra storage** (persistent) | Best of both worlds |
| Tool bridging | **ADK tool.Tool adapter** wrapping MCP | Agents access all 90+ Orchestra tools |
| Testing | **Part of agent-orchestrator** plugin | Keeps deps simple, shares engine code |
| Provider routing | **Orchestrator AI routing** via manifest `provides_ai` | Clean separation, no tool name conflicts |

---

## 9. Critical Files

### Modified
- `libs/plugin-tools-agentops/internal/store/account.go` — Add Provider field
- `libs/plugin-tools-agentops/internal/tools/usage.go` — Provider-aware buildEnvVars
- `libs/plugin-tools-sessions/internal/tools/chat.go` — Provider dispatch in send_message
- `libs/plugin-tools-sessions/internal/storage/client.go` — CallToolWithProvider
- `libs/orchestrator/internal/router.go` — AI routing (aiRoutes map)
- `libs/proto/orchestra/plugin/v1/plugin.proto` — provider field on ToolRequest
- `libs/cli/internal/serve.go` — New optional binaries
- `go.work` — New modules
- `Makefile` — New build targets

### New
- `libs/plugin-bridge-openai/` — OpenAI bridge (5 tools)
- `libs/plugin-bridge-gemini/` — Gemini bridge (5 tools)
- `libs/plugin-bridge-ollama/` — Ollama bridge (5 tools)
- `libs/plugin-agent-orchestrator/` — ADK engine + 20 tools

---

## 10. Sources

- [Google ADK Documentation](https://google.github.io/adk-docs/)
- [ADK for Go Announcement](https://developers.googleblog.com/announcing-the-agent-development-kit-for-go-build-powerful-ai-agents-with-your-favorite-languages/)
- [Multi-Agent Systems in ADK](https://google.github.io/adk-docs/agents/multi-agents/)
- [ADK Go Quickstart](https://google.github.io/adk-docs/get-started/go/)
- [Building Multi-Agent Systems with ADK](https://cloud.google.com/blog/products/ai-machine-learning/build-multi-agentic-systems-using-google-adk)
- [ADK Python GitHub](https://github.com/google/adk-python)

# UI Component Inventory -- Orchestra Reference

> Every frontend package: role, components, stores, hooks, dependencies. Extracted from `orch-ref/packages/@orchestra-mcp/` and `orch-ref/resources/`.

---

## Package Overview

| # | Package | Version | Role | Components | Stores | Hooks | Key Dependencies |
|---|---------|---------|------|-----------|--------|-------|-----------------|
| 1 | @orchestra-mcp/ui | 0.1.0 | Core component library | 30 | 0 | 0 | react (peer) |
| 2 | @orchestra-mcp/ai | 0.2.0 | AI chat interface | 30+ components, 50+ cards | 0 | 13 | ui, editor, icons |
| 3 | @orchestra-mcp/editor | 0.2.0 | Code editing | 7 | 0 | 3 | monaco-editor, theme, ui, widgets |
| 4 | @orchestra-mcp/tasks | 0.1.0 | Task management | 20+ components, 10 widgets | 3 | 8 | ui, editor, icons, widgets, zustand |
| 5 | @orchestra-mcp/theme | 0.2.0 | Theming engine | 1 (ThemePicker) | 0 | 0 | tailwindcss v4, postcss |
| 6 | @orchestra-mcp/icons | 0.1.0 | Icon library | 5 base + 20 code + 13 launcher | 0 | 1 | boxicons, @resvg/resvg-js |
| 7 | @orchestra-mcp/widgets | 0.1.0 | Dashboard widgets | 8 | 0 | 0 | icons, html-to-image |
| 8 | @orchestra-mcp/devtools | 0.1.0 | Developer tools | 12 sessions + 3 UI | 1 | 1 | ui, editor, icons, xterm |
| 9 | @orchestra-mcp/settings | 0.1.0 | Settings UI | 5 | 0 | 0 | react (peer) |
| 10 | @orchestra-mcp/search | 0.1.0 | Search spotlight | 1 | 0 | 0 | react (peer) |
| 11 | @orchestra-mcp/explorer | 0.1.0 | File explorer | 2 | 0 | 0 | react (peer) |
| 12 | @orchestra-mcp/marketplace | 0.1.0 | Extension marketplace | 1 | 0 | 0 | react (peer) |
| 13 | @orchestra-mcp/tracking | 0.1.0 | Time tracking | 2 | 0 | 0 | react (peer) |
| 14 | @orchestra-mcp/account-center | 0.1.0 | Account management | 2 | 0 | 0 | react (peer) |
| 15 | @orchestra-mcp/voice | 0.1.0 | Voice STT/TTS | 0 (engines only) | 0 | 0 | @huggingface/transformers, kokoro-js, idb |
| 16 | @orchestra-mcp/desktop-ui | 0.1.0 | Desktop shell components | 15+ | 4 | 6 | ui, theme, icons, zustand (peers) |
| 17 | @orchestra-mcp/chrome-ui | 0.1.0 | Chrome extension UI | 12+ | 0 | 3 | ui, theme, icons (peers) |
| 18 | @orchestra-mcp/cli | 0.1.0 | Interactive TUI | 3 Ink components | 0 | 0 | ink, commander, chalk, enquirer |
| 19 | @orchestra-mcp/shared | -- | Firebase + Sentry integration | 1 (ErrorBoundary) | 0 | 0 | @sentry/react, firebase |

### App-Level Packages (in `orch-ref/resources/`)

| Package | Version | Role | Pages | Stores | Hooks |
|---------|---------|------|-------|--------|-------|
| @orchestra-mcp/desktop (resources) | 1.0.0 | Vite desktop app | 20+ pages | 6 | 12 |
| @orchestra-mcp/chrome (resources) | 0.1.0 | Chrome extension app | 2 entry points | 0 | 0 |

---

## 1. @orchestra-mcp/ui

### Role
Core shared component library. All primitive UI components live here. Domain-specific components have been extracted to dedicated packages (editor, ai, search, tasks, tracking, explorer, widgets, account-center, settings).

### Package Info
- **Version:** 0.1.0
- **Entry:** `./dist/index.js`
- **Build:** `tsc` + CSS copy + `"use client"` injection
- **Peer Dependencies:** react >=18.0.0, react-dom >=18.0.0

### Components (30 component directories)

#### Layout & Overlay
| Component | File | Sub-components |
|-----------|------|----------------|
| Accordion | `src/Accordion/Accordion.tsx` | -- |
| Modal | `src/Modal/Modal.tsx` | Modal.example.tsx |
| Panel | `src/Panel/Panel.tsx` | -- |
| Popover | `src/Popover/Popover.tsx` | -- |
| Sidebar | `src/Sidebar/Sidebar.tsx` | -- |
| Tabs | `src/Tabs/Tabs.tsx` | -- |
| ContextMenu | `src/ContextMenu/ContextMenu.tsx` | -- |

#### Form Controls
| Component | File | Sub-components |
|-----------|------|----------------|
| Button | `src/Button/Button.tsx` | -- |
| Input | `src/Input/Input.tsx` | -- |
| Checkbox | `src/Checkbox/Checkbox.tsx` | CheckboxCard, TreeCheckbox |
| RadioGroup | `src/RadioGroup/RadioGroup.tsx` | -- |
| Select | `src/Select/Select.tsx` | -- |
| Switch | `src/Switch/Switch.tsx` | -- |
| TagInput | `src/TagInput/TagInput.tsx` | -- |
| PINInput | `src/PINInput/PINInput.tsx` | -- |
| DatePicker | `src/DatePicker/DatePicker.tsx` | DateRangePicker, TimePicker, TimezonePicker, CalendarGrid |
| EmojiPicker | `src/EmojiPicker/EmojiPicker.tsx` | EmojiGrid |
| IconPicker | `src/IconPicker/IconPicker.tsx` | -- |

#### Data Display
| Component | File | Sub-components |
|-----------|------|----------------|
| Alert | `src/Alert/Alert.tsx` | -- |
| Avatar | `src/Avatar/Avatar.tsx` | -- |
| Badge | `src/Badge/Badge.tsx` | -- |
| Breadcrumbs | `src/Breadcrumbs/Breadcrumbs.tsx` | -- |
| Notification | `src/Notification/Notification.tsx` | -- |
| ProgressBar | `src/ProgressBar/ProgressBar.tsx` | -- |
| Tooltip | `src/Tooltip/Tooltip.tsx` | -- |
| EmptyState | `src/EmptyState/EmptyState.tsx` | -- |

#### Loading States
| Component | File | Sub-components |
|-----------|------|----------------|
| Shimmer | `src/Shimmer/` | ShimmerGroup |
| Skeleton | `src/Skeleton/Skeleton.tsx` | -- |

#### Drag & Drop
| Component | File | Sub-components |
|-----------|------|----------------|
| Draggable | `src/Draggable/Draggable.tsx` | DragProvider, DragItem, DropZone, DragBoard, DragGroup |

#### Feedback
| Component | File | Sub-components |
|-----------|------|----------------|
| Toaster | `src/Toaster/` | ToasterProvider, useToaster |

### Exports
All components exported from `src/index.ts` with named exports. Types exported alongside using `export type`.

---

## 2. @orchestra-mcp/ai

### Role
AI chat interface components -- the richest package in the system. Contains the full chat experience: message rendering, streaming, tool result cards, mentions, command palette, code preview, sub-agent views, and 50+ event cards for rendering Claude Code tool outputs.

### Package Info
- **Version:** 0.2.0
- **Entry:** `./src/index.ts`
- **Build:** tsc
- **Test:** vitest
- **Dependencies:** @orchestra-mcp/editor, @orchestra-mcp/icons, @orchestra-mcp/ui
- **Peer Dependencies:** react ^19.0.0, react-dom ^19.0.0

### Components -- Chat Core

| Component | File | Description |
|-----------|------|-------------|
| ChatBox | `src/ChatBox/ChatBox.tsx` | Main orchestrator component |
| ChatHeader | `src/ChatHeader/ChatHeader.tsx` | Session title + controls |
| ChatBody | `src/ChatBody/ChatBody.tsx` | Message list container |
| ChatInput | `src/ChatInput/ChatInput.tsx` | Multi-line input with mentions |
| ChatMessage | `src/ChatMessage/ChatMessage.tsx` | Single message bubble |
| ChatMarkdown | `src/ChatMarkdown/ChatMarkdown.tsx` | Markdown renderer for messages |
| ChatStreamMessage | `src/ChatStreamMessage/ChatStreamMessage.tsx` | Streaming token renderer |
| ChatThinkingMessage | `src/ChatThinkingMessage/ChatThinkingMessage.tsx` | "Thinking..." animation |
| ChatThinkingToggle | `src/ChatThinkingToggle/ChatThinkingToggle.tsx` | Toggle thinking visibility |
| ChatTypingIndicator | `src/ChatTypingIndicator/ChatTypingIndicator.tsx` | Typing dots indicator |
| ChatMessageActions | `src/ChatMessageActions/ChatMessageActions.tsx` | Copy/edit/delete per message |
| ChatMessageContextMenu | `src/ChatMessageContextMenu/ChatMessageContextMenu.tsx` | Right-click message menu |
| ChatQuickActions | `src/ChatQuickActions/ChatQuickActions.tsx` | Quick action buttons |
| ChatStartupPrompts | `src/ChatStartupPrompts/ChatStartupPrompts.tsx` | Initial suggestion cards |
| ChatModelSelector | `src/ChatModelSelector/ChatModelSelector.tsx` | AI model dropdown |
| ChatModeSelector | `src/ChatModeSelector/ChatModeSelector.tsx` | Chat mode picker |

### Components -- Mentions & Commands

| Component | File | Description |
|-----------|------|-------------|
| MentionPopup | `src/MentionPopup/MentionPopup.tsx` | @ mention autocomplete popup |
| MentionToken | `src/MentionToken/MentionToken.tsx` | Inline chip for @mentions |
| MentionTokens | `src/MentionTokens/MentionTokens.tsx` | Mirror overlay for textarea |
| CommandPalette | `src/CommandPalette/CommandPalette.tsx` | / command autocomplete |
| AgentSelectorGrid | `src/AgentSelectorGrid/AgentSelectorGrid.tsx` | Grid popup for agent selection |

### Components -- File & Preview

| Component | File | Description |
|-----------|------|-------------|
| FilePreview | `src/FilePreview/FilePreview.tsx` | Attached file thumbnails |
| FileEditorPage | `src/FileEditorPage/FileEditorPage.tsx` | Full-page Monaco editor with tabs |
| FilesChangedPanel | `src/FilesChangedPanel/FilesChangedPanel.tsx` | Codex-style split-pane file diff |
| PreviewFrame | `src/Preview/PreviewFrame.tsx` | Sandboxed iframe renderer |
| PreviewViewportToolbar | `src/Preview/PreviewViewportToolbar.tsx` | Viewport size controls |

### Components -- Pages

| Component | File | Description |
|-----------|------|-------------|
| SubAgentPage | `src/SubAgentPage/SubAgentPage.tsx` | Full-page agent conversation |
| TerminalPage | `src/TerminalPage/TerminalPage.tsx` | Full-page terminal session |
| TimelineLayout | `src/TimelineLayout/TimelineLayout.tsx` | Timeline visualization |
| TimelineNode | `src/TimelineLayout/TimelineNode.tsx` | Individual timeline entry |

### Components -- Misc

| Component | File | Description |
|-----------|------|-------------|
| BubbleButton | `src/BubbleButton/BubbleButton.tsx` | Floating action bubble |
| ModeToggle | `src/ModeToggle/ModeToggle.tsx` | Window mode switcher |
| PromptCardEditor | `src/PromptCardEditor/PromptCardEditor.tsx` | Editable prompt card UI |
| SessionPickerDialog | `src/SessionPickerDialog/SessionPickerDialog.tsx` | Session selection dialog |

### Event Cards (50+ tool result renderers)

Cards render Claude Code tool call results inside the chat. Each card maps to one or more event types.

| Card | File | Renders |
|------|------|---------|
| CardBase | `src/cards/CardBase.tsx` | Base card wrapper |
| CardErrorBoundary | `src/cards/CardErrorBoundary.tsx` | Error fallback |
| CardRegistry | `src/cards/CardRegistry.ts` | Card type registry |
| EventCardRenderer | `src/cards/EventCardRenderer.tsx` | Event-to-card router |
| McpCardRouter | `src/cards/McpCardRouter.tsx` | MCP tool name router |
| **Tool Cards** | | |
| BashCard | `src/cards/BashCard.tsx` | Shell command results |
| GrepCard | `src/cards/GrepCard.tsx` | Grep search results |
| GlobCard | `src/cards/GlobCard.tsx` | File glob results |
| ReadCard | `src/cards/ReadCard.tsx` | File read results |
| EditCard | `src/cards/EditCard.tsx` | File edit diffs |
| CreateCard | `src/cards/CreateCard.tsx` | File creation results |
| SearchCard | `src/cards/SearchCard.tsx` | Search results |
| **MCP / Orchestra Cards** | | |
| McpCard | `src/cards/McpCard.tsx` | Generic MCP tool result |
| OrchestraCard | `src/cards/OrchestraCard.tsx` | Orchestra-specific tools |
| TaskCard | `src/cards/TaskCard.tsx` | Task create/update |
| TaskDetailCard | `src/cards/TaskDetailCard.tsx` | Task detail view |
| ListCard | `src/cards/ListCard.tsx` | List display |
| EpicCard | `src/cards/EpicCard.tsx` | Epic view |
| StoryCard | `src/cards/StoryCard.tsx` | Story view |
| TeamCard | `src/cards/TeamCard.tsx` | Team view |
| NoteCard | `src/cards/NoteCard.tsx` | Note view |
| **Sprint & Project Cards** | | |
| SprintCard | `src/cards/SprintCard.tsx` | Sprint details |
| BurndownChartCard | `src/cards/BurndownChartCard.tsx` | Burndown chart |
| VelocityCard | `src/cards/VelocityCard.tsx` | Velocity chart |
| StandupCard | `src/cards/StandupCard.tsx` | Standup summary |
| RetrospectiveCard | `src/cards/RetrospectiveCard.tsx` | Retro data |
| WipLimitCard | `src/cards/WipLimitCard.tsx` | WIP limit status |
| ProjectStatusCard | `src/cards/ProjectStatusCard.tsx` | Project status overview |
| ProjectTreeCard | `src/cards/ProjectTreeCard.tsx` | Project tree view |
| DependencyGraphCard | `src/cards/DependencyGraphCard.tsx` | Task dependency graph |
| WorkflowStatusCard | `src/cards/WorkflowStatusCard.tsx` | Workflow state display |
| **PRD Cards** | | |
| PRDSessionCard | `src/cards/PRDSessionCard.tsx` | PRD session |
| PRDPreviewCard | `src/cards/PRDPreviewCard.tsx` | PRD document preview |
| PRDQuestionCard | `src/cards/PRDQuestionCard.tsx` | PRD question form |
| PrdContentCard | `src/cards/PrdContentCard.tsx` | PRD content block |
| PrdPhasesCard | `src/cards/PrdPhasesCard.tsx` | PRD phase list |
| PrdValidationCard | `src/cards/PrdValidationCard.tsx` | PRD validation results |
| **Integration Cards** | | |
| GitHubPRCard | `src/cards/GitHubPRCard.tsx` | GitHub PR details |
| GitHubIssueCard | `src/cards/GitHubIssueCard.tsx` | GitHub issue details |
| CIStatusCard | `src/cards/CIStatusCard.tsx` | CI/CD status |
| FigmaCard | `src/cards/FigmaCard.tsx` | Figma design preview |
| HookEventCard | `src/cards/HookEventCard.tsx` | Webhook event |
| **Misc Cards** | | |
| SubAgentCard | `src/cards/SubAgentCard.tsx` | Sub-agent activity |
| PlanCard | `src/cards/PlanCard.tsx` | Plan/todo list |
| SkillCard | `src/cards/SkillCard.tsx` | Skill activation |
| AgentSwitchCard | `src/cards/AgentSwitchCard.tsx` | Agent hand-off |
| AgentBriefingCard | `src/cards/AgentBriefingCard.tsx` | Agent briefing |
| MemoryCard | `src/cards/MemoryCard.tsx` | Memory search results |
| SessionCard | `src/cards/SessionCard.tsx` | Session details |
| UsageCard | `src/cards/UsageCard.tsx` | Usage stats |
| TodoListCard | `src/cards/TodoListCard.tsx` | Todo list display |
| ConfirmationCard | `src/cards/ConfirmationCard.tsx` | Confirmation prompt |
| ExportCard | `src/cards/ExportCard.tsx` | Export dialog |
| ExportConfigDialog | `src/cards/ExportConfigDialog.tsx` | Export config modal |
| GateCard | `src/cards/GateCard.tsx` | Gate transition |
| QuestionCard | `src/cards/QuestionCard.tsx` | User question form |
| SmartComponentCard | `src/cards/SmartComponentCard.tsx` | Component preview |
| WebSearchCard | `src/cards/WebSearchCard.tsx` | Web search results |
| WebFetchCard | `src/cards/WebFetchCard.tsx` | Web fetch results |
| PreviewCard | `src/cards/PreviewCard.tsx` | Live code preview |
| RawCard | `src/cards/RawCard.tsx` | Raw JSON fallback |

### Hooks (13)

| Hook | File | Description |
|------|------|-------------|
| useDragPosition | `src/hooks/useDragPosition.ts` | Drag offset tracking |
| useAutoResize | `src/hooks/useAutoResize.ts` | Auto-resize textarea |
| useAutoScroll | `src/hooks/useAutoScroll.ts` | Auto-scroll to bottom |
| useStreamRenderer | `src/hooks/useStreamRenderer.ts` | Token-by-token streaming |
| useMentionTrigger | `src/hooks/useMentionTrigger.ts` | Detect @ in input |
| useMentionTokens | `src/hooks/useMentionTokens.ts` | Manage mention chips |
| useMentionSearch | `src/hooks/useMentionSearch.ts` | Search mentions |
| useCommandTrigger | `src/hooks/useCommandTrigger.ts` | Detect / in input |
| useAttachments | `src/hooks/useAttachments.ts` | File attachment state |
| useFileAttachments | `src/hooks/useFileAttachments.ts` | File preview + limits |
| useScreenshot | `src/hooks/useScreenshot.ts` | Screen capture |
| useTodoPin | `src/hooks/useTodoPin.ts` | Pin todo items |
| useTerminalSessions | `src/hooks/useTerminalSessions.ts` | Terminal session management |

### Types
- `src/types/message.ts` -- ChatMessage, AttachedFile, MentionReference, StreamChunk
- `src/types/events.ts` -- 25+ ClaudeCodeEvent types (BashEvent, GrepEvent, McpEvent, etc.)
- `src/types/models.ts` -- AIModel, ChatMode, CHAT_MODES, DEFAULT_MODELS

---

## 3. @orchestra-mcp/editor

### Role
Code editing and document rendering components built on Monaco Editor.

### Package Info
- **Version:** 0.2.0
- **Entry:** `./dist/index.js`
- **Dependencies:** @monaco-editor/react, monaco-editor, @orchestra-mcp/theme, @orchestra-mcp/ui, @orchestra-mcp/widgets, html-to-image
- **Optional:** monaco-languageclient, vscode-ws-jsonrpc (for LSP support)
- **Peer Dependencies:** react ^19.0.0, react-dom ^19.0.0

### Components

| Component | File | Description |
|-----------|------|-------------|
| CodeEditor | `src/CodeEditor/CodeEditor.tsx` | Main Monaco wrapper |
| LegacyCodeEditor | `src/CodeEditor/LegacyCodeEditor.tsx` | Legacy compat editor |
| CodeDiffEditor | `src/CodeEditor/CodeDiffEditor.tsx` | Side-by-side diff |
| MonacoLoader | `src/CodeEditor/MonacoLoader.tsx` | Lazy Monaco loader |
| CodeBlock | `src/CodeBlock/CodeBlock.tsx` | Syntax-highlighted static block |
| GitDiffView | `src/GitDiffView/GitDiffView.tsx` | Git diff renderer |
| MarkdownEditor | `src/MarkdownEditor/MarkdownEditor.tsx` | Markdown edit mode |
| MarkdownRenderer | `src/MarkdownRenderer/MarkdownRenderer.tsx` | Markdown display mode |

### Hooks

| Hook | File | Description |
|------|------|-------------|
| useLsp | `src/CodeEditor/useLsp.ts` | LSP connection management |
| useMonacoKeybindings | `src/CodeEditor/useMonacoKeybindings.ts` | Custom keybinding setup |
| useMonacoTheme | `src/CodeEditor/useMonacoTheme.ts` | Theme syncing for Monaco |

### Utilities
- `src/CodeEditor/language-map.ts` -- File extension to Monaco language mapping
- `src/CodeEditor/theme-bridge.ts` -- Orchestra theme to Monaco theme converter
- `src/CodeEditor/lsp-bridge.ts` -- LSP WebSocket bridge
- `src/CodeEditor/jetbrains-keymap.ts` -- JetBrains keymap support
- `src/CodeEditor/monaco-workers.ts` -- Web worker configuration
- `src/CodeBlock/highlighter.ts` -- Static syntax highlighter
- `src/MarkdownRenderer/parseMarkdown.ts` -- Markdown AST parser
- `src/MarkdownRenderer/inlineFormat.ts` -- Inline formatting helpers
- `src/utils/exportToImage.ts` -- Export editor content as image

---

## 4. @orchestra-mcp/tasks

### Role
Full task management UI: backlog boards, task lists, project dashboards, sprint widgets, detail panels, and MCP integration for live project data.

### Package Info
- **Version:** 0.1.0
- **Entry:** `./dist/index.js`
- **Dependencies:** @orchestra-mcp/editor, @orchestra-mcp/icons, @orchestra-mcp/ui, @orchestra-mcp/widgets, zustand ^5.0.0
- **Peer Dependencies:** react ^19.0.0, react-dom ^19.0.0

### Components -- Views & Boards

| Component | File | Description |
|-----------|------|-------------|
| BacklogView | `src/BacklogView/BacklogView.tsx` | Kanban-style backlog board |
| StatusGrid | `src/StatusGrid/StatusGrid.tsx` | Status count grid |
| TaskList | `src/TaskList/TaskList.tsx` | Flat task list |
| DetailPanel | `src/DetailPanel/DetailPanel.tsx` | Generic detail panel |
| DashboardGrid | `src/Dashboard/DashboardGrid.tsx` | Draggable widget grid |
| ResizeHandle | `src/Dashboard/ResizeHandle.tsx` | Grid resize handle |
| WidgetContextMenu | `src/Dashboard/WidgetContextMenu.tsx` | Widget right-click menu |

### Components -- Project Dashboard

| Component | File | Description |
|-----------|------|-------------|
| ProjectDashboard | `src/ProjectDashboard/ProjectDashboard.tsx` | Full dashboard page |
| DashboardHeader | `src/ProjectDashboard/DashboardHeader.tsx` | Dashboard top bar |
| DashboardSettings | `src/ProjectDashboard/DashboardSettings.tsx` | Dashboard config modal |
| TaskProgressHero | `src/ProjectDashboard/TaskProgressHero.tsx` | Progress hero section |
| WidgetToolbar | `src/ProjectDashboard/WidgetToolbar.tsx` | Widget add toolbar |

### Components -- Sidebars

| Component | File | Description |
|-----------|------|-------------|
| ProjectsSidebar | `src/ProjectsSidebar/ProjectsSidebar.tsx` | Project list sidebar |
| ProjectItem | `src/ProjectsSidebar/ProjectItem.tsx` | Single project item |
| TasksSidebar | `src/TasksSidebar/TasksSidebar.tsx` | Task list sidebar |
| BacklogTree | `src/BacklogTree/BacklogTree.tsx` | Hierarchical backlog tree |
| QuickActions | `src/BacklogTree/QuickActions.tsx` | Quick action buttons |

### Components -- Task Detail

| Component | File | Description |
|-----------|------|-------------|
| TaskDetailPanel | `src/TaskDetailPanel/TaskDetailPanel.tsx` | Full task detail view |
| StatusActions | `src/TaskDetailPanel/StatusActions.tsx` | Status transition buttons |
| LabelsSection | `src/LabelsSection/LabelsSection.tsx` | Label management |
| LinksSection | `src/LinksSection/LinksSection.tsx` | Link management |
| StatusBadge | `src/StatusBadge/StatusBadge.tsx` | Workflow status badge |
| PriorityIcon | `src/PriorityIcon/PriorityIcon.tsx` | Priority indicator |
| LifecycleProgressBar | `src/LifecycleProgressBar/LifecycleProgressBar.tsx` | 13-state lifecycle bar |
| EvidenceLog | `src/EvidenceLog/EvidenceLog.tsx` | Gate evidence log |
| StaleNotification | `src/StaleNotification/StaleNotification.tsx` | Stale data warning |
| ConnectionStatus | `src/ConnectionStatus/ConnectionStatus.tsx` | ConnectionStatusDot, ConnectionStatusBanner |
| TaskFilter | `src/TaskFilter/TaskFilter.tsx` | Filter controls |
| InlineEdit | `src/InlineEdit/InlineEdit.tsx` | EditableText, EditableTextarea, EditableSelect |

### Dashboard Widgets (10)

| Widget | File | Description |
|--------|------|-------------|
| BacklogTreeWidget | `src/widgets/BacklogTreeWidget.tsx` | Tree view of backlog |
| TaskDistributionWidget | `src/widgets/TaskDistributionWidget.tsx` | Status distribution chart |
| RecentActivityWidget | `src/widgets/RecentActivityWidget.tsx` | Activity feed |
| SprintProgressWidget | `src/widgets/SprintProgressWidget.tsx` | Sprint burndown |
| TeamWorkloadWidget | `src/widgets/TeamWorkloadWidget.tsx` | Assignee workload |
| ActiveTasksWidget | `src/widgets/ActiveTasksWidget.tsx` | In-progress tasks |
| ProgressDistributionWidget | `src/widgets/ProgressDistributionWidget.tsx` | Progress pie chart |
| SessionMetricsWidget | `src/widgets/SessionMetricsWidget.tsx` | Session KPIs |
| FeatureAdoptionWidget | `src/widgets/FeatureAdoptionWidget.tsx` | Feature usage stats |
| ProjectStatusWidget | `src/widgets/ProjectStatusWidget.tsx` | Overall project status |
| StatusOverviewWidget | `src/widgets/StatusOverviewWidget.tsx` | Status overview grid |
| TaskStatisticsWidget | `src/widgets/TaskStatisticsWidget.tsx` | Task stats summary |

### Stores (3 Zustand stores)

| Store | File | State |
|-------|------|-------|
| useDashboardStore | `src/stores/useDashboardStore.ts` | WidgetType, WidgetLayout |
| useIntegrationStore | `src/stores/useIntegrationStore.ts` | IntegrationProvider, IntegrationConnection |
| useTaskDetailStore | `src/stores/useTaskDetailStore.ts` | Selected task detail state |

### Hooks (8)

| Hook | File | Description |
|------|------|-------------|
| useProjects | `src/hooks/useProjects.ts` | Project list fetching |
| useProjectTree | `src/hooks/useProjectTree.ts` | Hierarchical project tree |
| useTaskActions | `src/hooks/useTaskActions.ts` | Task CRUD actions |
| useTaskDetail | `src/hooks/useTaskDetail.ts` | Single task detail + evidence |
| useFilteredTree | `src/hooks/useFilteredTree.ts` | Tree filtering + extractAssignees |
| useConnectionStatus | `src/hooks/useConnectionStatus.ts` | MCP connection state |
| useOptimisticTaskActions | `src/hooks/useOptimisticTaskActions.ts` | Optimistic UI updates |
| useDebounce | `src/hooks/useDebounce.ts` | Debounce utility |
| useGridDrag | `src/Dashboard/useGridDrag.ts` | Dashboard grid drag-and-drop |

### MCP Client
- `src/mcp/client.ts` -- `mcpCall()` function for MCP tool invocation

---

## 5. @orchestra-mcp/theme

### Role
Theming engine: 26 color themes in 4 groups (Orchestra, Material, Popular, Classic) and 3 component variants (default, compact, modern). Built on Tailwind CSS v4 Oxide engine with CSS custom properties.

### Package Info
- **Version:** 0.2.0
- **Entry:** Multiple exports (`./`, `./styles`, `./themes`, `./variants`, `./theme-switcher`)
- **Build:** tsup (CJS + ESM + DTS)
- **Dependencies:** tailwindcss ^4.0.0-alpha.30, postcss, @sentry/react, firebase
- **Peer Dependencies:** react >=18.0.0, react-dom >=18.0.0

### Components

| Component | File | Description |
|-----------|------|-------------|
| ThemePicker | `src/ThemePicker/ThemePicker.tsx` | Theme selection UI |

### Exports

#### Color Themes (`themes.ts`)
- `THEMES` -- Array of 26 Theme objects
- `THEME_GROUPS` -- Grouped by Orchestra / Material / Popular / Classic
- `getThemeById(id)`, `getThemesByGroup(group)`, `toCssVariables(theme)`
- Each Theme has: `id`, `label`, `group`, `colors` (ThemeColors), `syntax` (SyntaxColors), `isLight`

#### Component Variants (`variants.ts`)
- 3 variants: `compact`, `modern`, `default`
- `VARIANTS` -- Array of VariantDefinition objects
- `getVariantById(id)`, `isValidVariant(v)`

#### Theme Switcher (`theme-switcher.ts`)
- `initTheme()` -- Initialize from localStorage/system
- `setColorTheme(id)`, `getColorTheme()`
- `setComponentVariant(variant)`, `getComponentVariant()`, `toggleComponentVariant()`
- `onColorThemeChange(cb)`, `onVariantChange(cb)`

---

## 6. @orchestra-mcp/icons

### Role
React SVG icon components organized in three packs: code editor icons, launcher icons, and Boxicons integration. Provides an `IconProvider` for pluggable icon resolution.

### Package Info
- **Version:** 0.1.0
- **Entry:** Multiple exports (`./`, `./code`, `./launcher`, `./boxicons`)
- **Build:** tsc
- **Dependencies:** @resvg/resvg-js, boxicons
- **Peer Dependencies:** react >=18.0.0

### Components

#### Core
| Component | File | Description |
|-----------|------|-------------|
| Icon | `src/Icon.tsx` | Base icon component |
| OrchestraLogo | `src/OrchestraLogo.tsx` | Brand logo |
| UnifiedIcon | `src/UnifiedIcon.tsx` | Cross-pack icon resolver |
| IconProvider | `src/IconProvider.tsx` | Context-based icon pack provider |

#### Code Icons (`src/code/`) -- 20 icons
Add, Check, ChevronDown, ChevronRight, Close, Debug, Error, Extensions, File, Folder, Git, Menu, Output, Problems, Run, Search, Settings, Stop, Terminal, Warning

#### Launcher Icons (`src/launcher/`) -- 13 icons
Calculator, Calendar, Clipboard, Command, Copy, Download, Filter, Heart, Paste, Refresh, Search, Share, Star, Trash, Upload

#### Boxicons (`src/boxicons/`)
| Component | File | Description |
|-----------|------|-------------|
| BoxIcon | `src/boxicons/BoxIcon.tsx` | Boxicons wrapper |
| paths.ts | `src/boxicons/paths.ts` | SVG path data |

### Hook
| Hook | File | Description |
|------|------|-------------|
| useIconResolvers | `src/IconProvider.tsx` | Access registered icon packs |

---

## 7. @orchestra-mcp/widgets

### Role
Reusable dashboard widget containers and chart components. Used by @orchestra-mcp/tasks and @orchestra-mcp/editor.

### Package Info
- **Version:** 0.1.0
- **Entry:** `./dist/index.js`
- **Dependencies:** @orchestra-mcp/icons, html-to-image
- **Peer Dependencies:** react ^19.0.0, react-dom ^19.0.0

### Components

| Component | File | Description |
|-----------|------|-------------|
| Widget | `src/Widget/Widget.tsx` | Generic widget container with actions |
| DataTable | `src/DataTable/DataTable.tsx` | Sortable data table |
| LineChart | `src/Charts/LineChart.tsx` | Line chart (SVG) |
| BarChart | `src/Charts/BarChart.tsx` | Bar chart (SVG) |
| PieChart | `src/Charts/PieChart.tsx` | Pie chart (SVG) |
| AreaChart | `src/Charts/AreaChart.tsx` | Area chart (SVG) |
| DonutChart | `src/Charts/DonutChart.tsx` | Donut chart (SVG) |

### Utilities
- `src/utils/saveFile.ts` -- saveFile, saveBlob, uuidFilename
- `src/utils/exportToImage.ts` -- Widget to image export

---

## 8. @orchestra-mcp/devtools

### Role
Session-based developer tools sidebar. Each "session" is a persistent dev tool instance (terminal, database, SSH, etc.). Supports background WebSocket connections via a worker service.

### Package Info
- **Version:** 0.1.0
- **Entry:** `./src/index.ts`
- **Dependencies:** @orchestra-mcp/editor, @orchestra-mcp/icons, @orchestra-mcp/ui, @xterm/xterm, @xterm/addon-fit
- **Peer Dependencies:** react ^19.0.0, react-dom ^19.0.0, zustand ^5.0.0

### Components -- Shell UI

| Component | File | Description |
|-----------|------|-------------|
| DevToolsSidebar | `src/DevToolsSidebar/DevToolsSidebar.tsx` | Main sidebar |
| DevToolsSessionSidebar | `src/DevToolsSessionSidebar/DevToolsSessionSidebar.tsx` | Session list |
| NewSessionPicker | `src/NewSessionPicker/NewSessionPicker.tsx` | New session type picker |
| SessionContent | `src/SessionContent/SessionContent.tsx` | Dynamic session renderer |

### Components -- Session Providers (9 types)

| Session | File | Description |
|---------|------|-------------|
| TerminalSession | `src/TerminalSession/TerminalSession.tsx` | xterm.js terminal |
| DatabaseSession | `src/DatabaseSession/DatabaseSession.tsx` | SQL query interface |
| SSHSession | `src/SSHSession/SSHSession.tsx` | Remote SSH terminal |
| LogViewerSession | `src/LogViewerSession/LogViewerSession.tsx` | Log stream viewer |
| FileExplorerSession | `src/FileExplorerSession/FileExplorerSession.tsx` | File browser |
| ServiceManagerSession | `src/ServiceManagerSession/ServiceManagerSession.tsx` | Service control |
| DebuggerSession | `src/DebuggerSession/DebuggerSession.tsx` | DAP debugger |
| TestingSession | `src/TestingSession/TestingSession.tsx` | Test runner |
| CloudSession | `src/CloudSession/CloudSession.tsx` | Cloud resource manager |

### Store (1)

| Store | File | State |
|-------|------|-------|
| useDevToolsStore | `src/stores/useDevToolsStore.ts` | Active sessions, selected session |

### Hooks

| Hook | File | Description |
|------|------|-------------|
| useSessionWorker | `src/hooks/useSessionWorker.ts` | WebSocket worker lifecycle |

### Services
- `src/registry/SessionRegistry.ts` -- Plugin-style session type registry
- `src/workers/DevToolsWorkerService.ts` -- Background WebSocket connection manager
- `src/PlaywrightBridge/PlaywrightBridge.ts` -- Browser automation bridge (navigate, click, type, screenshot, evaluate, intercept)
- `src/FileExplorerSession/fileIconService.ts` -- File icon resolver

---

## 9. @orchestra-mcp/settings

### Role
Settings UI components: form renderer, field types, navigation, and validation utilities.

### Package Info
- **Version:** 0.1.0
- **Entry:** `./src/index.ts`
- **Peer Dependencies:** react ^19.0.0, react-dom ^19.0.0

### Components

| Component | File | Description |
|-----------|------|-------------|
| SettingsForm | `src/SettingsForm/SettingsForm.tsx` | Full settings form renderer |
| SettingsGroup | `src/SettingsForm/SettingsGroup.tsx` | Grouped settings section |
| SettingGroupShell | `src/SettingsForm/SettingGroupShell.tsx` | Group layout shell |
| SettingField | `src/SettingsForm/SettingField.tsx` | Individual setting input |
| SettingsNav | `src/SettingsNav/SettingsNav.tsx` | Settings navigation sidebar |

### Utilities
- `src/schema.ts` -- SettingsSchema, SchemaField types
- `src/validate.ts` -- validateSetting() function
- `src/SettingsForm/timezones.ts` -- Timezone list data
- `src/SettingsForm/types.ts` -- Setting, SettingType, SettingValue, SettingsState

---

## 10. @orchestra-mcp/search

### Role
Global search spotlight dialog (Cmd+K / Ctrl+K).

### Package Info
- **Version:** 0.1.0
- **Entry:** `./src/index.ts`
- **Peer Dependencies:** react ^19.0.0, react-dom ^19.0.0

### Components

| Component | File | Description |
|-----------|------|-------------|
| SearchSpotlight | `src/SearchSpotlight/SearchSpotlight.tsx` | Global search modal with categories |

### Exports
- SearchSpotlight, SearchSpotlightProps, SearchResult, SearchCategory

---

## 11. @orchestra-mcp/explorer

### Role
File system tree explorer component.

### Package Info
- **Version:** 0.1.0
- **Entry:** `./src/index.ts`
- **Peer Dependencies:** react ^19.0.0, react-dom ^19.0.0

### Components

| Component | File | Description |
|-----------|------|-------------|
| FileTree | `src/FileTree/FileTree.tsx` | Recursive file tree |
| TreeNode | `src/FileTree/TreeNode.tsx` | Single tree node |

### Exports
- FileTree, FileTreeProps, FileNode

---

## 12. @orchestra-mcp/marketplace

### Role
Extension marketplace card component for browsing and installing extensions.

### Package Info
- **Version:** 0.1.0
- **Entry:** `./src/index.ts`
- **Peer Dependencies:** react ^19.0.0, react-dom ^19.0.0

### Components

| Component | File | Description |
|-----------|------|-------------|
| MarketplaceCard | `src/MarketplaceCard/MarketplaceCard.tsx` | Extension card with install |

### Exports
- MarketplaceCard, MarketplaceCardProps

---

## 13. @orchestra-mcp/tracking

### Role
Time tracking components: timer and calendar view.

### Package Info
- **Version:** 0.1.0
- **Entry:** `./src/index.ts`
- **Peer Dependencies:** react ^19.0.0, react-dom ^19.0.0

### Components

| Component | File | Description |
|-----------|------|-------------|
| Timer | `src/Timer/Timer.tsx` | Start/stop/pause timer |
| FullCalendar | `src/FullCalendar/FullCalendar.tsx` | Calendar event view |

### Exports
- Timer, TimerProps, FullCalendar, FullCalendarProps, CalendarEvent

---

## 14. @orchestra-mcp/account-center

### Role
Account management components: integration cards and user dropdown.

### Package Info
- **Version:** 0.1.0
- **Entry:** `./src/index.ts`
- **Peer Dependencies:** react ^19.0.0, react-dom ^19.0.0

### Components

| Component | File | Description |
|-----------|------|-------------|
| AccountIntegration | `src/AccountIntegration/AccountIntegration.tsx` | Third-party integration card |
| UserDropdown | `src/UserDropdown/UserDropdown.tsx` | User avatar + menu dropdown |

### Exports
- AccountIntegration, AccountIntegrationProps, Integration
- UserDropdown, UserDropdownProps, UserDropdownUser, UserDropdownMenuItem

---

## 15. @orchestra-mcp/voice

### Role
Local speech-to-text and text-to-speech engines. All models run in-browser via WebGPU/WASM. Audio never leaves the device. Uses Hugging Face Transformers.js for STT and Kokoro.js for TTS with IndexedDB for session persistence.

### Package Info
- **Version:** 0.1.0
- **Entry:** `./src/index.ts`
- **Dependencies:** @huggingface/transformers ^3.4.0, kokoro-js ^1.2.0, idb ^8.0.0
- **No peer dependencies** (pure engine, no React)

### Modules (no React components)

| Module | File | Description |
|--------|------|-------------|
| OrchestraSTT | `src/stt-engine.ts` | Speech-to-text engine |
| OrchestraTTS | `src/tts-engine.ts` | Text-to-speech engine |
| db | `src/db.ts` | IndexedDB session storage (saveMeetingSession, listSessions, etc.) |
| types | `src/types.ts` | STTLanguage, TTSLanguage, VoiceSettings, MeetingSession, TranscriptChunk |

---

## 16. @orchestra-mcp/desktop-ui

### Role
Desktop application shell: layout system, chat UI, panel management, notification handling, and Wails IPC integration. This is the component library for the desktop app.

### Package Info
- **Version:** 0.1.0
- **Name:** `@orchestra-mcp/desktop-ui`
- **Entry:** `src/index.ts` (multiple sub-path exports: `./layout`, `./panels`, `./hooks`, `./stores`, `./lib`, `./chat`, `./notes`)
- **Peer Dependencies:** react ^19.0.0, react-dom ^19.0.0, @orchestra-mcp/ui, @orchestra-mcp/theme, @orchestra-mcp/icons, zustand ^5.0.0

### Components -- Layout

| Component | File | Description |
|-----------|------|-------------|
| MainLayout | `src/Layout/MainLayout.tsx` | Top-level app layout |
| PanelLayout | `src/Layout/PanelLayout.tsx` | Resizable panel layout |
| PanelTitleBar | `src/Layout/PanelTitleBar.tsx` | Panel header bar |
| Sidebar | `src/Layout/Sidebar.tsx` | Activity bar sidebar |
| StatusBar | `src/Layout/StatusBar.tsx` | Bottom status bar |
| Topbar | `src/Layout/Topbar.tsx` | Top navigation bar |

### Components -- Chat

| Component | File | Description |
|-----------|------|-------------|
| ChatLayout | `src/Chat/ChatLayout.tsx` | Chat page layout |
| IconSidebar | `src/Chat/IconSidebar.tsx` | Chat icon nav |
| LoadingLogo | `src/Chat/LoadingLogo.tsx` | Loading animation |
| SessionItem | `src/Chat/SessionItem.tsx` | Chat session list item |
| SessionSidebar | `src/Chat/SessionSidebar.tsx` | Session list sidebar |
| WelcomeContent | `src/Chat/WelcomeContent.tsx` | Empty state welcome |

### Components -- Notes

| Component | File | Description |
|-----------|------|-------------|
| NotesView | `src/Notes/NotesView.tsx` | Notes list + editor |
| NotesSidebar | `src/Notes/NotesSidebar.tsx` | Notes navigation |
| NoteItem | `src/Notes/NoteItem.tsx` | Single note item |
| NoteBody | `src/Notes/NoteBody.tsx` | Note content editor |
| NotesPreview | `src/Notes/NotesPreview.tsx` | Note preview card |

### Components -- Panels

| Component | File | Description |
|-----------|------|-------------|
| PanelContainer | `src/Panels/PanelContainer.tsx` | Dynamic panel loader |
| PanelError | `src/Panels/PanelError.tsx` | Panel error fallback |
| PanelNotFound | `src/Panels/PanelNotFound.tsx` | Panel 404 fallback |

### Components -- Root

| Component | File | Description |
|-----------|------|-------------|
| BootstrapProvider | `src/BootstrapProvider.tsx` | App initialization provider |

### Stores (4 Zustand stores)

| Store | File | State |
|-------|------|-------|
| useThemeStore | `src/stores/useThemeStore.ts` | Active theme + variant |
| useWindowStore | `src/stores/useWindowStore.ts` | Window dimensions, sidebar state |
| usePluginStore | `src/stores/usePluginStore.ts` | Loaded plugins, panel registrations |
| useNotificationStore | `src/stores/useNotificationStore.ts` | Notification queue |

### Hooks (6)

| Hook | File | Description |
|------|------|-------------|
| useIPC | `src/hooks/useIPC.ts` | Wails IPC bridge |
| useWebSocket | `src/hooks/useWebSocket.ts` | WebSocket connection manager |
| useThemeSync | `src/hooks/useThemeSync.ts` | Sync theme to Wails backend |
| useNotifications | `src/hooks/useNotifications.ts` | Notification display |
| useBootstrap | `src/hooks/useBootstrap.ts` | App bootstrap sequence |
| usePanels | `src/hooks/usePanels.ts` | Panel open/close/focus |

### Lib (registries)
- `src/lib/panelRegistry.ts` -- registerPanel, registerLazyPanel, PanelRegistration
- `src/lib/sidebarRegistry.ts` -- SidebarViewRegistration, SidebarActionRegistration

---

## 17. @orchestra-mcp/chrome-ui

### Role
Chrome extension UI components: sidebar views, settings panel, tab bar, and platform services for extension storage, desktop discovery, and sync.

### Package Info
- **Version:** 0.1.0
- **Name:** `@orchestra-mcp/chrome-ui`
- **Entry:** `src/index.ts` (sub-path exports: `./sidebar`, `./settings`, `./tabbar`, `./types`, `./hooks`, `./services`)
- **Peer Dependencies:** react ^19.0.0, react-dom ^19.0.0, @orchestra-mcp/ui, @orchestra-mcp/theme, @orchestra-mcp/icons

### Components -- Sidebar

| Component | File | Description |
|-----------|------|-------------|
| Sidebar | `src/Sidebar/Sidebar.tsx` | Main sidebar container |
| SidebarHeader | `src/Sidebar/SidebarHeader.tsx` | Sidebar top header |
| SidebarNav | `src/Sidebar/SidebarNav.tsx` | Navigation items |
| IconNav | `src/Sidebar/IconNav.tsx` | Icon-only sidebar rail |
| TopBar | `src/Sidebar/TopBar.tsx` | Top navigation bar |
| StatusBar | `src/Sidebar/StatusBar.tsx` | Bottom status |
| ViewBody | `src/Sidebar/ViewBody.tsx` | View content area |
| ViewTitle | `src/Sidebar/ViewTitle.tsx` | View title bar |
| ChromeSidebar | `src/Sidebar/ChromeSidebar.stories.tsx` | (Storybook only) |

### Components -- Settings

| Component | File | Description |
|-----------|------|-------------|
| SettingsPanel | `src/Settings/SettingsPanel.tsx` | Extension settings page |
| SettingsGroup | `src/Settings/SettingsGroup.tsx` | Grouped settings |
| SettingInput | `src/Settings/SettingInput.tsx` | Setting input control |

### Components -- Tab Bar

| Component | File | Description |
|-----------|------|-------------|
| TabBar | `src/TabBar/TabBar.tsx` | Tab switcher bar |

### Hooks (3)

| Hook | File | Description |
|------|------|-------------|
| useSidebarViews | `src/hooks/useSidebarViews.ts` | Sidebar view navigation |
| useTabManager | `src/hooks/useTabManager.ts` | Tab lifecycle |
| useDesktopConnection | `src/hooks/useDesktopConnection.ts` | Desktop app discovery |

### Services (8)

| Service | File | Description |
|---------|------|-------------|
| TabManager | `src/services/TabManager.ts` | Tab state management |
| SidebarViewRegistry | `src/services/SidebarViewRegistry.ts` | View registration |
| DesktopDiscovery | `src/services/DesktopDiscovery.ts` | Find local desktop app |
| StorageWrapper | `src/services/StorageWrapper.ts` | chrome.storage abstraction |
| WebSocketClient | `src/services/WebSocketClient.ts` | WS connection to desktop |
| SyncProtocol | `src/services/SyncProtocol.ts` | Sync message protocol |
| SettingsSyncService | `src/services/SettingsSync.ts` | Settings sync |
| ThemeSyncService | `src/services/ThemeSyncService.ts` | Theme sync |

### Types
- `src/types/settings.ts` -- Extension settings types
- `src/types/sidebar.ts` -- Sidebar view types
- `src/types/tabs.ts` -- Tab types
- `src/Sidebar/defaultViews.ts` -- Default sidebar view definitions

---

## 18. @orchestra-mcp/cli

### Role
Interactive command-line TUI built with Ink (React for terminals). Provides plugin browsing, theme switching, and Storybook management.

### Package Info
- **Version:** 0.1.0
- **Binary:** `orchestra` (via `./bin/orchestra.js`)
- **Dependencies:** chalk, commander, enquirer, ink, ink-text-input, ink-select-input, ink-spinner, ora, react
- **Peer Dependencies:** react ^18.0.0 || ^19.0.0

### Commands

| Command | File | Description |
|---------|------|-------------|
| welcome | `src/commands/welcome.tsx` | ASCII art welcome screen |
| plugin | `src/commands/plugin.tsx` | Plugin browser TUI |
| theme | `src/commands/theme.tsx` | Theme switcher TUI |
| storybook | `src/commands/storybook.ts` | Start Storybook server |

### Ink Components

| Component | File | Description |
|-----------|------|-------------|
| WelcomeScreen | `src/components/WelcomeScreen.tsx` | ASCII art + menu |
| PluginBrowser | `src/components/PluginBrowser.tsx` | Plugin list browser |
| ThemeSwitcher | `src/components/ThemeSwitcher.tsx` | Theme preview + switch |

---

## 19. @orchestra-mcp/shared

### Role
Cross-cutting infrastructure: Firebase analytics/messaging and Sentry error tracking. Used by all apps via @orchestra-mcp/theme (which depends on @sentry/react and firebase).

### Package Info
- **No package.json** (raw `src/` directory, no standalone module)
- **Located at:** `orch-ref/packages/@orchestra-mcp/shared/src/`

### Modules

#### Firebase (`src/firebase/`)
| Module | File | Exports |
|--------|------|---------|
| config | `config.ts` | initializeFirebase, getFirebaseApp, getFirebaseAnalytics, getFirebaseMessaging |
| analytics | `analytics.ts` | trackEvent, trackPageView, trackLogin, trackSignup, trackProjectCreated, trackTaskCompleted, trackChatMessage, trackAgentInvoked, trackError, trackSearch, trackSettingsChanged, setAnalyticsUserId |
| messaging | `messaging.ts` | requestNotificationPermission, listenForMessages, showNotification, handlePushNotification |

#### Sentry (`src/sentry/`)
| Module | File | Exports |
|--------|------|---------|
| config | `config.ts` | initializeSentry, isSentryInitialized, getSentryConfig |
| capture | `capture.ts` | captureException, captureMessage, captureErrorWithContext, setUser, setTag, addBreadcrumb, startTransaction, logApiError, logNavigationError, logComponentError |
| ErrorBoundary | `ErrorBoundary.tsx` | ErrorBoundary, withErrorBoundary |

---

## App Packages (resources/)

### @orchestra-mcp/desktop (resources/desktop)

The Vite-built desktop application. Imports from nearly every package.

**Dependencies:** @orchestra-mcp/ai, @orchestra-mcp/desktop-ui, @orchestra-mcp/devtools, @orchestra-mcp/editor, @orchestra-mcp/icons, @orchestra-mcp/search, @orchestra-mcp/tasks, @orchestra-mcp/theme, @orchestra-mcp/ui, @orchestra-mcp/widgets, @xterm/xterm, react-router-dom, zustand

#### Pages (20+)

| Page | File | Description |
|------|------|-------------|
| App | `src/App.tsx` | Root app with routes |
| MainChatPage | `src/pages/MainChatPage.tsx` | Primary chat interface |
| ChatPage | `src/pages/ChatPage.tsx` | Chat page variant |
| SessionChat | `src/pages/SessionChat.tsx` | Session-based chat |
| SessionPanel | `src/pages/SessionPanel.tsx` | Session panel view |
| BubblePage | `src/pages/BubblePage.tsx` | Floating bubble mode |
| NotesPage | `src/pages/NotesPage.tsx` | Notes editor |
| SettingsPage | `src/pages/SettingsPage.tsx` | Settings |
| SyncSettingsPage | `src/pages/SyncSettingsPage.tsx` | Sync configuration |
| ProjectsView | `src/pages/ProjectsView.tsx` | Project list |
| ProjectStatusPage | `src/pages/ProjectStatusPage.tsx` | Project dashboard |
| CodeViewerPage | `src/pages/CodeViewerPage.tsx` | Code file viewer |
| MarkdownViewerPage | `src/pages/MarkdownViewerPage.tsx` | Markdown viewer |
| PromptPage | `src/pages/PromptPage.tsx` | Prompt editor |
| ScreenshotPage | `src/pages/ScreenshotPage.tsx` | Screenshot capture |
| ComponentBrowserPage | `src/pages/ComponentBrowserPage.tsx` | Component explorer |
| ComponentEditorPage | `src/pages/ComponentEditorPage.tsx` | Component editor |
| IntegrationsPage | `src/pages/IntegrationsPage.tsx` | Integration management |
| OnboardingPage | `src/pages/OnboardingPage.tsx` | First-run onboarding |
| LoginPage | `src/pages/LoginPage.tsx` | Authentication |
| SpiritPage | `src/pages/SpiritPage.tsx` | Spirit mode |
| TeamsPage | `src/pages/TeamsPage.tsx` | Team list |
| TeamDetailPage | `src/pages/TeamDetailPage.tsx` | Team detail |
| PanelPage | `src/pages/PanelPage.tsx` | Dynamic panel page |

#### Integration Settings Pages
AIProvidersSettings, AppleSettings, DiscordSettings, FigmaSettings, GitHubSettings, JiraSettings, LinearSettings, NotionSettings, SlackSettings

#### Onboarding Steps
SplashStep, WelcomeStep, AuthStep, ThemeStep, VoiceStep, PermissionsStep, ServicesStep, DoneStep

#### Voice Components
VoiceChatButton, VoiceInputOverlay, VoiceSettingsPanel, WaveformVisualizer

#### App Stores (6 Zustand stores)
| Store | File |
|-------|------|
| authStore | `src/stores/authStore.ts` |
| chatStore | `src/stores/chatStore.ts` |
| chatUiStore | `src/stores/chatUiStore.ts` |
| notesStore | `src/stores/notesStore.ts` |
| workspaceStore | `src/stores/workspaceStore.ts` |
| agentRegistryStore | `src/stores/agentRegistryStore.ts` |
| syncMiddleware | `src/stores/syncMiddleware.ts` |

#### App Hooks (12)
| Hook | File |
|------|------|
| useAISync | `src/hooks/useAISync.ts` |
| useAgentRegistry | `src/hooks/useAgentRegistry.ts` |
| useAnalytics | `src/hooks/useAnalytics.ts` |
| useBrowserAwareness | `src/hooks/useBrowserAwareness.ts` |
| useCommandItems | `src/hooks/useCommandItems.ts` |
| useFCM | `src/hooks/useFCM.ts` |
| useNotificationTTS | `src/hooks/useNotificationTTS.ts` |
| useScreenAwareness | `src/hooks/useScreenAwareness.ts` |
| useSessionMetrics | `src/hooks/useSessionMetrics.ts` |
| useSyncEntity | `src/hooks/useSyncEntity.ts` |
| useSyncStores | `src/hooks/useSyncStores.ts` |
| useVoiceChat | `src/hooks/useVoiceChat.ts` |
| useVoiceSettings | `src/hooks/useVoiceSettings.ts` |

#### Providers
ThemeProvider, SettingsProvider, WebSocketProvider, OnboardingContext

### @orchestra-mcp/chrome (resources/chrome)

The Chrome extension build. Uses Vite with multiple entry points (sidebar panel, service worker, content scripts, offscreen page, preview page).

**Dependencies:** @orchestra-mcp/chrome-ui, @orchestra-mcp/theme, @orchestra-mcp/ui, react, react-dom

#### Entry Points

| Entry | File | Description |
|-------|------|-------------|
| Service Worker | `src/background/service-worker.ts` | Background script |
| Content Script | `src/content/content.ts` | Page injection |
| Page Extractor | `src/content/extract.ts` | Page content extraction |
| Meet Injector | `src/content/meet-injector.ts` | Google Meet integration |
| Example Plugin | `src/content/example-plugin.ts` | Plugin example |
| Offscreen | `src/offscreen/offscreen.ts` | Offscreen document |
| Preview App | `src/preview/PreviewApp.tsx` | Code preview iframe |
| Preview Viewport | `src/preview/PreviewViewportToolbar.tsx` | Preview controls |
| buildSrcdoc | `src/preview/buildSrcdoc.ts` | HTML doc builder |

---

## Inter-Package Dependency Graph

```
resources/desktop (Vite app)
  +-- @orchestra-mcp/desktop-ui
  |     +-- @orchestra-mcp/ui (peer)
  |     +-- @orchestra-mcp/theme (peer)
  |     +-- @orchestra-mcp/icons (peer)
  |     +-- zustand (peer)
  +-- @orchestra-mcp/ai
  |     +-- @orchestra-mcp/editor
  |     |     +-- @monaco-editor/react
  |     |     +-- @orchestra-mcp/theme
  |     |     +-- @orchestra-mcp/ui
  |     |     +-- @orchestra-mcp/widgets
  |     |           +-- @orchestra-mcp/icons
  |     +-- @orchestra-mcp/icons
  |     +-- @orchestra-mcp/ui
  +-- @orchestra-mcp/tasks
  |     +-- @orchestra-mcp/editor
  |     +-- @orchestra-mcp/icons
  |     +-- @orchestra-mcp/ui
  |     +-- @orchestra-mcp/widgets
  |     +-- zustand
  +-- @orchestra-mcp/devtools
  |     +-- @orchestra-mcp/editor
  |     +-- @orchestra-mcp/icons
  |     +-- @orchestra-mcp/ui
  |     +-- @xterm/xterm
  +-- @orchestra-mcp/search
  +-- @orchestra-mcp/widgets
  +-- @orchestra-mcp/theme
  +-- @orchestra-mcp/icons

resources/chrome (Vite extension)
  +-- @orchestra-mcp/chrome-ui
  |     +-- @orchestra-mcp/ui (peer)
  |     +-- @orchestra-mcp/theme (peer)
  |     +-- @orchestra-mcp/icons (peer)
  +-- @orchestra-mcp/theme
  +-- @orchestra-mcp/ui

Leaf packages (no internal deps):
  @orchestra-mcp/ui
  @orchestra-mcp/theme
  @orchestra-mcp/icons (depends on boxicons only)
  @orchestra-mcp/settings
  @orchestra-mcp/search
  @orchestra-mcp/explorer
  @orchestra-mcp/marketplace
  @orchestra-mcp/tracking
  @orchestra-mcp/account-center
  @orchestra-mcp/voice
  @orchestra-mcp/cli (Ink/Node only)
  @orchestra-mcp/shared (no package.json)
```

---

## Component Count Summary

| Package | React Components | Zustand Stores | Custom Hooks | Services/Utilities |
|---------|-----------------|----------------|--------------|-------------------|
| ui | 30 | 0 | 0 | 0 |
| ai | 30 + 50 cards = 80 | 0 | 13 | 3 type modules |
| editor | 8 | 0 | 3 | 8 utilities |
| tasks | 25 + 12 widgets = 37 | 3 | 9 | 1 MCP client |
| theme | 1 | 0 | 0 | 3 modules |
| icons | 5 + 20 + 13 = 38 | 0 | 1 | 2 scripts |
| widgets | 8 | 0 | 0 | 2 utilities |
| devtools | 13 | 1 | 1 | 3 services |
| settings | 5 | 0 | 0 | 2 utilities |
| search | 1 | 0 | 0 | 0 |
| explorer | 2 | 0 | 0 | 0 |
| marketplace | 1 | 0 | 0 | 0 |
| tracking | 2 | 0 | 0 | 0 |
| account-center | 2 | 0 | 0 | 0 |
| voice | 0 | 0 | 0 | 3 engines |
| desktop-ui | 17 | 4 | 6 | 2 registries |
| chrome-ui | 12 | 0 | 3 | 8 services |
| cli | 3 Ink components | 0 | 0 | 4 commands |
| shared | 1 | 0 | 0 | 6 modules |
| **TOTAL** | **~267 components** | **8 stores** | **36 hooks** | **~47 services** |

### App-Level (resources/)

| App | Pages | Stores | Hooks | Providers |
|-----|-------|--------|-------|-----------|
| desktop | 24+ pages | 6 | 13 | 4 |
| chrome | 2 entry points | 0 | 0 | 0 |

---

## Testing Coverage

All packages with tests use **vitest** + **@testing-library/react**. Storybook stories exist for most visual components (`.stories.tsx` files).

| Package | Test Files | Stories |
|---------|-----------|---------|
| ui | 30 `.test.tsx` | 30 `.stories.tsx` |
| ai | 3 `.test.tsx` | 20+ `.stories.tsx` |
| editor | 4 `.test.tsx` | 4 `.stories.tsx` |
| tasks | 15 `.test.tsx` | 8 `.stories.tsx` |
| settings | 3 `.test.tsx` | 3 `.stories.tsx` |
| search | 1 `.test.tsx` | 1 `.stories.tsx` |
| explorer | 1 `.test.tsx` | 1 `.stories.tsx` |
| marketplace | 1 `.test.tsx` | 1 `.stories.tsx` |
| tracking | 2 `.test.tsx` | 2 `.stories.tsx` |
| account-center | 2 `.test.tsx` | 2 `.stories.tsx` |
| desktop-ui | 12 `.test.ts(x)` | 8 `.stories.tsx` |
| chrome-ui | 15 `.test.ts(x)` | 10 `.stories.tsx` |
| icons | 4 `.test.tsx` | 2 `.stories.tsx` |
| widgets | 3 `.test.tsx` | 3 `.stories.tsx` |

# MCP Tool Catalog -- Orchestra Reference

> Extracted from `orch-ref/app/tools/`. Every MCP tool with full parameter detail.
> Total: **186 tools** across **34 categories** (plus 3 resources and 3 prompts).

---

## 1. Project Management (app/tools/project.go)

### 1. list_projects
- **Description:** List all projects
- **Parameters:** None
- **Returns:** Array of project summaries (project name, slug, status, description, timestamps)
- **Source:** app/tools/project.go

### 2. create_project
- **Description:** Create a new project with PRD
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | name | string | Yes | Project name |
  | description | string | No | Project description |

- **Returns:** Created project slug, key, and status
- **Source:** app/tools/project.go

### 3. delete_project
- **Description:** Delete a project and all its data
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |

- **Returns:** Deletion confirmation with slug
- **Source:** app/tools/project.go

### 4. get_project_status
- **Description:** Get project status and summary
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |

- **Returns:** Project metadata, epics/stories/tasks summaries, workflow stats (total, completed, blocked, in-progress, completion percentage)
- **Source:** app/tools/project.go

### 5. get_project_tree
- **Description:** Get full project tree: epics -> stories -> tasks
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |

- **Returns:** Nested tree structure with all epics, their stories, and their tasks
- **Source:** app/tools/project_tree.go

### 6. read_prd
- **Description:** Read project PRD document
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |

- **Returns:** PRD markdown content as text
- **Source:** app/tools/project.go

### 7. write_prd
- **Description:** Write/update project PRD document
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | content | string | Yes | PRD markdown |

- **Returns:** Confirmation text
- **Source:** app/tools/project.go

---

## 2. Epic Management (app/tools/epic.go, epic_handlers.go)

### 8. list_epics
- **Description:** List epics in a project
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |

- **Returns:** Array of epic issue data
- **Source:** app/tools/epic.go

### 9. create_epic
- **Description:** Create a new epic
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | title | string | Yes | Epic title |
  | description | string | No | Epic description |
  | priority | string | No | Priority: low, medium, high, critical |

- **Returns:** Created epic issue data
- **Source:** app/tools/epic.go

### 10. get_epic
- **Description:** Get epic details
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | epic_id | string | Yes | Epic ID |

- **Returns:** Epic issue data with description
- **Source:** app/tools/epic.go

### 11. update_epic
- **Description:** Update epic fields
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | epic_id | string | Yes | Epic ID |
  | title | string | No | New title |
  | description | string | No | New description |
  | status | string | No | New status |
  | priority | string | No | New priority |

- **Returns:** Updated epic issue data
- **Source:** app/tools/epic.go

### 12. delete_epic
- **Description:** Delete epic and all children
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | epic_id | string | Yes | Epic ID |

- **Returns:** Deletion confirmation text
- **Source:** app/tools/epic.go

---

## 3. Story Management (app/tools/story.go, story_handlers.go)

### 13. list_stories
- **Description:** List stories in an epic
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | epic_id | string | Yes | Epic ID |

- **Returns:** Array of story issue data
- **Source:** app/tools/story.go

### 14. create_story
- **Description:** Create a story under an epic
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | epic_id | string | Yes | Epic ID |
  | title | string | Yes | Story title |
  | user_story | string | Yes | As a... I want... So that... |
  | priority | string | No | Priority level |

- **Returns:** Created story issue data
- **Source:** app/tools/story.go

### 15. get_story
- **Description:** Get story with child tasks
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | epic_id | string | Yes | Epic ID |
  | story_id | string | Yes | Story ID |

- **Returns:** Story issue data with children array
- **Source:** app/tools/story.go

### 16. update_story
- **Description:** Update story fields
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | epic_id | string | Yes | Epic ID |
  | story_id | string | Yes | Story ID |
  | title | string | No | New title |
  | description | string | No | New description |
  | status | string | No | New status |
  | priority | string | No | New priority |

- **Returns:** Updated story issue data
- **Source:** app/tools/story.go

### 17. delete_story
- **Description:** Delete story and all tasks
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | epic_id | string | Yes | Epic ID |
  | story_id | string | Yes | Story ID |

- **Returns:** Deletion confirmation text
- **Source:** app/tools/story.go

---

## 4. Task Management (app/tools/task.go, task_update.go)

### 18. list_tasks
- **Description:** List tasks in a story
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | epic_id | string | Yes | Epic ID |
  | story_id | string | Yes | Story ID |

- **Returns:** Array of task issue data
- **Source:** app/tools/task.go

### 19. create_task
- **Description:** Create a task/bug/hotfix under a story
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | epic_id | string | Yes | Epic ID |
  | story_id | string | Yes | Story ID |
  | title | string | Yes | Task title |
  | type | string | Yes | Type: task, bug, hotfix |
  | description | string | No | Task description |
  | priority | string | No | Priority level |

- **Returns:** Created task issue data
- **Source:** app/tools/task.go

### 20. get_task
- **Description:** Get task details
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | epic_id | string | Yes | Epic ID |
  | story_id | string | Yes | Story ID |
  | task_id | string | Yes | Task ID |

- **Returns:** Task issue data
- **Source:** app/tools/task.go

### 21. update_task
- **Description:** Update task with workflow validation
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | epic_id | string | Yes | Epic ID |
  | story_id | string | Yes | Story ID |
  | task_id | string | Yes | Task ID |
  | title | string | No | New title |
  | description | string | No | New description |
  | status | string | No | New status (validated against workflow) |
  | priority | string | No | New priority |

- **Returns:** Updated task issue data
- **Source:** app/tools/task_update.go

### 22. delete_task
- **Description:** Delete a task
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | epic_id | string | Yes | Epic ID |
  | story_id | string | Yes | Story ID |
  | task_id | string | Yes | Task ID |

- **Returns:** Deletion confirmation text
- **Source:** app/tools/task_update.go

---

## 5. Workflow Management (app/tools/workflow.go, workflow_cascade.go)

### 23. get_next_task
- **Description:** Get highest priority actionable task. Optional filters: epic_id, story_id, assignee, label to scope results for parallel agent work.
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | epic_id | string | No | Only tasks under this epic |
  | story_id | string | No | Only tasks under this story |
  | assignee | string | No | Only tasks assigned to this person |
  | label | string | No | Only tasks with this label |

- **Returns:** Highest priority actionable task issue data (sorted by type priority then status priority)
- **Source:** app/tools/workflow.go

### 24. set_current_task
- **Description:** Set task to in-progress, cascade parents
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | epic_id | string | Yes | Epic ID |
  | story_id | string | Yes | Story ID |
  | task_id | string | Yes | Task ID |

- **Returns:** Updated task (with WIP limit enforcement, parent cascading to in-progress)
- **Source:** app/tools/workflow_cascade.go

### 25. complete_task
- **Description:** Complete task, cascade done if all siblings done
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | epic_id | string | Yes | Epic ID |
  | story_id | string | Yes | Story ID |
  | task_id | string | Yes | Task ID |

- **Returns:** Updated task (transitions to ready-for-testing, cascades parent completion)
- **Source:** app/tools/workflow_cascade.go

### 26. search
- **Description:** Search issues by text, optional type filter
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | query | string | Yes | Search text |
  | type | string | No | Type filter: epic, story, task, bug, hotfix |

- **Returns:** Array of matching issue data
- **Source:** app/tools/workflow.go

### 27. get_workflow_status
- **Description:** Get workflow stats: counts, blocked, completion %
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |

- **Returns:** Total/done counts, completion percentage, by-status/by-type breakdowns, blocked/in-progress/ready/testing/documenting/reviewing task IDs
- **Source:** app/tools/workflow.go

---

## 6. Lifecycle Management (app/tools/lifecycle.go, lifecycle_handlers.go)

### 28. advance_task
- **Description:** Advance task to next lifecycle stage. Gated transitions (from in-progress, in-testing, in-docs, in-review) require 'evidence' describing work done.
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | epic_id | string | Yes | Epic ID |
  | story_id | string | Yes | Story ID |
  | task_id | string | Yes | Task ID |
  | evidence | string | No | Required for gated transitions. Describe tests run, docs written, or review findings. |
  | files | array of strings | No | File paths as evidence (tests, docs, code). Must exist and be modified after task started. |

- **Returns:** Task data with from/to status, evidence, gate hints for next stage
- **Source:** app/tools/lifecycle.go

### 29. reject_task
- **Description:** Reject task from review, auto-creates bug
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | epic_id | string | Yes | Epic ID |
  | story_id | string | Yes | Story ID |
  | task_id | string | Yes | Task ID |
  | reason | string | No | Rejection reason |

- **Returns:** Rejected task and auto-created bug issue data
- **Source:** app/tools/lifecycle.go

---

## 7. Sprint Management (app/tools/sprint.go)

### 30. create_sprint
- **Description:** Create a new sprint
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | name | string | Yes | Sprint name |
  | goal | string | No | Sprint goal |
  | start_date | string | No | Start date |
  | end_date | string | No | End date |

- **Returns:** Created sprint data
- **Source:** app/tools/sprint.go

### 31. list_sprints
- **Description:** List all sprints in a project
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |

- **Returns:** Array of sprint entries
- **Source:** app/tools/sprint.go

### 32. get_sprint
- **Description:** Get sprint details with task data
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | sprint_id | string | Yes | Sprint ID |

- **Returns:** Sprint data with enriched tasks, total/completed points
- **Source:** app/tools/sprint.go

### 33. start_sprint
- **Description:** Start a sprint (sets status to active)
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | sprint_id | string | Yes | Sprint ID |

- **Returns:** Sprint data with promoted-to-todo task list
- **Source:** app/tools/sprint.go

### 34. end_sprint
- **Description:** End a sprint and calculate velocity
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | sprint_id | string | Yes | Sprint ID |

- **Returns:** Sprint data, velocity, done/incomplete task lists
- **Source:** app/tools/sprint.go

### 35. add_sprint_tasks
- **Description:** Add tasks to a sprint's task list
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | sprint_id | string | Yes | Sprint ID |
  | task_ids | array of strings | Yes | Task IDs to add to the sprint |

- **Returns:** Sprint ID with updated task IDs and added count
- **Source:** app/tools/sprint.go

### 36. remove_sprint_tasks
- **Description:** Remove tasks from a sprint's task list
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | sprint_id | string | Yes | Sprint ID |
  | task_ids | array of strings | Yes | Task IDs to remove from the sprint |

- **Returns:** Sprint ID with updated task IDs
- **Source:** app/tools/sprint.go

---

## 8. Scrum Ceremonies (app/tools/scrum_*.go)

### 37. get_standup_summary
- **Description:** Get daily standup summary: completed, in-progress, blocked tasks
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | since | string | No | ISO timestamp cutoff (default 24h ago) |

- **Returns:** Completed, in-progress, and blocked task lists with assignee info
- **Source:** app/tools/scrum_standup.go

### 38. get_burndown
- **Description:** Get sprint burndown chart data
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | sprint_id | string | Yes | Sprint ID |

- **Returns:** Total/completed/remaining points, daily progress entries, projected completion date
- **Source:** app/tools/scrum_burndown.go

### 39. create_retrospective
- **Description:** Create a sprint retrospective
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | sprint_id | string | Yes | Sprint ID |
  | went_well | array of strings | No | Things that went well |
  | didnt_go_well | array of strings | No | Things that did not go well |
  | action_items | array of objects | No | Action items (action, owner, due_date) |

- **Returns:** Retrospective data
- **Source:** app/tools/scrum_retro.go

### 40. reorder_backlog
- **Description:** Reorder backlog tasks by priority rank
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | task_ids | array of strings | Yes | Task IDs in desired rank order |

- **Returns:** Array of ranked tasks with ID, title, and rank
- **Source:** app/tools/scrum_backlog.go

---

## 9. WIP Limits (app/tools/scrum_wip.go)

### 41. set_wip_limits
- **Description:** Set WIP (work-in-progress) limits for the project
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | max_in_progress | number | No | Max tasks in-progress globally |
  | max_per_assignee | number | No | Max in-progress per assignee |
  | max_per_sprint | number | No | Max tasks per sprint |

- **Returns:** Updated WIP limits
- **Source:** app/tools/scrum_wip.go

### 42. get_wip_limits
- **Description:** Get current WIP limits
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |

- **Returns:** Current WIP limits configuration
- **Source:** app/tools/scrum_wip.go

### 43. check_wip_limit
- **Description:** Check if WIP limits would be exceeded
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | assignee | string | No | Optional assignee to check |

- **Returns:** Within-limits boolean, in-progress count, limits, violations (if any)
- **Source:** app/tools/scrum_wip.go

---

## 10. Dependency Management (app/tools/dependency.go)

### 44. add_dependency
- **Description:** Add a dependency (task depends on another task)
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | epic_id | string | Yes | Epic ID |
  | story_id | string | Yes | Story ID |
  | task_id | string | Yes | Task ID (the dependent) |
  | depends_on_id | string | Yes | Task ID it depends on |

- **Returns:** Updated task data (with cycle detection)
- **Source:** app/tools/dependency.go

### 45. remove_dependency
- **Description:** Remove a dependency from a task
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | epic_id | string | Yes | Epic ID |
  | story_id | string | Yes | Story ID |
  | task_id | string | Yes | Task ID |
  | depends_on_id | string | Yes | Dependency to remove |

- **Returns:** Updated task data
- **Source:** app/tools/dependency.go

### 46. get_dependency_graph
- **Description:** Get dependency graph for a project
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |

- **Returns:** Nodes (tasks with dependencies), edges (from/to), blocked task IDs
- **Source:** app/tools/dependency.go

---

## 11. Metadata & Templates (app/tools/metadata.go)

### 47. assign_task
- **Description:** Assign a task to a person
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | epic_id | string | Yes | Epic ID |
  | story_id | string | Yes | Story ID |
  | task_id | string | Yes | Task ID |
  | assignee | string | Yes | Person to assign |

- **Returns:** Updated task data
- **Source:** app/tools/metadata.go

### 48. unassign_task
- **Description:** Remove assignment from a task
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | epic_id | string | Yes | Epic ID |
  | story_id | string | Yes | Story ID |
  | task_id | string | Yes | Task ID |

- **Returns:** Updated task data
- **Source:** app/tools/metadata.go

### 49. my_tasks
- **Description:** List all tasks assigned to a person
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | assignee | string | Yes | Person name |

- **Returns:** Array of matching task data
- **Source:** app/tools/metadata.go

### 50. add_labels
- **Description:** Add labels to a task (deduped)
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | epic_id | string | Yes | Epic ID |
  | story_id | string | Yes | Story ID |
  | task_id | string | Yes | Task ID |
  | labels | array of strings | Yes | Labels to add |

- **Returns:** Updated task data
- **Source:** app/tools/metadata.go

### 51. remove_labels
- **Description:** Remove labels from a task
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | epic_id | string | Yes | Epic ID |
  | story_id | string | Yes | Story ID |
  | task_id | string | Yes | Task ID |
  | labels | array of strings | Yes | Labels to remove |

- **Returns:** Updated task data
- **Source:** app/tools/metadata.go

### 52. set_estimate
- **Description:** Set story point estimate on a task
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | epic_id | string | Yes | Epic ID |
  | story_id | string | Yes | Story ID |
  | task_id | string | Yes | Task ID |
  | points | number | Yes | Story point estimate |

- **Returns:** Updated task data
- **Source:** app/tools/metadata.go

### 53. get_velocity
- **Description:** Calculate average velocity over completed sprints
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |

- **Returns:** Average velocity, completed sprint count, total points
- **Source:** app/tools/metadata.go

### 54. add_link
- **Description:** Add a link to a task
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | epic_id | string | Yes | Epic ID |
  | story_id | string | Yes | Story ID |
  | task_id | string | Yes | Task ID |
  | type | string | Yes | Link type: pr, commit, url, issue |
  | url | string | Yes | URL |
  | title | string | No | Link title |

- **Returns:** Updated task data
- **Source:** app/tools/metadata.go

### 55. save_template
- **Description:** Save a reusable issue template
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | name | string | Yes | Template name |
  | type | string | Yes | Type: epic, story, task, bug |
  | title | string | No | Template title |
  | priority | string | No | Default priority |
  | labels | array of strings | No | Default labels |
  | description | string | No | Template description |

- **Returns:** Saved template data
- **Source:** app/tools/metadata.go

### 56. list_templates
- **Description:** List all saved templates
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |

- **Returns:** Array of template data
- **Source:** app/tools/metadata.go

### 57. create_from_template
- **Description:** Create a task from a saved template
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | epic_id | string | Yes | Epic ID |
  | story_id | string | Yes | Story ID |
  | template_name | string | Yes | Template name to use |

- **Returns:** Created task issue data
- **Source:** app/tools/metadata.go

---

## 12. PRD Sessions (app/tools/prd.go, prd_*.go)

### 58. start_prd_session
- **Description:** Start guided PRD creation
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | prd_type | string | No | PRD audience type: business, product, technical, qa (default: general) |

- **Returns:** First PRD question to answer
- **Source:** app/tools/prd.go

### 59. answer_prd_question
- **Description:** Answer current PRD question
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | answer | string | Yes | Answer text |

- **Returns:** Next question or completion result (with conditional follow-up questions)
- **Source:** app/tools/prd.go

### 60. get_prd_session
- **Description:** Get PRD session state
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |

- **Returns:** Full PRD session state
- **Source:** app/tools/prd.go

### 61. abandon_prd_session
- **Description:** Abandon PRD session
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |

- **Returns:** Confirmation text
- **Source:** app/tools/prd.go

### 62. skip_prd_question
- **Description:** Skip optional PRD question
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |

- **Returns:** Next question or error if question is required
- **Source:** app/tools/prd.go

### 63. back_prd_question
- **Description:** Go back to previous PRD question
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |

- **Returns:** Previous question
- **Source:** app/tools/prd.go

### 64. preview_prd
- **Description:** Preview PRD markdown
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |

- **Returns:** Generated PRD markdown text
- **Source:** app/tools/prd.go

### 65. split_prd
- **Description:** Split completed PRD into numbered phases
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | phases | array of strings | Yes | Phase names in order (at least 2) |

- **Returns:** Phase slugs and count
- **Source:** app/tools/prd_phases.go

### 66. list_prd_phases
- **Description:** List PRD phases for a project
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |

- **Returns:** Array of phase info (slug, name, phase number, status)
- **Source:** app/tools/prd_phases.go

### 67. validate_prd
- **Description:** Validate PRD completeness before splitting into phases
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |

- **Returns:** Completeness score, section scores, gaps, suggestions
- **Source:** app/tools/prd_validate.go

### 68. generate_backlog
- **Description:** Generate epic/story/task backlog from completed PRD
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |

- **Returns:** Instructions and PRD content for agent to create backlog using tools
- **Source:** app/tools/prd_backlog.go

### 69. get_agent_briefing
- **Description:** Get role-specific context summary from PRD
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | role | string | Yes | Role: backend, frontend, qa, devops, design, pm |

- **Returns:** Agent briefing with role, summary, key points, constraints
- **Source:** app/tools/prd_briefing.go

---

## 13. PRD Templates (app/tools/prd_templates.go)

### 70. save_prd_template
- **Description:** Save completed PRD as reusable template
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | name | string | Yes | Template name |

- **Returns:** Template name, path, answer count
- **Source:** app/tools/prd_templates.go

### 71. list_prd_templates
- **Description:** List saved PRD templates
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |

- **Returns:** Array of template summaries (name, type, answers count, created)
- **Source:** app/tools/prd_templates.go

### 72. load_prd_template
- **Description:** Load a PRD template to pre-fill a new session
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | name | string | Yes | Template name to load |

- **Returns:** Loaded template info with pre-filled answer count and current status
- **Source:** app/tools/prd_templates.go

---

## 14. Notes (app/tools/notes.go)

### 73. save_note
- **Description:** Save a new note
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | title | string | Yes | Note title |
  | content | string | Yes | Note content |
  | tags | array of strings | No | Tags |
  | pinned | boolean | No | Pin the note |
  | source_session_id | string | No | Source session ID |
  | source_message_id | string | No | Source message ID |

- **Returns:** Created note data
- **Source:** app/tools/notes.go

### 74. list_notes
- **Description:** List all notes for a project
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |

- **Returns:** Array of notes
- **Source:** app/tools/notes.go

### 75. search_notes
- **Description:** Search notes by title, content, or tags
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | query | string | Yes | Search text |

- **Returns:** Array of matching notes
- **Source:** app/tools/notes.go

### 76. update_note
- **Description:** Update an existing note
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | id | string | Yes | Note ID |
  | title | string | No | New title |
  | content | string | No | New content |
  | tags | array of strings | No | New tags |
  | pinned | boolean | No | Pin state |
  | startup_prompt | boolean | No | Show as startup prompt |
  | quick_action | boolean | No | Show as quick action |

- **Returns:** Updated note data
- **Source:** app/tools/notes.go

### 77. delete_note
- **Description:** Delete a note
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | id | string | Yes | Note ID |

- **Returns:** Deletion confirmation text
- **Source:** app/tools/notes.go

---

## 15. Bug Reporting (app/tools/bugfix.go)

### 78. report_bug
- **Description:** Report a bug under a story
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | story_id | string | Yes | Story ID to file bug under |
  | title | string | Yes | Bug title |
  | severity | string | Yes | Severity: critical, high, medium, low |
  | steps | string | No | Steps to reproduce |
  | expected | string | No | Expected behavior |
  | actual | string | No | Actual behavior |

- **Returns:** Bug ID and created status
- **Source:** app/tools/bugfix.go

### 79. log_request
- **Description:** Log a feature request or suggestion
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | type | string | Yes | Type: feature, bug, improvement, question |
  | description | string | Yes | Request description |

- **Returns:** Logged status and total count
- **Source:** app/tools/bugfix.go

---

## 16. Artifacts & Plans (app/tools/artifacts.go)

### 80. save_plan
- **Description:** Save a plan document as markdown
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | title | string | Yes | Plan title |
  | content | string | Yes | Markdown content |
  | issue_id | string | No | Related issue ID |

- **Returns:** File name and path
- **Source:** app/tools/artifacts.go

### 81. list_plans
- **Description:** List all plan documents for a project
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |

- **Returns:** Array of plan info (file, title, issue_id, created)
- **Source:** app/tools/artifacts.go

---

## 17. README Generation (app/tools/readme.go)

### 82. regenerate_readme
- **Description:** Regenerate project README from issues
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |

- **Returns:** Confirmation text
- **Source:** app/tools/readme.go

---

## 18. Memory & Sessions (app/tools/memory.go)

### 83. save_memory
- **Description:** Save a context chunk to project memory
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | content | string | Yes | Content to remember |
  | summary | string | Yes | Short summary |
  | source | string | No | Source type: task, prd, session, user |
  | source_id | string | No | Source ID (task ID, session ID) |
  | tags | array of strings | No | Tags |

- **Returns:** Confirmation (uses engine gRPC or markdown fallback)
- **Source:** app/tools/memory.go

### 84. search_memory
- **Description:** Search project memory by keyword
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | query | string | Yes | Search query |
  | limit | number | No | Max results (default 10) |

- **Returns:** Array of search results
- **Source:** app/tools/memory.go

### 85. get_context
- **Description:** Get relevant context for current work
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | query | string | Yes | What context do you need? |
  | limit | number | No | Max results (default 5) |

- **Returns:** Array of context chunks
- **Source:** app/tools/memory.go

### 86. save_session
- **Description:** Save a session summary to project memory
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | session_id | string | Yes | Session ID |
  | summary | string | Yes | Session summary |
  | events | array of objects | No | Session events |

- **Returns:** Confirmation
- **Source:** app/tools/memory.go

### 87. list_sessions
- **Description:** List recent sessions for a project
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | limit | number | No | Max results (default 20) |

- **Returns:** Array of session data
- **Source:** app/tools/memory.go

### 88. get_session
- **Description:** Get full session details
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | session_id | string | Yes | Session ID |

- **Returns:** Full session data
- **Source:** app/tools/memory.go

---

## 19. Claude Code Awareness (app/tools/claude.go, claude_install.go, claude_hooks.go)

### 89. list_skills
- **Description:** List available skills in the project
- **Parameters:** None
- **Returns:** Array of skill info (name, description)
- **Source:** app/tools/claude.go

### 90. list_agents
- **Description:** List available agents in the project
- **Parameters:** None
- **Returns:** Array of agent info (name, description)
- **Source:** app/tools/claude.go

### 91. receive_hook_event
- **Description:** Receive a Claude Code hook event
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | event_type | string | Yes | Hook event type |
  | session_id | string | No | Session ID |
  | tool_name | string | No | Tool name |
  | agent_type | string | No | Agent type |
  | data | object | No | Event data payload |

- **Returns:** Stored confirmation with event type
- **Source:** app/tools/claude_hooks.go

### 92. get_hook_events
- **Description:** Get recent Claude Code hook events
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | event_type | string | No | Filter by event type |
  | limit | number | No | Max events to return |

- **Returns:** Array of hook event results
- **Source:** app/tools/claude_hooks.go

### 93. search_memory (Claude variant)
- **Description:** Search past session memory for relevant context. Use this when you need to recall how something was implemented before, or find relevant past work.
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | query | string | Yes | What to search for |
  | limit | number | No | Maximum number of results (default: 10) |

- **Returns:** Query results with count
- **Source:** app/tools/claude_hooks.go

### 94. install_skills
- **Description:** Install bundled skills to project
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | names | array of strings | No | Skill names to install (empty = all) |

- **Returns:** Installed count and available skills list
- **Source:** app/tools/claude_install.go

### 95. install_agents
- **Description:** Install bundled agents to project
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | names | array of strings | No | Agent names to install (empty = all) |

- **Returns:** Installed count and available agents list
- **Source:** app/tools/claude_install.go

### 96. install_docs
- **Description:** Install CLAUDE.md, AGENTS.md, CONTEXT.md to project root
- **Parameters:** None
- **Returns:** Installed count and file list
- **Source:** app/tools/claude_install.go

---

## 20. Desktop Integration (app/tools/desktop.go)

### 97. open_desktop_window
- **Description:** Open a window in the Orchestra desktop app showing a specific view. Desktop app must be running.
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | name | string | Yes | Unique window identifier (e.g. 'project-status-myapp') |
  | title | string | Yes | Window title bar text |
  | route | string | Yes | SPA route to load (e.g. '/panels/project-status') |
  | width | integer | No | Window width in pixels (default: 1200) |
  | height | integer | No | Window height in pixels (default: 800) |
  | data | object | No | Data passed to the window (e.g. {"project": "myapp"}) |

- **Returns:** Confirmation text with window name
- **Source:** app/tools/desktop.go

---

## 21. Notifications & Sound (app/tools/notification.go)

### 98. send_notification
- **Description:** Send a desktop notification to the user via the Orchestra desktop app. Optionally plays a sound. The desktop app must be running. Use for agent-finished, task updates, alerts, and other user-facing events.
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | title | string | Yes | Notification title |
  | body | string | Yes | Notification body text |
  | subtitle | string | No | Optional subtitle shown below the title (macOS only) |
  | sound | string | No | Sound to play. Available: agent-finished, agent-question, agent-permission, sub-agent-start, sub-agent-finish, started, finished, updated, action, timer-start, timer-pause, timer-stop, break-start, break-ended, jira-update, linear-update, github-pull, github-push, meeting, meeting-start |

- **Returns:** Confirmation text
- **Source:** app/tools/notification.go

### 99. play_sound
- **Description:** Play a notification sound through the Orchestra desktop app without showing a notification popup. Use to provide audio feedback for events like agent completion, timers, or alerts.
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | name | string | Yes | Sound name (same options as send_notification sound parameter) |

- **Returns:** Confirmation text
- **Source:** app/tools/notification.go

---

## 22. Team Management (app/tools/team.go)

### 100. list_teams
- **Description:** List all teams the authenticated user belongs to
- **Parameters:** None
- **Returns:** Team list from settings server
- **Source:** app/tools/team.go

### 101. create_team
- **Description:** Create a new team
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | name | string | Yes | Team name |
  | description | string | No | Team description |

- **Returns:** Created team data
- **Source:** app/tools/team.go

### 102. get_team
- **Description:** Get team details including members
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | team_id | string | Yes | Team ID |

- **Returns:** Team details with members
- **Source:** app/tools/team.go

### 103. invite_member
- **Description:** Invite a member to a team by email
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | team_id | string | Yes | Team ID |
  | email | string | Yes | Email address of the person to invite |
  | role | string | No | Role for the invitee (default: "member") |

- **Returns:** Invitation data
- **Source:** app/tools/team.go

### 104. get_pending_invitations
- **Description:** Get all pending team invitations for the authenticated user
- **Parameters:** None
- **Returns:** Array of pending invitations
- **Source:** app/tools/team.go

### 105. share_with_team
- **Description:** Share an entity (project, epic, etc.) with a team
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | team_id | string | Yes | Team ID |
  | entity_type | string | Yes | Type of entity to share, e.g. "project", "epic" |
  | entity_id | string | Yes | ID of the entity to share |
  | message | string | No | Optional message to include with the share |

- **Returns:** Share confirmation data
- **Source:** app/tools/team.go

---

## 23. Usage Tracking (app/tools/usage.go)

### 106. get_usage
- **Description:** Get usage totals and recent sessions
- **Parameters:** None
- **Returns:** Totals and last 10 sessions
- **Source:** app/tools/usage.go

### 107. record_usage
- **Description:** Record token usage for current session
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | provider | string | No | Provider name |
  | model | string | No | Model name |
  | input_tokens | number | Yes | Input token count |
  | output_tokens | number | Yes | Output token count |
  | cost | number | No | Cost |

- **Returns:** Session input/output totals
- **Source:** app/tools/usage.go

### 108. reset_session_usage
- **Description:** End the current usage session
- **Parameters:** None
- **Returns:** Confirmation text
- **Source:** app/tools/usage.go

---

## 24. Session Metrics (app/tools/session_metrics.go)

### 109. start_session_metrics
- **Description:** Start tracking metrics for a new session
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | session_id | string | Yes | Unique session identifier |

- **Returns:** Confirmation text
- **Source:** app/tools/session_metrics.go

### 110. record_session_message
- **Description:** Record a message in an active session
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | session_id | string | Yes | Session identifier |
  | is_user | boolean | No | Whether the message is from a user (default true) |

- **Returns:** Confirmation text
- **Source:** app/tools/session_metrics.go

### 111. get_session_metrics
- **Description:** Get active sessions, recent completed sessions, and daily aggregates
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | days | number | No | Number of days to aggregate (default 7) |

- **Returns:** Active sessions, recent sessions, daily aggregates
- **Source:** app/tools/session_metrics.go

---

## 25. DevTools (app/tools/devtools.go)

### 112. create_dev_session
- **Description:** Create a new dev tools session (terminal, database, ssh, logs, testing, etc.)
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | type | string | Yes | Session type: terminal, database, ssh, cloud, logs, testing, services, debugger, file-explorer |
  | name | string | Yes | Display name for the session |
  | config | object | No | Connection config (host, port, user, etc.) |

- **Returns:** Session object with ID, type, name, status, created_at
- **Source:** app/tools/devtools.go

### 113. run_query
- **Description:** Execute a SQL query against a database session's connection
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | session_id | string | Yes | Database session ID |
  | query | string | Yes | SQL query to execute |
  | params | array | No | Optional query parameters |

- **Returns:** Query results from desktop app
- **Source:** app/tools/devtools.go

### 114. ssh_exec
- **Description:** Execute a command on a remote server via an SSH session
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | session_id | string | Yes | SSH session ID |
  | command | string | Yes | Command to execute on the remote server |

- **Returns:** Command output from desktop app
- **Source:** app/tools/devtools.go

### 115. terminal_exec
- **Description:** Execute a command in a terminal session and return the output
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | session_id | string | Yes | Terminal session ID |
  | command | string | Yes | Command to execute |

- **Returns:** Command output from desktop app
- **Source:** app/tools/devtools.go

### 116. list_databases
- **Description:** List all databases on a connected database server session
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | session_id | string | Yes | Database session ID |

- **Returns:** Database list from desktop app
- **Source:** app/tools/devtools.go

### 117. list_services
- **Description:** List all system services and their running status
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | session_id | string | Yes | Services session ID |

- **Returns:** Services list from desktop app
- **Source:** app/tools/devtools.go

### 118. control_service
- **Description:** Start or stop a system service
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | session_id | string | Yes | Services session ID |
  | service_name | string | Yes | Service name (e.g. postgresql, redis) |
  | action | string | Yes | Action: start or stop |

- **Returns:** Service control result from desktop app
- **Source:** app/tools/devtools.go

### 119. read_file_session
- **Description:** Read a file from a file explorer session
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | session_id | string | Yes | File explorer session ID |
  | path | string | Yes | File path relative to session root |

- **Returns:** File content from desktop app
- **Source:** app/tools/devtools.go

### 120. detect_logs
- **Description:** Auto-detect log files in the workspace for a logs session
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | session_id | string | Yes | Logs session ID |

- **Returns:** Detected log files from desktop app
- **Source:** app/tools/devtools.go

### 121. search_logs
- **Description:** Search log entries in an active logs session
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | session_id | string | Yes | Logs session ID |
  | pattern | string | No | Text pattern to search for |
  | level | string | No | Log level filter (error, warn, info, debug) |
  | limit | number | No | Max results (default 100) |

- **Returns:** Matching log entries from desktop app
- **Source:** app/tools/devtools.go

### 122. get_stacks
- **Description:** List available Docker stack templates and currently running stacks
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | session_id | string | Yes | Services session ID |

- **Returns:** Stack list from desktop app
- **Source:** app/tools/devtools.go

### 123. install_stack
- **Description:** Install and start a Docker stack template (e.g. web-dev, laravel, node-dev, data-stack)
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | session_id | string | Yes | Services session ID |
  | stack_id | string | Yes | Stack template ID (web-dev, laravel, node-dev, data-stack) |

- **Returns:** Stack installation result from desktop app
- **Source:** app/tools/devtools.go

### 124. stack_control
- **Description:** Control a running Docker stack: stop it or get its status/logs
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | session_id | string | Yes | Services session ID |
  | stack_name | string | Yes | Stack name |
  | action | string | Yes | Action: stop, status, logs |

- **Returns:** Stack control result from desktop app
- **Source:** app/tools/devtools.go

### 125. debug_launch
- **Description:** Launch a debug session with a specific adapter (delve, debugpy, node-debug, codelldb, xdebug)
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | session_id | string | Yes | Debugger session ID |
  | adapter | string | Yes | Debug adapter: delve, debugpy, node-debug, codelldb, xdebug |
  | config | object | No | Adapter-specific launch config (program, mode, cwd, etc.) |

- **Returns:** Debug session result from desktop app
- **Source:** app/tools/devtools.go

### 126. debug_control
- **Description:** Control a running debug session: continue, step_over, step_into, step_out, disconnect, stack_trace
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | session_id | string | Yes | Debugger session ID |
  | action | string | Yes | Action: continue, step_over, step_into, step_out, disconnect, stack_trace |

- **Returns:** Debug control result from desktop app
- **Source:** app/tools/devtools.go

### 127. set_breakpoint
- **Description:** Set or remove a breakpoint in a debug session
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | session_id | string | Yes | Debugger session ID |
  | file | string | Yes | Source file path |
  | line | number | Yes | Line number |
  | remove | boolean | No | Set true to remove the breakpoint |

- **Returns:** Breakpoint result from desktop app
- **Source:** app/tools/devtools.go

---

## 26. DevTools Operations (app/tools/devtools_ops.go)

### 128. run_tests
- **Description:** Run automated tests using a specified framework
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | framework | string | Yes | Test framework: go, rust, node, playwright, phpunit, pytest |
  | pattern | string | No | Test file or pattern to run |
  | config | object | No | Extra config (coverage, verbose, etc.) |

- **Returns:** Stub response (requires active testing session)
- **Source:** app/tools/devtools_ops.go

### 129. manage_service
- **Description:** Start, stop, restart, or check status of an OS service
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | name | string | Yes | Service name (postgresql, redis, nginx, etc.) |
  | action | string | Yes | Action: start, stop, restart, status, install |

- **Returns:** Stub response (requires active services session)
- **Source:** app/tools/devtools_ops.go

### 130. view_logs
- **Description:** View and search application logs from various sources
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | source | string | Yes | Log source: file, docker, cloud |
  | level | string | No | Filter level: debug, info, warn, error |
  | query | string | No | Search pattern |
  | lines | integer | No | Number of lines (default 100) |

- **Returns:** Stub response (requires active logs session)
- **Source:** app/tools/devtools_ops.go

### 131. cloud_deploy
- **Description:** Deploy to a cloud provider
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | provider | string | Yes | Cloud provider: aws, gcp, azure, laravel-cloud |
  | config | object | Yes | Deployment configuration |

- **Returns:** Stub response (requires active cloud session)
- **Source:** app/tools/devtools_ops.go

### 132. debug_attach
- **Description:** Attach a debugger to a running process or launch a debug target
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | target | string | Yes | Process name, PID, or file path |
  | language | string | Yes | Language: go, python, node, php, rust |
  | mode | string | No | Mode: launch or attach (default attach) |

- **Returns:** Stub response (requires active debugger session)
- **Source:** app/tools/devtools_ops.go

---

## 27. Figma Integration (app/tools/figma.go)

### 133. figma_get_file_meta
- **Description:** Get metadata for a Figma file (name, last modified, thumbnail URL, version). Use this to verify a file key is correct before fetching the full document tree. The file key is the alphanumeric ID in the Figma URL: figma.com/file/{KEY}/...
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | file_key | string | Yes | Figma file key from the URL (e.g. 'aBcDeFgHiJkL') |

- **Returns:** File metadata (name, last_modified, thumbnail_url, version, role)
- **Source:** app/tools/figma.go

### 134. figma_get_file
- **Description:** Get the full document tree of a Figma file by its key. Returns nodes, styles, and metadata. Use figma_list_files to find file keys.
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | file_key | string | Yes | Figma file key (from the file URL or figma_list_files) |

- **Returns:** Full Figma file document tree
- **Source:** app/tools/figma.go

### 135. figma_get_nodes
- **Description:** Get specific nodes from a Figma file by node ID. Useful for inspecting a frame, component, or layer without fetching the whole file.
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | file_key | string | Yes | Figma file key |
  | node_id | string | Yes | Node ID to fetch (e.g. '1:2' or '1-2') |

- **Returns:** Node data from Figma
- **Source:** app/tools/figma.go

### 136. figma_get_components
- **Description:** List all published components in a Figma file. Returns component names, IDs, descriptions, and thumbnails.
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | file_key | string | Yes | Figma file key |

- **Returns:** Array of component data
- **Source:** app/tools/figma.go

---

## 28. Component Library (app/tools/component.go, component_export.go)

### 137. save_component
- **Description:** Save a UI component to the component library. Provide name, framework, code (html/css/js/jsx), tags, and description.
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | name | string | Yes | Component name |
  | framework | string | Yes | Target framework: html, react, vue, svelte, angular |
  | html | string | No | HTML markup |
  | css | string | No | CSS styles |
  | js | string | No | Plain JavaScript |
  | jsx | string | No | JSX / component source |
  | tags | array of strings | No | Tags for categorization |
  | description | string | No | Component description |
  | user_id | string | No | User ID to own the component (optional, defaults to mcp-agent) |

- **Returns:** Created component with ID, name, framework, version, tags, created_at
- **Source:** app/tools/component.go

### 138. list_components
- **Description:** List and search components in the component library. Filter by framework, name, or tags.
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | framework | string | No | Filter by framework |
  | name | string | No | Fuzzy search by name |
  | tags | string | No | Comma-separated tag filter |
  | limit | integer | No | Max results (default 20) |
  | offset | integer | No | Pagination offset (default 0) |
  | user_id | string | No | Filter to this user's components (optional; omit to search all public components) |

- **Returns:** Components array with count
- **Source:** app/tools/component.go

### 139. export_component
- **Description:** Export a saved component from the library. Supports npm (package.json + index.js + styles.css), cdn (single self-contained HTML file), raw (separate html/css/js files as JSON), and zip (base64-encoded zip archive).
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | component_id | string | Yes | Component UUID to export |
  | format | string | Yes | Export format: npm, cdn, raw, zip |
  | package_name | string | No | npm package name (used when format=npm, defaults to kebab-case component name) |
  | version | string | No | Package version (used when format=npm, defaults to 1.0.0) |
  | user_id | string | No | User ID for ownership verification (optional, defaults to mcp-agent) |

- **Returns:** Exported component in the requested format
- **Source:** app/tools/component_export.go

---

## 29. Analytics (app/tools/analytics.go)

### 140. track_event
- **Description:** Track a feature analytics event (mode switch, context menu, mention, command palette, agent selection, screenshot, file attachment)
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | event_name | string | Yes | Event name (e.g. mode_switch, context_menu_open, mention_used) |
  | properties | object | No | Optional event properties |
  | distinct_id | string | No | Optional user identifier |

- **Returns:** Confirmation text
- **Source:** app/tools/analytics.go

### 141. get_analytics_events
- **Description:** Returns buffered analytics events for debugging
- **Parameters:** None
- **Returns:** Buffered event count and events array
- **Source:** app/tools/analytics.go

---

## 30. Preview (app/tools/preview.go)

### 142. preview_component
- **Description:** Create a live preview session for a UI component. Returns a session_id and WebSocket URL that streams hot-updates.
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | framework | string | Yes | Rendering framework: html, react, vue, svelte, angular, react-native, flutter |
  | html | string | No | HTML markup (for html framework or wrapper shell) |
  | css | string | No | CSS styles |
  | js | string | No | Plain JavaScript |
  | jsx | string | No | JSX / component source |

- **Returns:** Session ID, WebSocket URL, framework
- **Source:** app/tools/preview.go

### 143. update_preview
- **Description:** Hot-update the code in an active preview session. Only supply fields that changed.
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | session_id | string | Yes | Preview session ID returned by preview_component |
  | html | string | No | Updated HTML markup |
  | css | string | No | Updated CSS styles |
  | js | string | No | Updated JavaScript |
  | jsx | string | No | Updated JSX / component source |

- **Returns:** Session ID and updated confirmation
- **Source:** app/tools/preview.go

### 144. set_preview_viewport
- **Description:** Change the viewport size of a preview session. Use a preset or supply custom width/height.
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | session_id | string | Yes | Preview session ID |
  | preset | string | No | Viewport preset: mobile (375px), tablet (768px), desktop (1280px), custom |
  | width | integer | No | Custom width in pixels (required when preset=custom) |
  | height | integer | No | Custom height in pixels (required when preset=custom) |

- **Returns:** Session ID and resolved viewport
- **Source:** app/tools/preview.go

### 145. open_browser_preview
- **Description:** Signal the Chrome Extension (or any connected client) to open the preview session in a new browser tab.
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | session_id | string | Yes | Preview session ID to open in the browser |

- **Returns:** Session ID and signaled confirmation
- **Source:** app/tools/preview.go

---

## 31. GitHub Integration (app/tools/github_*.go)

### 146. github_auth_start
- **Description:** Get the GitHub OAuth authorization URL. Open this URL in the browser to connect a GitHub account.
- **Parameters:** None
- **Returns:** Authorization URL
- **Source:** app/tools/github_auth.go

### 147. github_auth_status
- **Description:** Check GitHub authentication status. Returns connected/disconnected with user info.
- **Parameters:** None
- **Returns:** Authentication state with user info
- **Source:** app/tools/github_auth.go

### 148. github_auth_pat
- **Description:** Sign in to GitHub with a Personal Access Token
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | token | string | Yes | GitHub Personal Access Token |

- **Returns:** Connection status and user info
- **Source:** app/tools/github_auth.go

### 149. github_sign_out
- **Description:** Disconnect GitHub account and remove stored credentials
- **Parameters:** None
- **Returns:** Confirmation text
- **Source:** app/tools/github_auth.go

### 150. github_list_issues
- **Description:** List issues for a GitHub repository
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | owner | string | Yes | Repository owner |
  | repo | string | Yes | Repository name |
  | state | string | No | State: open, closed, all |
  | assignee | string | No | Filter by assignee |
  | labels | string | No | Comma-separated labels |

- **Returns:** Array of issues
- **Source:** app/tools/github_issues.go

### 151. github_get_issue
- **Description:** Get issue details
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | owner | string | Yes | Repository owner |
  | repo | string | Yes | Repository name |
  | number | number | Yes | Issue number |

- **Returns:** Issue details
- **Source:** app/tools/github_issues.go

### 152. github_create_issue
- **Description:** Create a new GitHub issue
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | owner | string | Yes | Repository owner |
  | repo | string | Yes | Repository name |
  | title | string | Yes | Issue title |
  | body | string | No | Issue body |
  | assignees | array of strings | No | Assignees |
  | labels | array of strings | No | Labels |

- **Returns:** Created issue data
- **Source:** app/tools/github_issues.go

### 153. github_close_issue
- **Description:** Close a GitHub issue
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | owner | string | Yes | Repository owner |
  | repo | string | Yes | Repository name |
  | number | number | Yes | Issue number |

- **Returns:** Confirmation text
- **Source:** app/tools/github_issues.go

### 154. github_issue_comment
- **Description:** Add a comment to an issue
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | owner | string | Yes | Repository owner |
  | repo | string | Yes | Repository name |
  | number | number | Yes | Issue number |
  | body | string | Yes | Comment text |

- **Returns:** Confirmation text
- **Source:** app/tools/github_issues.go

### 155. github_ci_status
- **Description:** Get CI/CD status for a commit ref (branch, tag, or SHA)
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | owner | string | Yes | Repository owner |
  | repo | string | Yes | Repository name |
  | ref | string | Yes | Branch name, tag, or commit SHA |

- **Returns:** CI status data
- **Source:** app/tools/github_issues.go

### 156. github_list_prs
- **Description:** List pull requests for a GitHub repository
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | owner | string | Yes | Repository owner |
  | repo | string | Yes | Repository name |
  | state | string | No | State: open, closed, all |
  | sort | string | No | Sort by: created, updated, popularity |
  | direction | string | No | Direction: asc, desc |

- **Returns:** Array of pull requests
- **Source:** app/tools/github_pr.go

### 157. github_get_pr
- **Description:** Get pull request details
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | owner | string | Yes | Repository owner |
  | repo | string | Yes | Repository name |
  | number | number | Yes | PR number |

- **Returns:** Pull request details
- **Source:** app/tools/github_pr.go

### 158. github_pr_files
- **Description:** List files changed in a pull request
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | owner | string | Yes | Repository owner |
  | repo | string | Yes | Repository name |
  | number | number | Yes | PR number |

- **Returns:** Array of changed files
- **Source:** app/tools/github_pr.go

### 159. github_create_pr
- **Description:** Create a new pull request
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | owner | string | Yes | Repository owner |
  | repo | string | Yes | Repository name |
  | title | string | Yes | PR title |
  | head | string | Yes | Branch with changes |
  | base | string | Yes | Branch to merge into |
  | body | string | No | PR body |
  | draft | boolean | No | Create as draft |

- **Returns:** Created PR data
- **Source:** app/tools/github_pr.go

### 160. github_merge_pr
- **Description:** Merge a pull request
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | owner | string | Yes | Repository owner |
  | repo | string | Yes | Repository name |
  | number | number | Yes | PR number |
  | method | string | No | Merge method: merge, squash, rebase |

- **Returns:** Confirmation text
- **Source:** app/tools/github_pr.go

### 161. github_review_pr
- **Description:** Submit a review on a pull request
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | owner | string | Yes | Repository owner |
  | repo | string | Yes | Repository name |
  | number | number | Yes | PR number |
  | event | string | Yes | Review event: APPROVE, REQUEST_CHANGES, COMMENT |
  | body | string | No | Review body |

- **Returns:** Confirmation text
- **Source:** app/tools/github_pr.go

### 162. github_pr_comment
- **Description:** Add a comment to a pull request
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | owner | string | Yes | Repository owner |
  | repo | string | Yes | Repository name |
  | number | number | Yes | PR number |
  | body | string | Yes | Comment text |

- **Returns:** Confirmation text
- **Source:** app/tools/github_pr.go

### 163. github_list_repos
- **Description:** List tracked GitHub repositories for a project
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |

- **Returns:** Array of tracked repositories
- **Source:** app/tools/github_repos.go

### 164. github_track_repo
- **Description:** Add a GitHub repository to tracking for activity notifications
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | owner | string | Yes | Repository owner |
  | repo | string | Yes | Repository name |
  | default_branch | string | No | Default branch name |
  | watch_prs | boolean | No | Watch pull requests |
  | watch_issues | boolean | No | Watch issues |

- **Returns:** Tracked repository data
- **Source:** app/tools/github_repos.go

### 165. github_untrack_repo
- **Description:** Remove a GitHub repository from tracking
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Project slug |
  | repo_id | string | Yes | Tracked repo ID |

- **Returns:** Confirmation text
- **Source:** app/tools/github_repos.go

---

## 32. Jira Integration (app/tools/jira_*.go)

### 166. jira_auth_status
- **Description:** Check Jira authentication status. Returns connected/disconnected with user info.
- **Parameters:** None
- **Returns:** Authentication state
- **Source:** app/tools/jira_auth.go

### 167. jira_sign_out
- **Description:** Disconnect Jira account and remove stored credentials
- **Parameters:** None
- **Returns:** Confirmation text
- **Source:** app/tools/jira_auth.go

### 168. jira_search
- **Description:** Search Jira issues using JQL
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | jql | string | Yes | JQL query string |
  | max_results | number | No | Max results (default 50) |

- **Returns:** Array of Jira issues
- **Source:** app/tools/jira_issues.go

### 169. jira_get_issue
- **Description:** Get Jira issue details by key (e.g. PROJ-123)
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | key | string | Yes | Issue key (e.g. PROJ-123) |

- **Returns:** Jira issue details
- **Source:** app/tools/jira_issues.go

### 170. jira_create_issue
- **Description:** Create a new Jira issue
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project_key | string | Yes | Project key (e.g. PROJ) |
  | summary | string | Yes | Issue summary/title |
  | description | string | No | Issue description |
  | issue_type | string | No | Issue type (Task, Bug, Story, Epic) |
  | priority | string | No | Priority (Highest, High, Medium, Low, Lowest) |

- **Returns:** Created Jira issue
- **Source:** app/tools/jira_issues.go

### 171. jira_update_issue
- **Description:** Update a Jira issue
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | key | string | Yes | Issue key (e.g. PROJ-123) |
  | summary | string | No | New summary |
  | description | string | No | New description |

- **Returns:** Confirmation text
- **Source:** app/tools/jira_issues.go

### 172. jira_transition_issue
- **Description:** Transition a Jira issue to a new status
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | key | string | Yes | Issue key (e.g. PROJ-123) |
  | transition_id | string | Yes | Transition ID |

- **Returns:** Confirmation text
- **Source:** app/tools/jira_issues.go

### 173. jira_sync_pull
- **Description:** Pull issues from Jira into local TOON files
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Local project slug |
  | jql | string | Yes | JQL query to filter issues |
  | max_results | number | No | Max issues to sync (default 50) |

- **Returns:** Sync result data
- **Source:** app/tools/jira_sync.go

### 174. jira_sync_push
- **Description:** Push a local TOON issue to Jira
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Local project slug |
  | local_id | string | Yes | Local issue ID to push |

- **Returns:** Confirmation text
- **Source:** app/tools/jira_sync.go

### 175. jira_sync_status
- **Description:** Show sync mappings between local and Jira issues
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Local project slug |

- **Returns:** Sync configuration and mappings
- **Source:** app/tools/jira_sync.go

---

## 33. Linear Integration (app/tools/linear_*.go)

### 176. linear_auth_status
- **Description:** Check Linear authentication status. Returns connected/disconnected with user info.
- **Parameters:** None
- **Returns:** Authentication state
- **Source:** app/tools/linear_auth.go

### 177. linear_sign_out
- **Description:** Disconnect Linear account and remove stored credentials
- **Parameters:** None
- **Returns:** Confirmation text
- **Source:** app/tools/linear_auth.go

### 178. linear_list_issues
- **Description:** List Linear issues for a team
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | team_id | string | Yes | Linear team ID |
  | state | string | No | Filter by state type (backlog, unstarted, started, completed, cancelled) |

- **Returns:** Array of Linear issues
- **Source:** app/tools/linear_issues.go

### 179. linear_get_issue
- **Description:** Get Linear issue details by ID
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | id | string | Yes | Linear issue ID |

- **Returns:** Linear issue details
- **Source:** app/tools/linear_issues.go

### 180. linear_create_issue
- **Description:** Create a new Linear issue
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | team_id | string | Yes | Team ID |
  | title | string | Yes | Issue title |
  | description | string | No | Issue description (markdown) |
  | priority | number | No | Priority (0=none, 1=urgent, 2=high, 3=medium, 4=low) |

- **Returns:** Created Linear issue
- **Source:** app/tools/linear_issues.go

### 181. linear_update_issue
- **Description:** Update a Linear issue
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | id | string | Yes | Issue ID |
  | title | string | No | New title |
  | description | string | No | New description |

- **Returns:** Confirmation text
- **Source:** app/tools/linear_issues.go

### 182. linear_list_teams
- **Description:** List Linear teams
- **Parameters:** None
- **Returns:** Array of Linear teams
- **Source:** app/tools/linear_issues.go

### 183. linear_sync_pull
- **Description:** Pull issues from Linear into local TOON files
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Local project slug |
  | team_id | string | Yes | Linear team ID to sync from |

- **Returns:** Sync result data
- **Source:** app/tools/linear_sync.go

### 184. linear_sync_push
- **Description:** Push a local TOON issue to Linear
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Local project slug |
  | local_id | string | Yes | Local issue ID to push |

- **Returns:** Confirmation text
- **Source:** app/tools/linear_sync.go

### 185. linear_sync_status
- **Description:** Show sync mappings between local and Linear issues
- **Parameters:**

  | Name | Type | Required | Description |
  |------|------|----------|-------------|
  | project | string | Yes | Local project slug |

- **Returns:** Sync configuration and mappings
- **Source:** app/tools/linear_sync.go

---

## 34. CLI Commands (app/tools/cli_commands.go)

### 186. list_cli_commands
- **Description:** List available CLI commands for the chat interface
- **Parameters:** None
- **Returns:** Array of command info (name, description, shortcut)
- **Source:** app/tools/cli_commands.go

---

## MCP Resources (app/tools/resources.go)

Resources are read-only data exposed via `resources/list` and `resources/read` in the MCP protocol.

### R1. project_prd
- **URI:** `orchestra://project/{slug}/prd`
- **Title:** Project PRD Document
- **Description:** The Product Requirements Document for a project
- **MIME Type:** text/markdown

### R2. project_status
- **URI:** `orchestra://project/{slug}/status`
- **Title:** Project Status
- **Description:** Current project status with epic/story/task summaries
- **MIME Type:** application/json

### R3. task_detail
- **URI:** `orchestra://project/{slug}/task/{epicId}/{storyId}/{taskId}`
- **Title:** Task Detail
- **Description:** Full detail of a specific task
- **MIME Type:** application/json

---

## MCP Prompts (app/tools/prompts.go)

Prompts are templated conversation starters exposed via `prompts/list` and `prompts/get` in the MCP protocol.

### P1. create_prd
- **Title:** Create PRD
- **Description:** Guided product requirements document creation
- **Arguments:**

  | Name | Required | Description |
  |------|----------|-------------|
  | project_name | Yes | Name of the project |
  | description | No | Brief project description |

### P2. review_task
- **Title:** Review Task
- **Description:** Generate a code review prompt for a specific task
- **Arguments:**

  | Name | Required | Description |
  |------|----------|-------------|
  | project | Yes | Project slug |
  | epic_id | Yes | Epic ID |
  | story_id | Yes | Story ID |
  | task_id | Yes | Task ID to review |

### P3. plan_sprint
- **Title:** Plan Sprint
- **Description:** Generate a sprint planning prompt with current backlog
- **Arguments:**

  | Name | Required | Description |
  |------|----------|-------------|
  | project | Yes | Project slug |

---

## Summary by Category

| # | Category | Tool Count | Source Files |
|---|----------|-----------|-------------|
| 1 | Project Management | 7 | project.go, project_tree.go |
| 2 | Epic Management | 5 | epic.go, epic_handlers.go |
| 3 | Story Management | 5 | story.go, story_handlers.go |
| 4 | Task Management | 5 | task.go, task_update.go |
| 5 | Workflow Management | 5 | workflow.go, workflow_cascade.go |
| 6 | Lifecycle Management | 2 | lifecycle.go, lifecycle_handlers.go |
| 7 | Sprint Management | 7 | sprint.go |
| 8 | Scrum Ceremonies | 4 | scrum_standup.go, scrum_burndown.go, scrum_retro.go, scrum_backlog.go |
| 9 | WIP Limits | 3 | scrum_wip.go |
| 10 | Dependency Management | 3 | dependency.go |
| 11 | Metadata & Templates | 11 | metadata.go |
| 12 | PRD Sessions | 12 | prd.go, prd_phases.go, prd_validate.go, prd_backlog.go, prd_briefing.go |
| 13 | PRD Templates | 3 | prd_templates.go |
| 14 | Notes | 5 | notes.go |
| 15 | Bug Reporting | 2 | bugfix.go |
| 16 | Artifacts & Plans | 2 | artifacts.go |
| 17 | README Generation | 1 | readme.go |
| 18 | Memory & Sessions | 6 | memory.go |
| 19 | Claude Code Awareness | 8 | claude.go, claude_install.go, claude_hooks.go |
| 20 | Desktop Integration | 1 | desktop.go |
| 21 | Notifications & Sound | 2 | notification.go |
| 22 | Team Management | 6 | team.go |
| 23 | Usage Tracking | 3 | usage.go |
| 24 | Session Metrics | 3 | session_metrics.go |
| 25 | DevTools | 16 | devtools.go |
| 26 | DevTools Operations | 5 | devtools_ops.go |
| 27 | Figma Integration | 4 | figma.go |
| 28 | Component Library | 3 | component.go, component_export.go |
| 29 | Analytics | 2 | analytics.go |
| 30 | Preview | 4 | preview.go |
| 31 | GitHub Integration | 20 | github_auth.go, github_issues.go, github_pr.go, github_repos.go |
| 32 | Jira Integration | 10 | jira_auth.go, jira_issues.go, jira_sync.go |
| 33 | Linear Integration | 10 | linear_auth.go, linear_issues.go, linear_sync.go |
| 34 | CLI Commands | 1 | cli_commands.go |
| -- | **Resources** | 3 | resources.go |
| -- | **Prompts** | 3 | prompts.go |
| | **TOTAL** | **186 tools + 3 resources + 3 prompts** | |

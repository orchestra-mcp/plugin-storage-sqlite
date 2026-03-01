---
created_at: "2026-02-28T02:28:58Z"
description: 'Version control plugin providing full Git operations + GitHub integration via gh CLI / GitHub API. 20 tools: git_status, git_diff, git_log, git_commit, git_branch, git_merge, git_stash, git_blame, git_checkout, git_tag, gh_pr_create, gh_pr_list, gh_pr_review, gh_pr_merge, gh_issue_create, gh_issue_list, gh_issue_comment, gh_actions_status, gh_release_create, gh_repo_info. Scaffold via /plugin-generator, separate GitHub repo. Uses os/exec for git + gh CLI, with optional GitHub API fallback via go-github.'
id: FEAT-LXW
labels:
    - phase-6
    - devtools
    - git
    - github
    - high-value
priority: P0
project_id: orchestra-tools
status: done
title: Git & GitHub integration plugin (devtools.git)
updated_at: "2026-02-28T05:25:00Z"
version: 0
---

# Git & GitHub integration plugin (devtools.git)

Version control plugin providing full Git operations + GitHub integration via gh CLI / GitHub API. 20 tools: git_status, git_diff, git_log, git_commit, git_branch, git_merge, git_stash, git_blame, git_checkout, git_tag, gh_pr_create, gh_pr_list, gh_pr_review, gh_pr_merge, gh_issue_create, gh_issue_list, gh_issue_comment, gh_actions_status, gh_release_create, gh_repo_info. Scaffold via /plugin-generator, separate GitHub repo. Uses os/exec for git + gh CLI, with optional GitHub API fallback via go-github.

---
**in-progress -> done**: 22 tests passing: 9 in git/exec_test.go (init, status, log, diff, branch, tag, blame, stash, invalid command) + 13 in tools/tools_test.go (git_status, git_log, git_diff, git_commit, git_branch list/create/missing-action, git_tag create/list/missing-action). Real temp git repos used. Binary built to bin/devtools-git. Wired into plugins.yaml.

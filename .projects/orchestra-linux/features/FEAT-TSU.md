---
created_at: "2026-02-28T02:52:49Z"
description: 'LocalCache class at ~/.local/share/orchestra/cache.db. Tables: projects (slug TEXT PK, data TEXT, updated_at INTEGER), notes (id TEXT PK, data TEXT, updated_at INTEGER), sessions (id TEXT PK, data TEXT, updated_at INTEGER). Methods: cache_projects(), cached_projects(), cache_notes(), cached_notes(), cache_sessions(), invalidate(entity). Used to pre-populate UI when orchestrator is not yet connected. Cache refreshed after each successful tool call response. 24-hour TTL per entity type.'
id: FEAT-TSU
priority: P1
project_id: orchestra-linux
status: backlog
title: Local SQLite cache
updated_at: "2026-02-28T02:52:49Z"
version: 0
---

# Local SQLite cache

LocalCache class at ~/.local/share/orchestra/cache.db. Tables: projects (slug TEXT PK, data TEXT, updated_at INTEGER), notes (id TEXT PK, data TEXT, updated_at INTEGER), sessions (id TEXT PK, data TEXT, updated_at INTEGER). Methods: cache_projects(), cached_projects(), cache_notes(), cached_notes(), cache_sessions(), invalidate(entity). Used to pre-populate UI when orchestrator is not yet connected. Cache refreshed after each successful tool call response. 24-hour TTL per entity type.

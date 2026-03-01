---
created_at: "2026-02-28T02:51:18Z"
description: 'High-level client wrapping QUICConnection. Handles MCP tool call serialization/deserialization via JSON-glib. Methods: send_tool_call(name, arguments) → Json.Object, subscribe_events(topic), publish_event(topic, payload). Manages connection lifecycle and reconnection. Exposes connection_state property (GObject notify signal) for UI binding.'
id: FEAT-GMK
priority: P0
project_id: orchestra-linux
status: backlog
title: OrchestraClient high-level orchestrator service
updated_at: "2026-02-28T02:51:18Z"
version: 0
---

# OrchestraClient high-level orchestrator service

High-level client wrapping QUICConnection. Handles MCP tool call serialization/deserialization via JSON-glib. Methods: send_tool_call(name, arguments) → Json.Object, subscribe_events(topic), publish_event(topic, payload). Manages connection lifecycle and reconnection. Exposes connection_state property (GObject notify signal) for UI binding.

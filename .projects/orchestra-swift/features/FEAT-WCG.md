---
created_at: "2026-02-28T02:35:26Z"
description: 'AppState: root ObservableObject holding QUICConnection, PluginRegistry, ToolService, SettingsService, cached data. OrchestraClient: wraps QUICConnection + StreamFramer, callTool(name:arguments:) async, listTools() async, connection state publisher for SwiftUI binding. ToolService: high-level MCP tool proxy with convenience methods (listProjects, getProjectStatus, createFeature, advanceFeature, listNotes, saveNote, aiPrompt with provider/model).'
id: FEAT-WCG
priority: P0
project_id: orchestra-swift
status: done
title: AppState + OrchestraClient + ToolService
updated_at: "2026-02-28T05:15:00Z"
version: 0
---

# AppState + OrchestraClient + ToolService

AppState: root ObservableObject holding QUICConnection, PluginRegistry, ToolService, SettingsService, cached data. OrchestraClient: wraps QUICConnection + StreamFramer, callTool(name:arguments:) async, listTools() async, connection state publisher for SwiftUI binding. ToolService: high-level MCP tool proxy with convenience methods (listProjects, getProjectStatus, createFeature, advanceFeature, listNotes, saveNote, aiPrompt with provider/model).

---
created_at: "2026-02-28T02:51:44Z"
description: 'ChatPlugin implementing OrchestraPlugin (id: "chat", section: SIDEBAR, icon: "chat-message-new-symbolic"). ChatView with AdwNavigationSplitView: sidebar (session list GtkListBox) + ChatBox (AdwNavigationPage). Session list: GtkSearchEntry, New Chat GtkButton (+), AdwActionRow per session (name, provider pill, model pill, timestamp). ChatBox: AdwHeaderBar (session name, provider badge, model badge, Live dot), GtkScrolledWindow with GtkListBox (messages), StatusLine, ChatInput. ChatInput: GtkTextView (auto-resize 1-6 lines) + GtkBox (provider picker, model picker, tools toggle, attach, send/stop button). Send via ToolService.spawn_session().'
id: FEAT-BVE
priority: P0
project_id: orchestra-linux
status: backlog
title: AI Chat plugin — session list + multi-LLM conversation
updated_at: "2026-02-28T02:51:44Z"
version: 0
---

# AI Chat plugin — session list + multi-LLM conversation

ChatPlugin implementing OrchestraPlugin (id: "chat", section: SIDEBAR, icon: "chat-message-new-symbolic"). ChatView with AdwNavigationSplitView: sidebar (session list GtkListBox) + ChatBox (AdwNavigationPage). Session list: GtkSearchEntry, New Chat GtkButton (+), AdwActionRow per session (name, provider pill, model pill, timestamp). ChatBox: AdwHeaderBar (session name, provider badge, model badge, Live dot), GtkScrolledWindow with GtkListBox (messages), StatusLine, ChatInput. ChatInput: GtkTextView (auto-resize 1-6 lines) + GtkBox (provider picker, model picker, tools toggle, attach, send/stop button). Send via ToolService.spawn_session().

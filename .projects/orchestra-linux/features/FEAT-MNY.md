---
created_at: "2026-02-28T02:55:10Z"
description: 'Implement ai.browser-context tool integration via WebSocket server listening on localhost:8765. GLib.SocketService + Soup.WebsocketConnection to handle Chrome extension connection. Protocol: extension sends PageContent JSON (title, url, selectedText, mainContent, metaDescription, headings, language, isArticle). Store last received page context in AppState. Tools: get_page_content(), get_page_dom(), get_selected_text(), get_open_tabs(), get_page_screenshot(), navigate_to(), execute_script() — all route via WebSocket message to extension and await response with GLib timeout. Show connection status badge in chat input tray.'
id: FEAT-MNY
priority: P2
project_id: orchestra-linux
status: backlog
title: Browser context bridge (Chrome extension WebSocket)
updated_at: "2026-02-28T02:55:10Z"
version: 0
---

# Browser context bridge (Chrome extension WebSocket)

Implement ai.browser-context tool integration via WebSocket server listening on localhost:8765. GLib.SocketService + Soup.WebsocketConnection to handle Chrome extension connection. Protocol: extension sends PageContent JSON (title, url, selectedText, mainContent, metaDescription, headings, language, isArticle). Store last received page context in AppState. Tools: get_page_content(), get_page_dom(), get_selected_text(), get_open_tabs(), get_page_screenshot(), navigate_to(), execute_script() — all route via WebSocket message to extension and await response with GLib timeout. Show connection status badge in chat input tray.

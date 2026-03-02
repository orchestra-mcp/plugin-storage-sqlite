---
created_at: "2026-03-01T12:32:05Z"
description: Transform the monolithic 36-plugin binary into a core+optional architecture. 4 core plugins (storage.markdown, transport.stdio, tools.features, tools.marketplace) stay compiled in-process. 32 optional plugins become separate pre-built binaries downloaded from GitHub, managed via `orchestra plugin install/remove/enable/disable` CLI commands. External plugins connect via QUIC and are spawned as child processes on serve startup.
features:
    - FEAT-PUO
    - FEAT-GUX
    - FEAT-HLD
    - FEAT-LNU
    - FEAT-REI
    - FEAT-PLK
id: PLAN-YPA
project_id: orchestra-tools
status: completed
title: Selective Plugin Loading System
updated_at: "2026-03-01T12:56:51Z"
version: 0
---

# Selective Plugin Loading System

Transform the monolithic 36-plugin binary into a core+optional architecture. 4 core plugins (storage.markdown, transport.stdio, tools.features, tools.marketplace) stay compiled in-process. 32 optional plugins become separate pre-built binaries downloaded from GitHub, managed via `orchestra plugin install/remove/enable/disable` CLI commands. External plugins connect via QUIC and are spawned as child processes on serve startup.

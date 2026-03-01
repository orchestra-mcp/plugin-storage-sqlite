---
created_at: "2026-02-28T02:55:39Z"
description: |-
    Implement `Orchestra.Core/Transport/` — the QUIC connection layer that connects `ui.windows` to the orchestrator.

    **Files:**
    - `QUICConnection.cs` — `System.Net.Quic` client, connect/disconnect, exponential backoff reconnect (1s→30s), `ConnectionState` events
    - `StreamFramer.cs` — length-delimited Protobuf framing: `[4-byte big-endian uint32 length][N bytes payload]`, max 16MB. Must match `libs/sdk-go/plugin/framing.go` byte-for-byte
    - `MTLSConfig.cs` — load `%USERPROFILE%\.orchestra\certs\{pluginId}.crt/key` + `ca.crt`, build `X509Certificate2`, validate server cert against CA chain

    **Key API:** `ConnectAsync(host, port, pluginId, ct)`, `WriteAsync(IMessage, ct)`, `ReadAsync<T>(ct)`, `DisposeAsync()`

    **ALPN:** `orchestra-plugin`

    **ConnectionState enum:** `Disconnected`, `Connecting`, `Connected`, `Reconnecting`

    **Platform:** `net8.0` (System.Net.Quic requires .NET 8+ on Windows 10 19041+)
id: FEAT-CUN
priority: P0
project_id: orchestra-win
status: backlog
title: QUIC transport — System.Net.Quic client
updated_at: "2026-02-28T02:55:39Z"
version: 0
---

# QUIC transport — System.Net.Quic client

Implement `Orchestra.Core/Transport/` — the QUIC connection layer that connects `ui.windows` to the orchestrator.

**Files:**
- `QUICConnection.cs` — `System.Net.Quic` client, connect/disconnect, exponential backoff reconnect (1s→30s), `ConnectionState` events
- `StreamFramer.cs` — length-delimited Protobuf framing: `[4-byte big-endian uint32 length][N bytes payload]`, max 16MB. Must match `libs/sdk-go/plugin/framing.go` byte-for-byte
- `MTLSConfig.cs` — load `%USERPROFILE%\.orchestra\certs\{pluginId}.crt/key` + `ca.crt`, build `X509Certificate2`, validate server cert against CA chain

**Key API:** `ConnectAsync(host, port, pluginId, ct)`, `WriteAsync(IMessage, ct)`, `ReadAsync<T>(ct)`, `DisposeAsync()`

**ALPN:** `orchestra-plugin`

**ConnectionState enum:** `Disconnected`, `Connecting`, `Connected`, `Reconnecting`

**Platform:** `net8.0` (System.Net.Quic requires .NET 8+ on Windows 10 19041+)

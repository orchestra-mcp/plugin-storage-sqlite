---
created_at: "2026-02-28T02:48:59Z"
description: 'Implement QUIC transport client using ngtcp2 C library with Vala VAPI bindings. Connect to orchestrator at localhost:50100. Support mTLS authentication using certificates from ~/.orchestra/certs/ (ca.crt, plugin cert/key). ALPN protocol: "orchestra-plugin". Include exponential backoff reconnection (1s to 30s max). Integrate with GLib main loop for async I/O.'
id: FEAT-ILH
priority: P0
project_id: orchestra-linux
status: backlog
title: QUIC transport client (ngtcp2)
updated_at: "2026-02-28T02:48:59Z"
version: 0
---

# QUIC transport client (ngtcp2)

Implement QUIC transport client using ngtcp2 C library with Vala VAPI bindings. Connect to orchestrator at localhost:50100. Support mTLS authentication using certificates from ~/.orchestra/certs/ (ca.crt, plugin cert/key). ALPN protocol: "orchestra-plugin". Include exponential backoff reconnection (1s to 30s max). Integrate with GLib main loop for async I/O.

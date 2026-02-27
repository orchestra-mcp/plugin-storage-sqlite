# QUIC Mesh Architecture — Orchestra Rebuild

> Event-driven microservice communication over QUIC protocol with peer-to-peer mesh networking.
> No central broker. Services discover each other, connect directly, and exchange Protobuf events.

---

## 1. Core Decisions

| Decision | Choice | Why |
|----------|--------|-----|
| Transport | QUIC (RFC 9000) | Multiplexed streams, 0-RTT reconnect, built-in TLS 1.3 |
| Browser transport | WebTransport over HTTP/3 | Native QUIC for browsers (Chrome, Firefox) |
| Message format | Protobuf | Typed, fast, versioned, generated for 6 languages |
| Auth | mTLS (mutual TLS) | Every service authenticates to every other service |
| Discovery (local) | Seed file (`~/.orchestra/mesh.json`) | Simple, explicit, works everywhere |
| Discovery (LAN) | mDNS/Bonjour (`_orch-mesh._udp.local.`) | Zero-config for multi-machine setups |
| Topology | Peer-to-peer mesh | No single point of failure, no central broker |
| Delivery | Best-effort default, at-least-once for state changes | Local outbox for offline peers |
| Ordering | Lamport sequence + CausalRef | Partial causal ordering, no vector clocks needed |

---

## 2. Technology Choices

### Go QUIC: `quic-go` (github.com/quic-go/quic-go)

Production-ready, pure-Go QUIC implementation (RFC 9000/9001/9002) with HTTP/3 support.

Key capabilities:
- **Bidirectional streams** (`quic.Stream`): Both endpoints read and write. Used for control channel.
- **Unidirectional streams** (`quic.SendStream`/`quic.ReceiveStream`): One-way. Used for event delivery.
- **Unreliable datagrams** (RFC 9221): Encrypted, congestion-controlled, NOT retransmitted. Used for heartbeats.
- **Connection multiplexing**: Unlimited concurrent streams over a single QUIC connection.
- **0-RTT**: After first handshake, subsequent connections resume instantly.

### Rust QUIC: `quinn`

Standard Rust QUIC crate (86M+ downloads). Used by `orch-engine`. The `tonic-h3` crate enables running existing Tonic gRPC services over HTTP/3/QUIC.

### mDNS: `grandcat/zeroconf`

Chosen over `hashicorp/mdns` for:
- Full RFC 6762 (mDNS) and RFC 6763 (DNS-SD) support
- Better IPv6 handling
- Channel-based API for Go concurrency
- More actively maintained

### WebTransport: `quic-go/webtransport-go`

Official companion to quic-go. Implements WebTransport over HTTP/3 for browser connectivity.

Browser-side:
```javascript
const wt = new WebTransport("https://localhost:4443/mesh");
await wt.ready;
const stream = await wt.createBidirectionalStream();
```

### mTLS over QUIC

QUIC mandates TLS 1.3. Certificate management:

1. First run generates a **local CA** at `~/.orchestra/certs/ca.crt` + `ca.key`
2. Each service gets a signed cert: `{service-name}.crt` + `{service-name}.key`
3. Services load CA cert (verify peers) + own cert (present to peers)
4. `tls.Config.ClientAuth = tls.RequireAndVerifyClientCert`
5. CA fingerprint broadcast via mDNS TXT records for cluster matching
6. Certificates auto-rotate every 90 days

```go
tlsConf := &tls.Config{
    Certificates: []tls.Certificate{serviceCert},
    ClientCAs:    caCertPool,
    RootCAs:      caCertPool,
    ClientAuth:   tls.RequireAndVerifyClientCert,
    NextProtos:   []string{"orch-mesh/1"},
}
listener, err := quic.ListenAddr(":0", tlsConf, &quic.Config{
    EnableDatagrams: true,
})
```

---

## 3. Protobuf Event Envelope

Universal wrapper for ALL inter-service communication. Location: `proto/orchestra/mesh/v1/envelope.proto`

```protobuf
syntax = "proto3";
package orchestra.mesh.v1;

import "google/protobuf/timestamp.proto";
import "google/protobuf/any.proto";

// ================================================================
// Event Envelope — wraps every message on the QUIC mesh
// ================================================================

message Envelope {
  string id = 1;                          // UUIDv7 (time-sortable)
  string topic = 2;                       // Dotted: "task.created", "memory.stored"
  string source = 3;                      // Originating service ID
  google.protobuf.Timestamp timestamp = 4; // Wall-clock at source
  uint64 sequence = 5;                    // Monotonic per source (Lamport-style)
  CausalRef cause = 6;                    // Dependency on another event
  google.protobuf.Any payload = 7;        // Typed payload via Any
  DeliveryMode delivery = 8;              // How to handle offline peers
  uint32 ttl_seconds = 9;                 // 0 = no expiry
  string target = 10;                     // Empty = broadcast to all subscribers
  string trace_parent = 11;              // W3C traceparent for distributed tracing
  map<string, string> metadata = 12;     // project_id, workspace, user_id, etc.
}

message CausalRef {
  string source = 1;
  uint64 sequence = 2;
}

enum DeliveryMode {
  BEST_EFFORT = 0;      // Drop if receiver offline (heartbeats, metrics, UI)
  AT_LEAST_ONCE = 1;    // Store-and-forward via outbox (task state changes)
  REQUEST_REPLY = 2;    // Sender expects response with matching cause (MCP tool calls)
}

// ================================================================
// Control Messages — used on control stream (stream 0)
// ================================================================

message ControlMessage {
  oneof message {
    Handshake handshake = 1;
    HandshakeAck handshake_ack = 2;
    Subscribe subscribe = 3;
    Unsubscribe unsubscribe = 4;
    Heartbeat heartbeat = 5;
    TopologyUpdate topology = 6;
  }
}

message Handshake {
  string service_id = 1;                  // e.g., "orch-mcp"
  string instance_id = 2;                // UUIDv4 per process instance
  string version = 3;                    // semver
  repeated string capabilities = 4;      // ["events", "rpc", "sync"]
  repeated string publish_topics = 5;    // topics this service publishes
  repeated string subscribe_topics = 6;  // topics this service wants
  string cert_fingerprint = 7;           // SHA-256 of service certificate
}

message HandshakeAck {
  bool accepted = 1;
  string reject_reason = 2;
  repeated string matched_topics = 3;    // intersection of publish/subscribe
}

message Subscribe {
  repeated string topics = 1;
}

message Unsubscribe {
  repeated string topics = 1;
}

message Heartbeat {
  string instance_id = 1;
  google.protobuf.Timestamp timestamp = 2;
  uint64 last_sequence = 3;             // highest sequence seen from this peer
  ServiceLoad load = 4;
}

message ServiceLoad {
  float cpu_percent = 1;
  uint64 memory_bytes = 2;
  uint32 active_streams = 3;
  uint32 pending_events = 4;
}

message TopologyUpdate {
  repeated PeerInfo peers = 1;
}

message PeerInfo {
  string service_id = 1;
  string instance_id = 2;
  string address = 3;                   // host:port
  repeated string topics = 4;
  ServiceStatus status = 5;
  google.protobuf.Timestamp last_seen = 6;
}

enum ServiceStatus {
  UNKNOWN = 0;
  HEALTHY = 1;
  DEGRADED = 2;
  OFFLINE = 3;
}
```

**Design rationale:**
- **UUIDv7**: Time-sortable, globally unique, no coordination
- **Lamport sequence**: Monotonic counter per source gives partial ordering without central clock
- **`google.protobuf.Any`**: Wraps any Protobuf message; receivers decode via `type_url`
- **TTL**: Prevents stale events during offline catch-up
- **DeliveryMode**: Most events best-effort; task state changes at-least-once; tool calls request-reply

---

## 4. Mesh Connection Protocol

### 4.1 Connection Lifecycle

```
Service A starts up
  │
  ├──→ 1. Read ~/.orchestra/mesh.json (seed file)
  ├──→ 2. Register on mDNS: _orch-mesh._udp.local.
  ├──→ 3. Browse mDNS for other _orch-mesh._udp.local. services
  │
  For each discovered peer:
  │
  ├──→ 4. Open QUIC connection (mTLS handshake)
  ├──→ 5. Open bidirectional control stream (stream 0)
  ├──→ 6. Send Handshake message
  ├──→ 7. Receive HandshakeAck
  ├──→ 8. If accepted, start heartbeat loop (datagrams every 5s)
  └──→ 9. Begin publishing/receiving events on topic-scoped streams
```

### 4.2 Stream Allocation Strategy

| Stream Type | Purpose | Lifetime |
|-------------|---------|----------|
| Bidirectional stream 0 | Control channel (handshake, subscribe, heartbeat) | Connection lifetime |
| Unidirectional (initiator → peer) | Event delivery (one stream per event) | Single event |
| Bidirectional (for REQUEST_REPLY) | RPC-style calls (tool invocations) | Single request/response |
| QUIC Datagrams | Heartbeats, presence, metrics | Unreliable, no stream |

**Why one stream per event?** QUIC streams are cheap (few bytes of framing):
- No head-of-line blocking between events
- Natural backpressure via QUIC flow control
- Clean lifecycle: open → write → close

For high-throughput scenarios (rapid indexing events), batch multiple envelopes on one unidirectional stream with varint-prefixed length encoding.

### 4.3 Topic Matching

Dotted hierarchical names with wildcards:
- Exact: `task.created` matches `task.created`
- Wildcard: `task.*` matches `task.created`, `task.updated`, `task.deleted`
- Global: `*` matches everything (for logging/debugging)

### 4.4 Connection Ordering

**Rule: Lower service_id connects to higher service_id** (lexicographic). Prevents duplicate connections.

Example:
- `orch-engine` connects TO `orch-gateway`, `orch-mcp`, `orch-server`, `orch-sync`
- `orch-gateway` connects TO `orch-mcp`, `orch-server`, `orch-sync`

If connection drops, the "lower" side reconnects with exponential backoff (1s → 2s → 4s → 8s → max 30s).

### 4.5 Wire Format

Events on QUIC streams: length-delimited Protobuf
```
[4 bytes: big-endian uint32 message length][N bytes: Protobuf-encoded Envelope]
```

Control stream (stream 0): same format with `ControlMessage` protobuf.

QUIC datagrams (heartbeats): raw Protobuf `Heartbeat` bytes (datagrams have own framing).

---

## 5. Service Discovery

### 5.1 Seed File: `~/.orchestra/mesh.json`

```json
{
  "version": 1,
  "ca_fingerprint": "sha256:a1b2c3d4e5f6...",
  "services": {
    "orch-mcp": {
      "enabled": true,
      "address": "127.0.0.1:0",
      "auto_start": true,
      "binary": "orch-mcp",
      "args": ["--workspace", "."],
      "env": {}
    },
    "orch-engine": {
      "enabled": true,
      "address": "127.0.0.1:0",
      "auto_start": true,
      "binary": "orch-engine",
      "args": [],
      "env": {}
    },
    "orch-server": {
      "enabled": true,
      "address": "127.0.0.1:0",
      "auto_start": true,
      "binary": "orch-server",
      "args": [],
      "env": { "ORCHESTRA_PORT": "8443" }
    },
    "orch-sync": {
      "enabled": true,
      "address": "127.0.0.1:0",
      "auto_start": false,
      "binary": "orch-sync",
      "args": [],
      "env": {}
    },
    "orch-gateway": {
      "enabled": true,
      "address": "127.0.0.1:0",
      "auto_start": true,
      "binary": "orch-gateway",
      "args": [],
      "env": {}
    }
  },
  "discovery": {
    "mdns_enabled": true,
    "mdns_domain": "local.",
    "static_peers": []
  },
  "tls": {
    "ca_path": "~/.orchestra/certs/ca.crt",
    "ca_key_path": "~/.orchestra/certs/ca.key",
    "cert_dir": "~/.orchestra/certs/services/",
    "auto_generate": true,
    "rotation_days": 90
  }
}
```

Notes:
- `"address": "127.0.0.1:0"` = localhost, OS-assigned port. Actual port broadcast via mDNS after startup.
- `"auto_start": true` = first Orchestra process to launch spawns this service if not running.
- `"static_peers"` = manually specify addresses for other machines (team dev, cross-subnet).

### 5.2 mDNS Registration

Each service registers:
- **Service type**: `_orch-mesh._udp`
- **Instance name**: `{service_id}-{short_instance_id}` (e.g., `orch-mcp-a1b2c3d4`)
- **Domain**: `local.`
- **Port**: QUIC listener port
- **TXT records**:
  - `v=1` (protocol version)
  - `id={service_id}`
  - `iid={instance_id}` (full UUID)
  - `ca={ca_fingerprint_short}` (first 16 chars SHA-256, for cluster matching)
  - `topics={comma_separated_publish_topics}`
  - `ver={service_version}`

```go
server, err := zeroconf.Register(
    "orch-mcp-a1b2c3d4",
    "_orch-mesh._udp",
    "local.",
    quicPort,
    []string{
        "v=1",
        "id=orch-mcp",
        "iid=a1b2c3d4-...",
        "ca=a1b2c3d4e5f6...",
        "topics=task.*,sprint.*,prd.*",
        "ver=1.0.0",
    },
    nil,
)
```

---

## 6. Event Flow Design

### 6.1 Example: "task.created" Propagation

```
orch-mcp                    orch-server              orch-sync
    │                            │                       │
Agent calls create_task          │                       │
    │                            │                       │
1. Task created in storage files    │                       │
    │                            │                       │
2. Build Envelope:               │                       │
   topic: "task.created"         │                       │
   source: "orch-mcp"            │                       │
   sequence: 42                  │                       │
   delivery: AT_LEAST_ONCE       │                       │
    │                            │                       │
3. Check subscriptions:          │                       │
   orch-server: "task.*" ────────┤                       │
   orch-sync: "task.*" ─────────┤───────────────────────┤
    │                            │                       │
4. Open unidirectional ─────────→│                       │
   stream, write envelope     5. Receives envelope       │
    │                         Updates REST cache          │
6. Open unidirectional ─────────┤──────────────────────→ │
   stream, write envelope        │                    7. Queues for
    │                            │                       cloud sync
```

### 6.2 Offline Handling

**BEST_EFFORT events**: Dropped. Connection detected as down via heartbeat timeout (15s).

**AT_LEAST_ONCE events** (like `task.created`):

1. Detect peer offline (heartbeat timeout)
2. Write event to local outbox: `~/.orchestra/outbox/{peer}/{sequence}.pb`
3. When peer reconnects and completes handshake, it sends `last_sequence` in Heartbeat
4. Replay outbox events with sequence > `last_sequence`
5. Peer acknowledges via updated Heartbeat `last_sequence`
6. Prune acknowledged events from outbox

**Outbox limits**: 1000 events or 50MB per peer (whichever first). Oldest evicted with log warning.

### 6.3 Event Ordering

**Within a single source**: Strictly ordered by `sequence` field. Receivers process in order.

**Across sources**: No total ordering. **Causal ordering** via `CausalRef`:
- When orch-server receives `task.created` (seq 42 from orch-mcp) and publishes `task.notification.sent`, it sets `cause = {source: "orch-mcp", sequence: 42}`
- Any receiver of `task.notification.sent` knows it must first process orch-mcp seq 42

**No vector clocks needed**: With 5 services, tracking per-source sequences (5 integers per receiver) is negligible. Each receiver maintains `map[source]last_seen_sequence` to detect gaps.

---

## 7. Service Designs

### 7.1 `orch-mcp` — MCP Server (The Brain)

Dual interface: stdio JSON-RPC (for agents, per MCP spec) + QUIC mesh.

```
                     stdio (JSON-RPC)
Agent/Claude ──────→ orch-mcp ──────→ QUIC mesh
                       │                 │
                   MCP tools         publishes: task.*, sprint.*, prd.*, workflow.*
                   (186 tools)       subscribes: memory.*, search.*, sync.*
```

### 7.2 `orch-engine` — Rust Data Engine

Existing gRPC services + new QUIC mesh module via `quinn`.

```
Existing gRPC (Tonic)        New QUIC mesh (quinn)
ParseService  ─────┐        publishes: search.indexed, memory.stored, engine.health
SearchService ─────┤─ :50051 subscribes: task.*, file.changed
MemoryService ─────┤
HealthService ─────┘
```

New: LanceDB replaces custom vector storage for embeddings.

### 7.3 `orch-server` — Go API Server

HTTP/3 over QUIC (Fiber v3) + QUIC mesh participant.

```
Browser/Desktop ── HTTP/3 (QUIC) ──→ orch-server ──→ QUIC mesh
                                        │
                                    REST API, Auth, AI Bridge
```

Publishes: `api.request`, `auth.login`, `settings.changed`
Subscribes: `task.*`, `sprint.*`, `workflow.*`, `sync.*`

### 7.4 `orch-sync` — Background Sync Daemon

```
QUIC mesh ──→ orch-sync ──→ PostgreSQL (cloud)
                 │
              Background sync, conflict resolution
```

Subscribes to all state-changing events. Syncs to cloud. Publishes `sync.completed`.

### 7.5 `orch-gateway` — QUIC-to-WebTransport Bridge

```
Browser ── WebTransport (HTTP/3) ──→ orch-gateway ──→ QUIC mesh
```

Authenticates browsers (JWT), maps topic subscriptions, rate-limits, optionally transcodes Protobuf ↔ JSON.

---

## 8. Shared Mesh Library: `libs/go/quic/`

All Go services import this shared package:

```
libs/go/quic/
  mesh.go         — Top-level Mesh struct
  config.go       — Load ~/.orchestra/mesh.json
  certs.go        — CA generation, cert signing, loading
  discovery.go    — mDNS registration + browsing (zeroconf)
  connection.go   — QUIC connection management (quic-go)
  control.go      — Control stream (handshake, subscribe, heartbeat)
  publisher.go    — Event publishing (envelope creation, stream writing)
  subscriber.go   — Event receiving (stream reading, dispatch)
  outbox.go       — Store-and-forward for AT_LEAST_ONCE
  topology.go     — Peer tracking, connection ordering
```

Usage:
```go
m, _ := mesh.New("orch-mcp",
    mesh.WithConfig("~/.orchestra/mesh.json"),
    mesh.WithPublishTopics("task.*", "sprint.*", "prd.*"),
    mesh.WithSubscribeTopics("memory.*", "search.*"),
)
m.Start(ctx)
m.Publish(ctx, "task.created", &taskCreatedPayload, mesh.AtLeastOnce)
m.Subscribe("memory.*", func(env *meshv1.Envelope) { ... })
m.Stop()
```

---

## 9. Practical Considerations

| Concern | Answer |
|---------|--------|
| **Resource usage** | 5 services = max 10 connections (5 choose 2), ~20 goroutines, ~10 UDP sockets. Negligible. |
| **Startup order** | Doesn't matter. Services start independently, discover peers via mDNS. Outbox handles race conditions. |
| **Port conflicts** | All bind to port 0 (OS-assigned), advertise via mDNS. No configuration needed. |
| **Debugging** | `--mesh-debug` flag dumps all traffic to stderr (topic, source, sequence, payload type, size). |
| **Graceful shutdown** | QUIC GOAWAY frame + stream drain. Peers detect disconnect within 5s (one heartbeat). |

---

## 10. Port Allocation

| Service | Protocol | Port | Notes |
|---------|----------|------|-------|
| `orch-mcp` | JSON-RPC stdio | none | Communicates over stdin/stdout |
| `orch-server` | HTTP/3 (QUIC) | 8443 | API + WebSocket upgrade |
| `orch-engine` | gRPC (QUIC) | 50051 | Parse + Search + Memory |
| `orch-gateway` | WebTransport | 4433 | Browser QUIC bridge |
| `orch-sync` | daemon | none | Background process, no listener |
| PostgreSQL | TCP | 5432 | Docker (dev only) |
| Redis | TCP | 6379 | Docker (dev only) |
| ClickHouse | HTTP/Native | 8123/9000 | Docker (dev only) |

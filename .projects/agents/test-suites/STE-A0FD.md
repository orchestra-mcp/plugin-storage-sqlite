---
created_at: "2026-03-01T02:24:35Z"
id: STE-A0FD
name: Summarizer Tests
target_id: AGT-RTHD
target_type: agent
updated_at: "2026-03-01T02:24:40Z"
---

[{"name":"basic summary","prompt":"Summarize: Go is a statically typed language created at Google.","state":"","contains":["Go"],"not_contains":null,"regex":null,"min_length":10},{"name":"no-hallucination check","prompt":"Summarize: Rust is a memory-safe systems language.","state":"","contains":["Rust"],"not_contains":["Python","JavaScript"],"regex":null,"min_length":0}]
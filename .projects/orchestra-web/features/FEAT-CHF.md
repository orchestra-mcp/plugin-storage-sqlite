---
created_at: "2026-02-28T05:13:25Z"
description: |-
    Standalone Next.js 15 app at `apps/web/next/` (apps/next/).

    ## Stack
    - Next.js 15, React 19, TypeScript 5, Tailwind CSS v4, Zustand 5
    - Standalone — no pnpm workspace, own npm install
    - API proxies to Go backend at localhost:8080 via next.config.ts rewrites

    ## Completed (Phase 1 — Small Win)
    - Landing page (`/`) — hero, feature grid, nav, footer
    - Login page (`/login`) — JWT auth, Zustand store, persist to localStorage
    - Register page (`/register`) — create account
    - App shell (`(app)/layout.tsx`) — sidebar, auth guard, user avatar
    - Dashboard (`/dashboard`) — project count + recent notes stats
    - Projects list (`/projects`) — card grid
    - Project detail (`/projects/[id]`) — tree view (epics → stories → tasks)
    - Notes (`/notes`) — list with search
    - Settings (`/settings`) — user profile display
    - Storybook — component stories for Button, Input, Card, Landing, Login, Register pages

    ## Files
    - `apps/next/` — all source files
    - `apps/next/.storybook/` — Storybook config (main.ts, preview.ts)
    - `apps/next/src/components/ui/` — Button, Input, Card with stories

    ## Commands
    ```bash
    cd apps/next && npm install && npm run dev        # http://localhost:3000
    cd apps/next && npm run storybook                  # http://localhost:6006
    make dev-next
    make storybook-next
    ```
id: FEAT-CHF
priority: P0
project_id: orchestra-web
status: backlog
title: Next.js Web App — Landing, Auth & Dashboard
updated_at: "2026-02-28T05:13:25Z"
version: 0
---

# Next.js Web App — Landing, Auth & Dashboard

Standalone Next.js 15 app at `apps/web/next/` (apps/next/).

## Stack
- Next.js 15, React 19, TypeScript 5, Tailwind CSS v4, Zustand 5
- Standalone — no pnpm workspace, own npm install
- API proxies to Go backend at localhost:8080 via next.config.ts rewrites

## Completed (Phase 1 — Small Win)
- Landing page (`/`) — hero, feature grid, nav, footer
- Login page (`/login`) — JWT auth, Zustand store, persist to localStorage
- Register page (`/register`) — create account
- App shell (`(app)/layout.tsx`) — sidebar, auth guard, user avatar
- Dashboard (`/dashboard`) — project count + recent notes stats
- Projects list (`/projects`) — card grid
- Project detail (`/projects/[id]`) — tree view (epics → stories → tasks)
- Notes (`/notes`) — list with search
- Settings (`/settings`) — user profile display
- Storybook — component stories for Button, Input, Card, Landing, Login, Register pages

## Files
- `apps/next/` — all source files
- `apps/next/.storybook/` — Storybook config (main.ts, preview.ts)
- `apps/next/src/components/ui/` — Button, Input, Card with stories

## Commands
```bash
cd apps/next && npm install && npm run dev        # http://localhost:3000
cd apps/next && npm run storybook                  # http://localhost:6006
make dev-next
make storybook-next
```

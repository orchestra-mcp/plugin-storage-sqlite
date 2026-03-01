---
created_at: "2026-02-28T03:34:03Z"
depends_on:
    - FEAT-SHP
description: |-
    Full team management — create teams, invite members, assign roles, share projects, and team-scoped sync.

    **Web Pages**:
    - `resources/js/pages/Teams/Index.tsx` — grid of teams user belongs to, create team button
    - `resources/js/pages/Teams/Show.tsx` — team detail: members list with roles, pending invitations, team settings, danger zone (leave/delete)
    - `resources/js/pages/Teams/AcceptInvitation.tsx` — accept/decline team invitation via signed URL

    **Team CRUD** (`TeamController`):
    - `GET /teams` — list teams (owned + member of)
    - `POST /teams` — create team (name, slug auto-generated, email optional)
    - `GET /teams/{team}` — team detail with members
    - `PUT /teams/{team}` — update name, description, settings
    - `DELETE /teams/{team}` — delete team (owner only, transfers projects to owner first)
    - `POST /teams/{team}/avatar` — upload team avatar (Spatie MediaLibrary → S3/local)

    **Member Management** (`TeamMemberController`):
    - `GET /teams/{team}/members` — paginated member list with roles
    - `PUT /teams/{team}/members/{user}` — update role (owner, admin, member, viewer)
    - `DELETE /teams/{team}/members/{user}` — remove member

    **Invitations** (`TeamInvitationController`):
    - `POST /teams/{team}/invitations` — invite by email, send signed invitation URL (expires 7 days)
    - `GET /invitations/{token}` — show invitation (public, no auth required)
    - `POST /invitations/{token}/accept` — accept (creates team_user record)
    - `POST /invitations/{token}/decline` — decline

    **Sharing** (`TeamShareController`):
    - `GET /teams/{team}/shares` — list share links
    - `POST /teams/{team}/shares` — create share link (expiry, permissions: view|edit)
    - `DELETE /teams/{team}/shares/{share}` — revoke share

    **API endpoints** (for desktop/mobile):
    - Full CRUD mirroring web routes at `/api/teams/`
    - `GET /api/pending-invitations` — invites for current user's email

    **Role permissions matrix**:
    - `owner` — full access, delete team, manage roles
    - `admin` — manage members, invite, create/delete projects
    - `member` — create/edit projects and notes, read all
    - `viewer` — read-only access to team resources

    **Project-level team scope**:
    - Projects with `team_id` set → accessible to all team members per role
    - `Project::isAccessibleBy(User)` checks team membership
    - Sync pull includes team members' changes when `team_id` provided

    Acceptance: team creation works, invitations send + accept via signed URL, roles enforce permissions, project sharing via team_id works, API endpoints usable from desktop/mobile
id: FEAT-NXS
priority: P1
project_id: orchestra-web
status: backlog
title: Team Collaboration (Teams, Members, Roles, Invitations, Sharing)
updated_at: "2026-02-28T03:36:06Z"
version: 0
---

# Team Collaboration (Teams, Members, Roles, Invitations, Sharing)

Full team management — create teams, invite members, assign roles, share projects, and team-scoped sync.

**Web Pages**:
- `resources/js/pages/Teams/Index.tsx` — grid of teams user belongs to, create team button
- `resources/js/pages/Teams/Show.tsx` — team detail: members list with roles, pending invitations, team settings, danger zone (leave/delete)
- `resources/js/pages/Teams/AcceptInvitation.tsx` — accept/decline team invitation via signed URL

**Team CRUD** (`TeamController`):
- `GET /teams` — list teams (owned + member of)
- `POST /teams` — create team (name, slug auto-generated, email optional)
- `GET /teams/{team}` — team detail with members
- `PUT /teams/{team}` — update name, description, settings
- `DELETE /teams/{team}` — delete team (owner only, transfers projects to owner first)
- `POST /teams/{team}/avatar` — upload team avatar (Spatie MediaLibrary → S3/local)

**Member Management** (`TeamMemberController`):
- `GET /teams/{team}/members` — paginated member list with roles
- `PUT /teams/{team}/members/{user}` — update role (owner, admin, member, viewer)
- `DELETE /teams/{team}/members/{user}` — remove member

**Invitations** (`TeamInvitationController`):
- `POST /teams/{team}/invitations` — invite by email, send signed invitation URL (expires 7 days)
- `GET /invitations/{token}` — show invitation (public, no auth required)
- `POST /invitations/{token}/accept` — accept (creates team_user record)
- `POST /invitations/{token}/decline` — decline

**Sharing** (`TeamShareController`):
- `GET /teams/{team}/shares` — list share links
- `POST /teams/{team}/shares` — create share link (expiry, permissions: view|edit)
- `DELETE /teams/{team}/shares/{share}` — revoke share

**API endpoints** (for desktop/mobile):
- Full CRUD mirroring web routes at `/api/teams/`
- `GET /api/pending-invitations` — invites for current user's email

**Role permissions matrix**:
- `owner` — full access, delete team, manage roles
- `admin` — manage members, invite, create/delete projects
- `member` — create/edit projects and notes, read all
- `viewer` — read-only access to team resources

**Project-level team scope**:
- Projects with `team_id` set → accessible to all team members per role
- `Project::isAccessibleBy(User)` checks team membership
- Sync pull includes team members' changes when `team_id` provided

Acceptance: team creation works, invitations send + accept via signed URL, roles enforce permissions, project sharing via team_id works, API endpoints usable from desktop/mobile

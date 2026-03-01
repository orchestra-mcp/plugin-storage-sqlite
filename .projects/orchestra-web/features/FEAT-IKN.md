---
blocks:
    - FEAT-NAA
created_at: "2026-02-28T03:33:46Z"
depends_on:
    - FEAT-SHP
description: |-
    Subscription management — GitHub Sponsors webhook integration, plan enforcement, admin subscription CRUD, and user-facing subscribe/subscription pages.

    **GitHub Sponsors Webhook** (`POST /webhooks/github-sponsors`):
    - Verify HMAC-SHA256 signature from GitHub
    - Handle events: `sponsorship.created` → activate subscription, `sponsorship.cancelled` → expire subscription, `sponsorship.tier_changed` → update plan
    - Map GitHub sponsorship tier to plan: $5/mo → standard, $25/mo → team_sponsor, $50/mo → sponsor
    - Create/update `Subscription` record with `github_sponsor_id`, `amount_cents`, `status`

    **Subscribe page** (`resources/js/pages/subscribe.tsx`):
    - Show available plans (from admin settings: pricing page config)
    - GitHub Sponsors button (link to github.com/sponsors/orchestra-mcp)
    - Manual request form: `POST /api/request-subscription` with message to admin
    - If already subscribed: redirect to /subscription

    **Subscription page** (`resources/js/pages/subscription.tsx`):
    - Current plan Badge, status Badge, dates
    - GitHub Sponsors profile link
    - Cancel/manage link (back to GitHub Sponsors)

    **Admin Subscription management** (`admin/subscriptions/`):
    - Index: DataTable of all subscriptions with user, plan, status, dates, amount
    - Create: manually grant subscription (admin bypass) — user lookup + plan select + dates
    - Edit: change plan, extend dates, toggle status
    - Alerts: subscriptions expiring in 7 days → admin notification

    **Subscription model** (`app/Models/Subscription.php`):
    - `isActive()` — status=active AND end_date > now()
    - Scopes: `active()`, `expired()`, `expiringSoon(int $days)`
    - `plan` enum: sponsor, team_sponsor, standard

    **`subscribed` middleware** (`app/Http/Middleware/EnsureSubscribed.php`):
    - Check `auth()->user()->subscription()->isActive()`
    - Admin role bypass (always passes)
    - Redirect to /subscribe if no active subscription
    - Return 402 for API requests

    **API endpoint** (`POST /api/request-subscription`):
    - Log request with user message
    - Notify admin via email
    - Return 200 with "Request sent" message

    Acceptance: GitHub webhook creates/cancels subscriptions correctly, middleware blocks unsubscribed users, admin can manually grant subscriptions, subscribe page shows plans
id: FEAT-IKN
priority: P1
project_id: orchestra-web
status: backlog
title: Subscription + Billing (GitHub Sponsors, Plan Management)
updated_at: "2026-02-28T03:36:18Z"
version: 0
---

# Subscription + Billing (GitHub Sponsors, Plan Management)

Subscription management — GitHub Sponsors webhook integration, plan enforcement, admin subscription CRUD, and user-facing subscribe/subscription pages.

**GitHub Sponsors Webhook** (`POST /webhooks/github-sponsors`):
- Verify HMAC-SHA256 signature from GitHub
- Handle events: `sponsorship.created` → activate subscription, `sponsorship.cancelled` → expire subscription, `sponsorship.tier_changed` → update plan
- Map GitHub sponsorship tier to plan: $5/mo → standard, $25/mo → team_sponsor, $50/mo → sponsor
- Create/update `Subscription` record with `github_sponsor_id`, `amount_cents`, `status`

**Subscribe page** (`resources/js/pages/subscribe.tsx`):
- Show available plans (from admin settings: pricing page config)
- GitHub Sponsors button (link to github.com/sponsors/orchestra-mcp)
- Manual request form: `POST /api/request-subscription` with message to admin
- If already subscribed: redirect to /subscription

**Subscription page** (`resources/js/pages/subscription.tsx`):
- Current plan Badge, status Badge, dates
- GitHub Sponsors profile link
- Cancel/manage link (back to GitHub Sponsors)

**Admin Subscription management** (`admin/subscriptions/`):
- Index: DataTable of all subscriptions with user, plan, status, dates, amount
- Create: manually grant subscription (admin bypass) — user lookup + plan select + dates
- Edit: change plan, extend dates, toggle status
- Alerts: subscriptions expiring in 7 days → admin notification

**Subscription model** (`app/Models/Subscription.php`):
- `isActive()` — status=active AND end_date > now()
- Scopes: `active()`, `expired()`, `expiringSoon(int $days)`
- `plan` enum: sponsor, team_sponsor, standard

**`subscribed` middleware** (`app/Http/Middleware/EnsureSubscribed.php`):
- Check `auth()->user()->subscription()->isActive()`
- Admin role bypass (always passes)
- Redirect to /subscribe if no active subscription
- Return 402 for API requests

**API endpoint** (`POST /api/request-subscription`):
- Log request with user message
- Notify admin via email
- Return 200 with "Request sent" message

Acceptance: GitHub webhook creates/cancels subscriptions correctly, middleware blocks unsubscribed users, admin can manually grant subscriptions, subscribe page shows plans

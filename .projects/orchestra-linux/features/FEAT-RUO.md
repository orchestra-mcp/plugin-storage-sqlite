---
created_at: "2026-02-28T02:54:16Z"
description: 'SecretService class using libsecret (Secret.password_store/lookup/clear async APIs). Schema: "dev.orchestra.desktop" with "provider" string attribute. Methods: save_api_key(provider, key), load_api_key(provider) → string?, delete_api_key(provider). Providers: claude, openai, gemini, ollama, elevenlabs, deepgram. On first launch: prompt for API keys via AdwPreferencesWindow. Works with GNOME Keyring (org.gnome.keyring) and KDE Wallet (org.kde.kwalletd6) transparently via libsecret.'
id: FEAT-RUO
priority: P1
project_id: orchestra-linux
status: backlog
title: libsecret credential storage (API keys)
updated_at: "2026-02-28T02:54:16Z"
version: 0
---

# libsecret credential storage (API keys)

SecretService class using libsecret (Secret.password_store/lookup/clear async APIs). Schema: "dev.orchestra.desktop" with "provider" string attribute. Methods: save_api_key(provider, key), load_api_key(provider) → string?, delete_api_key(provider). Providers: claude, openai, gemini, ollama, elevenlabs, deepgram. On first launch: prompt for API keys via AdwPreferencesWindow. Works with GNOME Keyring (org.gnome.keyring) and KDE Wallet (org.kde.kwalletd6) transparently via libsecret.

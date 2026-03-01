---
created_at: "2026-02-28T03:14:40Z"
description: |-
    Implement `Orchestra.Desktop/Services/WindowsHelloService.cs` — biometric authentication via Windows Hello (fingerprint / face / PIN) as a second layer before exposing API keys or sensitive project data.

    **`WindowsHelloService` API:**
    ```csharp
    Task<bool> IsAvailableAsync()               // KeyCredentialManager.IsSupportedAsync
    Task<bool> EnrollAsync()                    // KeyCredentialManager.RequestCreateAsync
    Task<bool> AuthenticateAsync(string reason) // KeyCredentialManager.OpenAsync
    Task DeleteCredentialAsync()                // KeyCredentialManager.DeleteAsync
    ```

    **Usage points:**
    - App launch (optional, enabled in Settings → Security → "Require Windows Hello on launch")
    - Before revealing API keys in Settings
    - Before destructive actions (delete project, bulk delete features)

    **`KeyCredentialManager`** — `Windows.Security.Credentials` — uses TPM/secure enclave, works with fingerprint reader, IR camera (Face ID equivalent), or PIN fallback

    **UI:** `ContentDialog` with lock icon, reason text, and biometric prompt. On failure after 3 tries → fall back to Windows account password

    **Platform:** Desktop (TPM 2.0 required for Hello, otherwise disabled gracefully)
id: FEAT-HNF
priority: P3
project_id: orchestra-win
status: backlog
title: Windows Hello — biometric auth
updated_at: "2026-02-28T03:14:40Z"
version: 0
---

# Windows Hello — biometric auth

Implement `Orchestra.Desktop/Services/WindowsHelloService.cs` — biometric authentication via Windows Hello (fingerprint / face / PIN) as a second layer before exposing API keys or sensitive project data.

**`WindowsHelloService` API:**
```csharp
Task<bool> IsAvailableAsync()               // KeyCredentialManager.IsSupportedAsync
Task<bool> EnrollAsync()                    // KeyCredentialManager.RequestCreateAsync
Task<bool> AuthenticateAsync(string reason) // KeyCredentialManager.OpenAsync
Task DeleteCredentialAsync()                // KeyCredentialManager.DeleteAsync
```

**Usage points:**
- App launch (optional, enabled in Settings → Security → "Require Windows Hello on launch")
- Before revealing API keys in Settings
- Before destructive actions (delete project, bulk delete features)

**`KeyCredentialManager`** — `Windows.Security.Credentials` — uses TPM/secure enclave, works with fingerprint reader, IR camera (Face ID equivalent), or PIN fallback

**UI:** `ContentDialog` with lock icon, reason text, and biometric prompt. On failure after 3 tries → fall back to Windows account password

**Platform:** Desktop (TPM 2.0 required for Hello, otherwise disabled gracefully)

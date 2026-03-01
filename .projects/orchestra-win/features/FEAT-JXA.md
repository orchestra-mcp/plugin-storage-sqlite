---
created_at: "2026-02-28T03:19:00Z"
description: |-
    MSIX packaging, signing, and multi-channel distribution pipeline.

    Files:
    - `src/Orchestra.Desktop/Package.appxmanifest` — identity, capabilities, extensions
    - `src/Orchestra.Desktop/Package.wapproj` — Windows Application Packaging project
    - `scripts/sign-msix.ps1` — code signing with EV certificate via signtool.exe
    - `scripts/submit-winget.ps1` — generate winget manifest + submit PR to winget-pkgs
    - `.github/workflows/release-windows.yml` — GitHub Actions release pipeline

    Capabilities declared in manifest:
    - `runFullTrust` — required for QUIC + Win32 P/Invoke
    - `broadFileSystemAccess` — file explorer plugin
    - `webcam`, `microphone` — voice/vision plugins
    - `userNotificationListener` — notification relay

    Distribution channels:
    1. Microsoft Store (MSIX via Partner Center)
    2. winget — `orchestra-mcp.Orchestra` package ID
    3. Chocolatey — `orchestra` package
    4. GitHub Releases — direct MSIX download
    5. Scoop — `orchestra` manifest in `extras` bucket

    GitHub Actions pipeline:
    ```yaml
    - dotnet publish → self-contained
    - msbuild /t:Publish → MSIX bundle
    - signtool sign (EV cert from GitHub Secrets)
    - Upload artifact + create GitHub Release
    - Auto-submit winget PR via winget-create
    ```

    Auto-update: `PackageManager.AddPackageByUriAsync()` polls GitHub Releases API every 24h
id: FEAT-JXA
priority: P2
project_id: orchestra-win
status: backlog
title: MSIX packaging — CI/CD pipeline + Store + winget submission
updated_at: "2026-02-28T03:19:00Z"
version: 0
---

# MSIX packaging — CI/CD pipeline + Store + winget submission

MSIX packaging, signing, and multi-channel distribution pipeline.

Files:
- `src/Orchestra.Desktop/Package.appxmanifest` — identity, capabilities, extensions
- `src/Orchestra.Desktop/Package.wapproj` — Windows Application Packaging project
- `scripts/sign-msix.ps1` — code signing with EV certificate via signtool.exe
- `scripts/submit-winget.ps1` — generate winget manifest + submit PR to winget-pkgs
- `.github/workflows/release-windows.yml` — GitHub Actions release pipeline

Capabilities declared in manifest:
- `runFullTrust` — required for QUIC + Win32 P/Invoke
- `broadFileSystemAccess` — file explorer plugin
- `webcam`, `microphone` — voice/vision plugins
- `userNotificationListener` — notification relay

Distribution channels:
1. Microsoft Store (MSIX via Partner Center)
2. winget — `orchestra-mcp.Orchestra` package ID
3. Chocolatey — `orchestra` package
4. GitHub Releases — direct MSIX download
5. Scoop — `orchestra` manifest in `extras` bucket

GitHub Actions pipeline:
```yaml
- dotnet publish → self-contained
- msbuild /t:Publish → MSIX bundle
- signtool sign (EV cert from GitHub Secrets)
- Upload artifact + create GitHub Release
- Auto-submit winget PR via winget-create
```

Auto-update: `PackageManager.AddPackageByUriAsync()` polls GitHub Releases API every 24h

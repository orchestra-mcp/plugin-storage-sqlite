---
created_at: "2026-02-28T02:55:30Z"
description: 'GitHub Actions workflows for orchestra-linux repo. ci.yml: on push/PR — matrix build (ubuntu-24.04 + fedora-40 containers), meson setup + compile + test (GTest unit tests), Vala compilation check, lint with vala-lint. release.yml: on git tag v* — build all package formats in parallel: Flatpak (flatpak-builder), Snap (snapcraft), .deb (dpkg-buildpackage), .rpm (rpmbuild), AppImage (linuxdeploy), tarball (meson dist). Upload all artifacts to GitHub Release. Notify via D-Bus notification (for developers). Cache: meson build dir, downloaded dependencies.'
id: FEAT-GEH
priority: P2
project_id: orchestra-linux
status: backlog
title: GitHub Actions CI/CD pipeline
updated_at: "2026-02-28T02:55:30Z"
version: 0
---

# GitHub Actions CI/CD pipeline

GitHub Actions workflows for orchestra-linux repo. ci.yml: on push/PR — matrix build (ubuntu-24.04 + fedora-40 containers), meson setup + compile + test (GTest unit tests), Vala compilation check, lint with vala-lint. release.yml: on git tag v* — build all package formats in parallel: Flatpak (flatpak-builder), Snap (snapcraft), .deb (dpkg-buildpackage), .rpm (rpmbuild), AppImage (linuxdeploy), tarball (meson dist). Upload all artifacts to GitHub Release. Notify via D-Bus notification (for developers). Cache: meson build dir, downloaded dependencies.

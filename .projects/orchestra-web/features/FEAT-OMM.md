---
blocks:
    - FEAT-HDD
created_at: "2026-02-28T03:15:42Z"
depends_on:
    - FEAT-ZCC
    - FEAT-SOD
    - FEAT-SLN
    - FEAT-HXU
    - FEAT-SOM
    - FEAT-KIT
    - FEAT-HUY
    - FEAT-KFF
    - FEAT-UFQ
    - FEAT-WBV
    - FEAT-EMV
    - FEAT-SYQ
description: |-
    Wire apps/web/ and libs/plugin-transport-webtransport/ into the build system.

    Makefile changes:
    - Add variable after existing BIN_DIR line: `TRANSPORT_WEBTRANSPORT := $(BIN_DIR)/transport-webtransport`
    - Add to .PHONY: `build-dashboard build-transport-webtransport`
    - Add `build-transport-webtransport` to the main `build:` target dependencies
    - New targets:
      ```
      build-dashboard:
      	cd apps/web && pnpm install && pnpm build
      	rm -rf libs/plugin-transport-webtransport/internal/dist
      	cp -r apps/web/dist libs/plugin-transport-webtransport/internal/dist

      build-transport-webtransport: build-dashboard
      	@mkdir -p $(BIN_DIR)
      	go build -o $(TRANSPORT_WEBTRANSPORT) ./libs/plugin-transport-webtransport/cmd/
      ```
    - Add to test-unit target: `go test ./libs/plugin-transport-webtransport/... -v`
    - Add `transport-webtransport` to the BINARIES install list

    go.work changes:
    - Add `./libs/plugin-transport-webtransport` to the use ( block

    pnpm-workspace.yaml (root-level):
    - Ensure `apps/web` is in the packages list alongside existing apps
    - Pattern: `- 'apps/*'` covers it if already using wildcard

    Acceptance: `make build-transport-webtransport` produces bin/transport-webtransport with embedded React dist, `make test-unit` includes gateway tests, `pnpm install` from root resolves workspace packages in apps/web/
id: FEAT-OMM
priority: P0
project_id: orchestra-web
status: backlog
title: Makefile + go.work + Workspace Config
updated_at: "2026-02-28T03:28:20Z"
version: 0
---

# Makefile + go.work + Workspace Config

Wire apps/web/ and libs/plugin-transport-webtransport/ into the build system.

Makefile changes:
- Add variable after existing BIN_DIR line: `TRANSPORT_WEBTRANSPORT := $(BIN_DIR)/transport-webtransport`
- Add to .PHONY: `build-dashboard build-transport-webtransport`
- Add `build-transport-webtransport` to the main `build:` target dependencies
- New targets:
  ```
  build-dashboard:
  	cd apps/web && pnpm install && pnpm build
  	rm -rf libs/plugin-transport-webtransport/internal/dist
  	cp -r apps/web/dist libs/plugin-transport-webtransport/internal/dist

  build-transport-webtransport: build-dashboard
  	@mkdir -p $(BIN_DIR)
  	go build -o $(TRANSPORT_WEBTRANSPORT) ./libs/plugin-transport-webtransport/cmd/
  ```
- Add to test-unit target: `go test ./libs/plugin-transport-webtransport/... -v`
- Add `transport-webtransport` to the BINARIES install list

go.work changes:
- Add `./libs/plugin-transport-webtransport` to the use ( block

pnpm-workspace.yaml (root-level):
- Ensure `apps/web` is in the packages list alongside existing apps
- Pattern: `- 'apps/*'` covers it if already using wildcard

Acceptance: `make build-transport-webtransport` produces bin/transport-webtransport with embedded React dist, `make test-unit` includes gateway tests, `pnpm install` from root resolves workspace packages in apps/web/

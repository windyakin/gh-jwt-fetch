# CLAUDE.md

## Project Overview

`gh-jwt-fetch` is a single-purpose Go CLI tool that downloads a file from a GitHub repository using GitHub App authentication (JWT → Installation Token → Contents API).

## Build & Test

```bash
go build -o gh-jwt-fetch .
go test -v -race ./...
go vet ./...
```

## Architecture

Flat `package main` structure — no internal packages:

- `config.go` — Environment variable parsing and validation (`LoadConfig`)
- `jwt.go` — RS256 JWT generation using stdlib only (`GenerateJWT`)
- `github.go` — GitHub API client (`GetInstallationToken`, `DownloadFile`)
- `main.go` — Entry point, orchestration, loop mode with signal handling

## Key Design Decisions

- **Zero external dependencies** — JWT RS256 is implemented with `crypto/rsa`, `crypto/x509`, `encoding/pem`
- **Two modes** — One-shot (default) and loop (`GH_INTERVAL` set). Loop mode re-generates JWT and token each iteration
- **GHE support** — `GH_API_BASE_URL` overrides the API endpoint
- **Distroless container** — `gcr.io/distroless/static-debian12:nonroot` with static binary (`CGO_ENABLED=0`)

## Environment Variables

All configuration is via environment variables. See README.md for the full list.
Required: `GH_APP_ID`, `GH_APP_INSTALLATION_ID`, `GH_REPO_OWNER`, `GH_REPO_NAME`, `GH_FILE_PATH`, `GH_OUTPUT_PATH`, and one of `GH_APP_PRIVATE_KEY` / `GH_APP_PRIVATE_KEY_PATH`.

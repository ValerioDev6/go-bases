# Repository Guidelines

## Project Structure & Module Organization

This repository contains the Go module `web`. `cmd/server/main.go` is the server entry point. Configuration lives in `internal/config`; routes and request middleware live in `internal/transport/http`. Shared logging utilities live in `pkg/logger`.

The currently empty `internal/app/controllers`, `internal/app/services`, `internal/domain/user`, and `internal/infrastructure/{db,repository}` directories reserve space for application, domain, and persistence code. Keep HTTP concerns in the transport layer. `internal/transport/http/dtos` and `pkg/validator` are also placeholders. Use `web/templates` and `web/static` for future templates and assets; neither currently contains files.

## Build, Test, and Development Commands

Run commands from the repository root with a toolchain compatible with the `go 1.27.1` directive in `go.mod`:

- `go run ./cmd/server` — start the HTTP server.
- `go build ./...` — compile all packages.
- `go test ./...` — run all package tests.
- `go test -cover ./...` — report test coverage.
- `go vet ./...` — run Go's static checks.
- `gofmt -w cmd internal pkg` — format Go source files.

The entry point creates the logger and injects it into the HTTP router. Keep this dependency explicit when adding handlers or middleware.

## Coding Style & Naming Conventions

Use `gofmt` formatting, including tabs for indentation. Use lowercase package names, exported `PascalCase` identifiers, and unexported `camelCase` identifiers. Prefer the standard library and existing `net/http` and `log/slog` patterns. Alias the local HTTP package when importing it alongside `net/http`. No separate linter configuration is present.

## Testing Guidelines

No tests or coverage threshold are currently defined. Add colocated `*_test.go` files with `TestXxx` functions using Go's `testing` package. Use `net/http/httptest` for routes and middleware, and `t.Setenv` for configuration tests. Cover successful requests, invalid methods, and configuration defaults.

## Commit & Pull Request Guidelines

Git metadata is unavailable in this workspace, so existing commit conventions cannot be verified. Use short, imperative commit subjects, such as `Fix HTTP router initialization`. Keep changes focused. PRs should describe behavior changes, link relevant issues, and report validation results and known failures. Include screenshots when adding visible UI changes.

## Configuration & Secrets

Configuration reads `PORT` and `ENV` from the process environment. Defaults include port `8080` and environment `dev`. The `.env` file is not automatically loaded. Keep credentials out of commits if database or external-service configuration is added later.

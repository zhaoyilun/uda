# Repository Guidelines

## Project Structure & Module Organization
UDA is a Go CLI under `github.com/uda/uda`.

- `main.go` bootstraps the CLI.
- `cmd/` contains command handlers (`create`, `list`, `remove`, `activate`, `deactivate`, `install`, `run`, `self`, `init`).
- `internal/config`: runtime paths (`~/.uda`, `~/.uda/envs`, `~/.uda/config.toml`).
- `internal/env`: environment lifecycle operations.
- `internal/shell`: shell integration script generation (`init`, prompt, activate/deactivate).
- `internal/uv`: uv bootstrap and wrapper execution.
- `docs/`: design notes and architecture references.
- `scripts/build.sh`, `BUILD.md`, `install.sh`: build/install flow.

## Build, Test, and Development Commands
- `./scripts/build.sh` — standard local build, embeds git describe output into version.
- `go build -o uda .` / `go build -ldflags="-s -w" -o uda .` — compile releases.
- `go test ./...` — run all tests (set `GOCACHE=/tmp/go-cache` in restricted envs).
- `go run . --help` — quick command wiring check.
- `./uda init bash|zsh|fish` — print shell integration; eval output in shell startup.

## Coding Style & Naming Conventions
- Keep Go formatted with `gofmt`.
- Use clear early argument validation in `cmd/`.
- Prefer explicit errors and no panics for user flows.
- Use `camelCase` for locals/funcs and `PascalCase` for exported identifiers.
- Shell snippets should be minimal and idempotent (safe to re-run).

## Testing Guidelines
- Add `_test.go` for changed behavior.
- Use temp directories + env overrides when testing filesystem side effects.
- For shell behavior, assert generated script fragments instead of executing real shells when possible.
- Before merging, run tests and at least one smoke path: build + init + create/activate/deactivate.

## Commit & Pull Request Guidelines
- Use Conventional Commit style (`feat:`, `fix:`, `docs:`, `chore:`).
- PR body should include what changed, validation steps, and platform notes (linux/macOS).
- Include sample outputs for user-visible UX changes.

## Security & Configuration Tips
- Never commit binary artifacts (`uda`) or machine-specific paths.
- Keep generated state under `~/.uda`.
- Read mirrors from `UV_MIRROR` / `~/.uda/config.toml` and fail safely when unavailable.

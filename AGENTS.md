# Repository Guidelines

## Project Structure & Module Organization
UDA is a Go CLI in `github.com/uda/uda`.
- `main.go`: program entrypoint.
- `cmd/`: command definitions (`create`, `list`, `remove`, `activate`, `run`, `self`, etc.).
- `internal/config`: runtime config paths (`~/.uda`, `~/.uda/envs`, `~/.uda/config.toml`).
- `internal/env`: environment CRUD helpers.
- `internal/shell`: shell init/activate/deactivate script generation.
- `internal/uv`: uv discovery, install, and wrapper helpers.
- `docs/`: design and implementation notes.
- `scripts/build.sh`: build wrapper and release defaults.
- `BUILD.md` / `install.sh`: user-facing build and install guidance.
- `uda` (root) is a build artifact.

## Build, Test, and Development Commands
- `go test ./...` — run all tests.
- `go build -o uda .` — build a local binary.
- `go build -ldflags="-s -w" -o uda .` — release-style build.
- `./scripts/build.sh` — uses git version metadata when available.
- `go run . --help` — verify CLI wiring quickly.

## Coding Style & Naming Conventions
- Use standard Go formatting (`gofmt`) before commits.
- Use `camelCase` for local variables/functions, `PascalCase` for exported symbols.
- Keep commands in `cmd/` small and explicit; validate required flags/args early.
- Prefer explicit error returns over panics for CLI workflows.
- Shell scripts should keep `set -e` and clear log messages.

## Testing Guidelines
- There are currently no existing tests; add `*_test.go` when touching behavior.
- For command-level tests, keep side effects isolated with temp directories and mocked env vars.
- Before PR: run `go test ./...` and one command smoke test such as `go run . --help`.
- Record any manual checks in PR notes for `UV_MIRROR`, install path, and shell init.

## Commit & Pull Request Guidelines
History uses Conventional Commit style (`feat:`, `fix:`, `chore:`, `docs:`). Match that pattern and keep subjects imperative.
- PRs should include:
  - What changed and why.
  - Commands executed to validate.
  - Notes on platform impact (`linux/macOS` differences).
  - Related issue/task references.
- For UX or output changes, include sample terminal output.

## Security & Configuration Tips
- Do not commit generated files (`uda`) or user-machine paths.
- Changes affecting filesystem writes should avoid hard-coded home paths and validate targets.
- If adding mirror support, validate defaults (`UV_MIRROR`, `~/.uda/config.toml`) and fail safely.

# UDA Documentation

## 1. Runtime Architecture

UDA is a single Go binary that delegates actual Python/venv/package work to `uv`. The code is organized as:

- `cmd/`: CLI command routing (`urfave/cli/v3`) and argument validation.
- `internal/config`: local state under `~/.uda` and runtime paths.
- `internal/uv`: uv discovery, download, and command execution.
- `internal/env`: environment directory operations.
- `internal/shell`: shell helper script generation.

The tool intentionally avoids hidden state outside `~/.uda` and writes minimal side effects to the current shell through script output.

## 2. Runtime Filesystem Layout

- `~/.uda/` base directory
- `~/.uda/envs/` all environments (each env folder is `<name>`)
- `~/.uda/uv` local uv binary
- `~/.uda/config.toml` optional mirror config

## 3. Command Semantics

| Command | Purpose |
|---|---|
| `create <name>` | Create environment folder and call `uv venv`. |
| `list` | List directories under `~/.uda/envs`. |
| `remove <name>` | Remove environment directory recursively. |
| `activate <name>` | Emit `export VIRTUAL_ENV=...` and PATH adjustment commands. |
| `deactivate` | Emit shell cleanup commands for `VIRTUAL_ENV` and PATH. |
| `install` | Run `uv pip install` in selected environment with optional `-r` file. |
| `pip install ...` | Proxied to `uda install` when an environment is active (bash/zsh/fish init). |
| `run` | Run arbitrary command via uv with selected environment python. |
| `self install` | Download and install uv to `~/.uda/uv`, with mirror fallback. |
| `init [bash|zsh|fish]` | Output shell init function/alias script. |

## 4. Mirror Rules

- `UV_MIRROR` env var has highest priority.
- Fallback to `~/.uda/config.toml`:
  ```toml
  [mirror]
  url = "https://pypi.tuna.tsinghua.edu.cn"
  ```
- If unavailable, use built-in mirror list and network test.

## 5. Development Guide

### Build

```bash
./scripts/build.sh
# or
go build -o uda .
```

### Test

```bash
go test ./...
```

### Shell Integration Smoke Test

```bash
eval "$(./uda init bash)"
uda create testenv --python 3.11
uda activate testenv
python --version
deactivate
```

## 6. Release Notes & Compatibility

- Target Go version: `go1.22`.
- Current CLI version tracked in `cmd/root.go`.
- Keep command behavior backwards-compatible in minor releases where possible.

## 7. Known Caveats

- `activate`/`deactivate` output is shell text; when embedding, callers should `eval` command output only as shown in `init`.
- PATH manipulation is intentionally simple and assumes non-empty `VIRTUAL_ENV`.
- Windows paths differ (`Scripts\python.exe`), command behavior still flows through common wrappers.

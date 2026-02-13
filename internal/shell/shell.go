package shell

import (
	"fmt"
	"os"
	"strconv"

	"github.com/uda/uda/internal/config"
)

func Init(shellType string, binaryPath string) string {
	switch shellType {
	case "bash":
		return bashInit(binaryPath)
	case "zsh":
		return zshInit(binaryPath)
	case "fish":
		return fishInit(binaryPath)
	default:
		return bashInit(binaryPath)
	}
}

func bashInit(binaryPath string) string {
	quotedPath := strconv.Quote(binaryPath)
	return fmt.Sprintf(`_UDA_BIN=%s
_UDA_BASE_PS1="${_UDA_BASE_PS1-}"
_UDA_ACTIVE_ENV="${_UDA_ACTIVE_ENV-base}"

_uda_set_prompt() {
    local env_name="$1"
    if [ -z "$env_name" ]; then
        env_name="base"
    fi
    if [ -n "$_UDA_BASE_PS1" ]; then
        PS1="(${env_name}) ${_UDA_BASE_PS1}"
    fi
}

_uda_remove_path_entry() {
    local entry="$1"
    if [ -z "$entry" ]; then
        return
    fi

    local new_path=""
    local old_ifs="$IFS"
    local part
    IFS=":"
    for part in $PATH; do
        if [ "$part" != "$entry" ]; then
            if [ -z "$new_path" ]; then
                new_path="$part"
            else
                new_path="${new_path}:$part"
            fi
        fi
    done
    IFS="$old_ifs"
    PATH="$new_path"
}

if [ -z "$_UDA_BASE_PS1" ]; then
    _UDA_BASE_PS1="${PS1}"
fi

if [ -z "$PS1" ]; then
    PS1=""
fi

_uda_set_prompt "$_UDA_ACTIVE_ENV"

uda() {
    if [ $# -eq 0 ]; then
        "$_UDA_BIN"
        return
    fi

    local cmd="$1"
    shift

    case "$cmd" in
        activate)
            eval "$("$_UDA_BIN" activate "$@")"
            ;;
        deactivate)
            eval "$("$_UDA_BIN" deactivate)"
            ;;
        pip)
            if [ "$1" = "install" ]; then
                uda install "$@"
            else
                command pip "$@"
            fi
            ;;
        pip3)
            if [ "$1" = "install" ]; then
                uda install "$@"
            else
                command pip3 "$@"
            fi
            ;;
        *)
            "$_UDA_BIN" "$cmd" "$@"
            ;;
    esac
}

# Alias for conda compatibility
alias conda=uda
`, quotedPath)
}

func zshInit(binaryPath string) string {
	return bashInit(binaryPath)
}

func fishInit(binaryPath string) string {
	quotedPath := strconv.Quote(binaryPath)
	return fmt.Sprintf(`# UDA fish functions
set -l _UDA_BIN %s

function uda
    if test (count $argv) -eq 0
        eval "$_UDA_BIN"
        return
    end

    set cmd $argv[1]
    set argv $argv[2..-1]
    switch $cmd
        case activate
            eval "$_UDA_BIN activate $argv"
        case deactivate
            eval "$_UDA_BIN deactivate"
        case pip
            if test (count $argv) -gt 0
                if test $argv[1] = install
                    $_UDA_BIN install $argv
                else
                    command pip $argv
                end
            else
                command pip
            end
        case pip3
            if test (count $argv) -gt 0
                if test $argv[1] = install
                    $_UDA_BIN install $argv
                else
                    command pip3 $argv
                end
            else
                command pip3
            end
        case '*'
            $_UDA_BIN $cmd $argv
    end
end

alias conda uda
`, quotedPath)
}

// GenerateActivateScript generates activation commands for a specific environment
func GenerateActivateScript(envName string) (string, error) {
	envPath := config.EnvPath(envName)

	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		return "", fmt.Errorf("environment %s does not exist", envName)
	}

	script := fmt.Sprintf(`if [ -n "$VIRTUAL_ENV" ] && command -v _uda_remove_path_entry >/dev/null 2>&1; then
    _uda_remove_path_entry "$VIRTUAL_ENV/bin"
fi

export VIRTUAL_ENV="%s"
export _UDA_ACTIVE_ENV="%s"
export PATH="$VIRTUAL_ENV/bin:$PATH"
if command -v _uda_set_prompt >/dev/null 2>&1; then
    _uda_set_prompt "$_UDA_ACTIVE_ENV"
fi
`, envPath, envName)

	return script, nil
}

// GenerateDeactivateScript generates deactivation commands
func GenerateDeactivateScript() string {
	return `if [ -n "$VIRTUAL_ENV" ]; then
    if command -v _uda_remove_path_entry >/dev/null 2>&1; then
        _uda_remove_path_entry "$VIRTUAL_ENV/bin"
    else
        export PATH="${PATH#*:$VIRTUAL_ENV/bin}"
    fi
    unset VIRTUAL_ENV
fi

export _UDA_ACTIVE_ENV="base"
if [ -n "$_UDA_BASE_PS1" ] && command -v _uda_set_prompt >/dev/null 2>&1; then
    _uda_set_prompt "$_UDA_ACTIVE_ENV"
fi
`
}

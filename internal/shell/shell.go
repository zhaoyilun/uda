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

	script := fmt.Sprintf(`export VIRTUAL_ENV="%s"
export PATH="$VIRTUAL_ENV/bin:$PATH"
`, envPath)

	return script, nil
}

// GenerateDeactivateScript generates deactivation commands
func GenerateDeactivateScript() string {
	return `if [ -n "$VIRTUAL_ENV" ]; then
    export PATH="${PATH#*:$VIRTUAL_ENV/bin}"
    unset VIRTUAL_ENV
fi
`
}

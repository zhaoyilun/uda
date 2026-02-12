package shell

import (
	"fmt"
	"os"

	"github.com/uda/uda/internal/config"
)

func Init(shellType string) string {
	switch shellType {
	case "bash":
		return bashInit()
	case "zsh":
		return zshInit()
	case "fish":
		return fishInit()
	default:
		return bashInit()
	}
}

func bashInit() string {
	return `uda() {
    local cmd="$1"
    shift
    case "$cmd" in
        activate)
            eval "$(command uda activate "$@")"
            ;;
        deactivate)
            eval "$(command uda deactivate)"
            ;;
        *)
            command uda "$cmd" "$@"
            ;;
    esac
}

# Alias for conda compatibility
alias conda=uda
`
}

func zshInit() string {
	return bashInit()
}

func fishInit() string {
	return `# UDA fish functions
function uda
    set cmd (command uda $argv)
    eval $cmd
end

alias conda uda
`
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
	return `unset VIRTUAL_ENV
export PATH="${PATH#*:$VIRTUAL_ENV/bin}"
`
}

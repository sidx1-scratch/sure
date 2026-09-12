# sure - Shell integration for Bash
# Add to ~/.bashrc: eval "$(sure shell-init --shell bash)"

__sure_enabled=1

sure_disable() {
    __sure_enabled=0
    echo "sure: disabled"
}

sure_enable() {
    __sure_enabled=1
    echo "sure: enabled"
}

__sure_preexec() {
    [ "$__sure_enabled" != "1" ] && return 0
    [ -n "$SURE_SKIP" ] && return 0
    [ "$BASH_COMMAND" = "$PROMPT_COMMAND" ] && return 0

    local cmd="$BASH_COMMAND"

    # Skip if this is the DEBUG trap itself or prompt
    [[ "$cmd" == __sure_* ]] && return 0
    [[ "$cmd" == "$PROMPT_COMMAND" ]] && return 0

    # Skip simple builtins and common harmless commands
    case "${cmd%% *}" in
        cd|pushd|popd|dirs|export|unset|alias|unalias|source|.|history|pwd|exit|logout|true|false|:|echo|printf|type|which|help|set|shopt|complete|compgen|declare|local|readonly|return|shift|trap|wait|jobs|fg|bg|disown|eval|exec|builtin|command|let|read|mapfile|readarray|test|"["|"[["|sure) return 0 ;;
    esac

    # Skip empty commands
    [ -z "$cmd" ] && return 0

    export SURE_SKIP=1
    if ! command sure analyze --command "$cmd" 2>/dev/tty </dev/tty; then
        unset SURE_SKIP
        return 1
    fi
    unset SURE_SKIP
    return 0
}

# Enable extdebug so the DEBUG trap return value controls execution
shopt -s extdebug
trap '__sure_preexec' DEBUG

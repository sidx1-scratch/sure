# sure - Shell integration for Zsh
# Add to ~/.zshrc: eval "$(sure shell-init --shell zsh)"

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
    [ -n "$SURE_SKIP" ] && return

    local cmd="$1"

    # Skip simple builtins
    case "${cmd%% *}" in
        cd|pushd|popd|dirs|export|unset|alias|unalias|source|.|history|pwd|exit|logout|true|false|:|echo|printf|type|which|whence|where|help|set|setopt|unsetopt|autoload|zmodload|declare|local|readonly|return|shift|trap|wait|jobs|fg|bg|disown|eval|exec|builtin|command|let|read|test|"["|"[["|sure|noglob|nocorrect|rehash) return 0 ;;
    esac

    [ -z "$cmd" ] && return 0

    export SURE_SKIP=1
    if ! command sure analyze --command "$cmd" 2>/dev/tty </dev/tty; then
        unset SURE_SKIP
        kill -INT $$
        return 1
    fi
    unset SURE_SKIP
}

autoload -Uz add-zsh-hook
add-zsh-hook preexec __sure_preexec

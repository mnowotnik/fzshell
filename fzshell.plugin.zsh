#!/usr/bin/zsh
export FZSHELL_BIN="${0:a:h}/fzshell"
fzshell_widget() {
    emulate -L zsh
    local completion
    completion=$($FZSHELL_BIN --cursor $CURSOR "$BUFFER")
    if [[ $? != 0 ]]; then
        return 1
    fi
    if [[ -n $completion ]]; then
        LBUFFER="$completion"
        zle reset-prompt
    fi
    return
}

zle -N fzshell_widget
bindkey "^n" fzshell_widget
# vim:ft=zsh:sw=2:

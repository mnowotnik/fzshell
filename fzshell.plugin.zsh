#!/usr/bin/zsh
if [[ -z $FZSHELL_BIN ]]; then
    export FZSHELL_BIN="${0:a:h}/fzshell"
fi
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
if [[ -n $FZSHELL_BIND_KEY ]]; then
    bindkey "$FZSHELL_BIND_KEY" fzshell_widget
else
    bindkey "^n" fzshell_widget
fi
# vim:ft=zsh:sw=2:

#!/usr/bin/zsh
if [[ -z $FZSHELL_BIN ]]; then
    export FZSHELL_BIN="${0:a:h}/fzshell"
fi
fzshell_widget() {
    emulate -L zsh
    local completion
    completion=$($FZSHELL_BIN --cursor $CURSOR "$BUFFER" 2>&1)
    if [[ $? != 0 ]]; then
        zle -I 
        echo fzshell: $completion
        return 1
    fi
    if [[ -z "$completion" ]]; then
        return
    fi

    LBUFFER="$completion"
    zle reset-prompt
}

zle -N fzshell_widget
if [[ -n $FZSHELL_BIND_KEY ]]; then
    bindkey "$FZSHELL_BIND_KEY" fzshell_widget
else
    bindkey "^n" fzshell_widget
fi
# vim:ft=zsh:sw=2:

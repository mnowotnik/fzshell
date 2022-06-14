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

__bind_fzshell_key() {
    bindkey -M emacs "$1" fzshell_widget
    bindkey -M vicmd "$1" fzshell_widget
    bindkey -M viins "$1" fzshell_widget
}

zle -N fzshell_widget
if [[ -n $FZSHELL_BIND_KEY ]]; then
    __bind_fzshell_key "$FZSHELL_BIND_KEY"
else
    __bind_fzshell_key "^N"
fi

unfunction __bind_fzshell_key
# vim:ft=zsh:sw=2:

#!/usr/bin/zsh
export FZSHELL_BIN="${0:a:h}/fzshell"
fzshell_widget() {
    # autoload -U split-shell-arguments
    # local reply REPLY REPLY2
    # split-shell-arguments
    emulate -L zsh
    local completion
    IFS= read -r -d '' completion < <($FZSHELL_BIN "$BUFFER" $CURSOR)
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

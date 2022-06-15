#!/bin/bash
if [[ -z $FZSHELL_BIN ]]; then
    export FZSHELL_BIN="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/fzshell"
fi
fuzzycompl_widget() {
    local completion
    completion=$("$FZSHELL_BIN" --cursor $READLINE_POINT "${READLINE_LINE}")
    if [[ $? != 0 ]]; then
        echo $completion
        return 1
    fi
    if [[ -z $completion ]]; then
        return
    fi
    READLINE_LINE="$completion${READLINE_LINE:$READLINE_POINT}"
    READLINE_POINT=${#completion}
}

__bind_fzshell_key() {
    bind -m vi-command -x "\"${1}\": fuzzycompl_widget"
    bind -m vi-insert -x "\"${1}\": fuzzycompl_widget"
    bind -m emacs-standard -x "\"${1}\": fuzzycompl_widget"
}

if [[ -n $FZSHELL_BIND_KEY ]]; then
    __bind_fzshell_key "$FZSHELL_BIND_KEY"
else
    __bind_fzshell_key "\C-n"
fi

unset -f __bind_fzshell_key
# vim:ft=bash:sw=2:

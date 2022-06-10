#!/bin/bash
if [[ -z $FZSHELL_BIN ]]; then
    export FZSHELL_BIN="$(dirname "${BASH_SOURCE[0]}")/fzshell"
fi
fuzzycompl_widget() {
    local completion=$("$FZSHELL_BIN" --cursor $READLINE_POINT "${READLINE_LINE}")
    READLINE_LINE="XXX"

    local ret=$?
    if [[ $ret != 0 ]]; then
        return $ret
    fi
    if [[ -z $completion ]]; then
        return
    fi
    READLINE_LINE="$completion${READLINE_LINE:$READLINE_POINT}"
    READLINE_POINT=${#completion}
}

if [[ -n $FZSHELL_BIND_KEY ]]; then
    bind -x "\"${FZSHELL_BIND_KEY}\": \"fuzzycompl_widget\""
else
    bind -x '"\C-n": "fuzzycompl_widget"'
fi
# vim:ft=bash:sw=2:

#!/bin/bash
export FZSHELL_BIN="$(dirname "${BASH_SOURCE[0]}")/fzshell"
fuzzycompl_widget() {
    local completion=$("$FZSHELL_BIN" "${READLINE_LINE}" $READLINE_POINT)

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

bind -x '"\C-n": "fuzzycompl_widget"'
# vim:ft=bash:sw=2:

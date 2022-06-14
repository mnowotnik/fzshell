#!/bin/bash
if [[ -z $FZSHELL_BIN ]]; then
    export FZSHELL_BIN="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/fzshell"
fi
fuzzycompl_widget() {
    local ret
    local completion=$("$FZSHELL_BIN" --cursor $READLINE_POINT "${READLINE_LINE}" 2>&1; ret=$?)
    if [[ $ret != 0 ]]; then
        echo fzshell: $completion
        return 1
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

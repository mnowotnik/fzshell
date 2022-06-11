
function fzshell-widget
  set -l cursor (commandline -C)
  set -l line (commandline -b)
  if test -z "$line"
    return 0
  end
  set -l lbuffer ("$FZSHELL_BIN" --cursor "$cursor" "$line")
  set -l cursor  (math $cursor + 1)
  set -l rbuffer (string sub -s$cursor $line)
  commandline -- $lbuffer$rbuffer
  commandline -f repaint
end

if [ -n "$FZSHELL_BIND_KEY" ]
  bind "$FZSHELL_BIND_KEY" fzshell-widget
else
  bind \cn fzshell-widget
end

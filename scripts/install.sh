#!/usr/bin/env bash

# Attribution: https://github.com/junegunn/fzshell/blob/master/install
# License (MIT): https://github.com/junegunn/fzshell/blob/master/LICENSE
# modified for fzshell

set -u
STYLE='\e[1;4m'
NOCOLOR='\033[0m'

version=0.3.4
revision=$(git rev-parse --short HEAD)
cd "$(dirname "${BASH_SOURCE[0]}")/.."
base_dir=$(pwd)
no_instructions=0

help() {
  cat <<EOF

Installation script that either downloads or compiles fzshell binary. Then
adds or prints instructions to add a shell initialization script.

usage: $0 [OPTIONS]
    --help               Show this message
    --no-instructions    Do not print instructions to add lines to your shell config
EOF
}

for opt in "$@"; do
  case $opt in
  --no-instructions)
    no_instructions=1
    ;;
  *)
    echo "Unknown option: $opt"
    help
    exit 1
    ;;
  esac

done

ask() {
  while true; do
    read -p "$1 ([y]/n) " -r
    REPLY=${REPLY:-"y"}
    if [[ $REPLY =~ ^[Yy]$ ]]; then
      return 1
    elif [[ $REPLY =~ ^[Nn]$ ]]; then
      return 0
    fi
  done
}

try_curl() {
  command -v curl >/dev/null &&
    if [[ $1 =~ tar.gz$ ]]; then
      curl -fL $1 | tar -xzf -
    else
      local temp=${TMPDIR:-/tmp}/fzf.zip
      curl -fLo "$temp" $1 && unzip -o "$temp" && rm -f "$temp"
    fi
}

try_wget() {
  command -v wget >/dev/null &&
    if [[ $1 =~ tar.gz$ ]]; then
      wget -O - $1 | tar -xzf -
    else
      local temp=${TMPDIR:-/tmp}/fzf.zip
      wget -O "$temp" $1 && unzip -o "$temp" && rm -f "$temp"
    fi
}

check_binary() {
  echo -n "  - Checking fzshell executable ... "
  local output
  output=$("$base_dir/fzshell" "--version" 2>&1)

  if [ $? -ne 0 ]; then
    echo "Error: $output"
    binary_error="Invalid binary"
  else
    if [ "v$version" != "${output% *}" ]; then
      echo "${output% *} != v$version"
      binary_error="Invalid version"
    else
      echo "$output" OK
      binary_error=""
      return 0
    fi
  fi
  rm -f "$base_dir"/fzshell
  return 1
}

link_fzshell_in_path() {
  if which_fzshell="$(command -v fzshell)"; then
    echo "  - Found in \$PATH"
    echo "  - Creating symlink: fzshell -> $which_fzshell"
    (cd "$sript_base" && rm -f fzshell && ln -sf "$which_fzshell" fzshell)
    check_binary && return
  fi
  return 1
}

download() {
  echo "Downloading fzshell ..."
  if [ -x "$base_dir"/fzshell ]; then
    echo "  - Already exists"
    check_binary && return
  fi
  link_fzshell_in_path && return

  local url
  url=https://github.com/mnowotnik/fzshell/releases/download/v$version/${1}
  set -o pipefail
  if ! (try_curl $url || try_wget $url); then
    set +o pipefail
    binary_error="Failed to download with curl and wget"
    return
  fi
  set +o pipefail

  if [ ! -f fzshell ]; then
    binary_error="Failed to download ${1}"
    return
  fi

  chmod +x fzshell && check_binary
}

archi=$(uname -sm)
binary_available=1
binary_error=""
case "$archi" in
Darwin\ arm64) download fzshell-v$version-darwin-arm64.tar.gz ;;
Darwin\ x86_64) download fzshell-v$version-darwin-amd64.tar.gz ;;
Linux\ armv8*) download fzshell-v$version-linux-arm64.tar.gz ;;
Linux\ aarch64*) download fzshell-v$version-linux-arm64.tar.gz ;;
Linux\ *64) download fzshell-v$version-linux-amd64.tar.gz ;;
FreeBSD\ *64) download fzshell-v$version-freebsd-amd64.tar.gz ;;
OpenBSD\ *64) download fzshell-v$version-openbsd-amd64.tar.gz ;;
CYGWIN*\ *64) download fzshell-v$version-windows-amd64.tar.gz ;;
MINGW*\ *64) download fzshell-v$version-windows-amd64.tar.gz ;;
*\ *64) download fzshell-v$version-windows-amd64.tar.gz ;;
Windows*\ *64) download fzshell-v$version-windows-amd64.tar.gz ;;
*) binary_available=0 binary_error=1 ;;
esac

if [ -n "$binary_error" ]; then
  if [ $binary_available -eq 0 ]; then
    echo "No prebuilt binary for $archi ..."
  else
    echo "  - $binary_error !!!"
  fi
  if command -v go >/dev/null; then
    echo "Attempting to build from source..."
    if go build -ldflags "-s -w -X main.version=$version -X main.revision=$revision"; then
      echo "OK"
    else
      echo "Failed to build binary. Installation failed."
      exit 1
    fi
  else
    echo "go executable not found. Installation failed."
    exit 1
  fi
fi

other_setup() {
  if [[ $no_instructions -eq 1 ]]; then
    return
  fi
  if [[ -n "${BASH_VERSION:-}" ]]; then
    echo "Add the following line to your .bashrc:"
    echo -e ${STYLE}source "\"${base_dir}/fzshell.bash\""$NOCOLOR
    echo
  elif [[ -n "${ZSH_VERSION:-}" ]]; then
    echo "Add the following line to your .zshrc, if you don't use plugin manager:"
    echo -e ${STYLE}source \""${base_dir}/fzshell.plugin.zsh${NOCOLOR}"\"
    echo
  else
    echo "I'm sorry. Your shell is not supported at the moment."
    echo
  fi
}

fish_setup() {
  local fish_dir=${XDG_CONFIG_HOME:-$HOME/.config}/fish/conf.d
  local fish_binding="${fish_dir}/fzshell_key_bindings.fish"
  local fish_binding_src="${base_dir}/fzshell.fish"
  if [[ "$1" -eq 1 ]]; then
    echo "Copying key bindings script to ${fish_binding}..."
    rm -f "$fish_binding"
    cp "${fish_binding_src}" "${fish_binding}" && 
    sed -i "/#to-be-replaced/c\  set FZSHELL_BIN \"$base_dir/fzshell\"" "${fish_binding}" &&
    echo "OK" || 
    echo "Failed"
  else
    if [[ -e "$fish_binding" ]]; then
      echo -n "Removing $fish_binding ... "
      rm -f "$fish_binding"
      echo "OK"
      echo
    fi
    if [[ $no_instructions -eq 1 ]]; then
      return
    fi
    echo Add this line to your config.fish:
    echo -e ${STYLE}source \"${fish_binding_src}\"${NOCOLOR}
    echo
  fi
}
if command -v fish &>/dev/null; then
  ask "Do you want to install key bindings for fish?"
  fish_setup $?
fi

if ! [[ $SHELL =~ fish ]]; then
  other_setup
fi

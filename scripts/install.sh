#!/usr/bin/env bash

# Attribution: https://github.com/junegunn/fzshell/blob/master/install
# License: https://github.com/junegunn/fzshell/blob/master/LICENSE
# modified for fzshell

set -u

version=0.2.0
revision=$(git rev-parse --short HEAD)
cd "$(dirname "${BASH_SOURCE[0]}")/.."
base_dir=$(pwd)

try_curl() {
  command -v curl > /dev/null &&
  curl -# -fL $1 | tar -xzf -
}

try_wget() {
  command -v wget > /dev/null &&
  wget -O - $1 | tar -xzf -
}

check_binary() {
  echo -n "  - Checking fzshell executable ... "
  local output
  output=$("$base_dir/fzshell" "--version" 2>&1)
  if [ $? -ne 0 ]; then
    echo "Error: $output"
    binary_error="Invalid binary"
  else
    if [ "$version" != "${output% *}" ]; then
      echo "$output != $version"
      binary_error="Invalid version"
    else
      echo "$output"
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
  url=https://github.com/mnowotnik/fzshell/releases/download/$version/${1}
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
Darwin\ arm64) download fzshell-$version-darwin-arm64.zip ;;
Darwin\ x86_64) download fzshell-$version-darwin-amd64.zip ;;
Linux\ armv8*) download fzshell-$version-linux-arm64.tar.gz ;;
Linux\ aarch64*) download fzshell-$version-linux-arm64.tar.gz ;;
Linux\ *64) download fzshell-$version-linux-amd64.tar.gz ;;
FreeBSD\ *64) download fzshell-$version-freebsd-amd64.tar.gz ;;
OpenBSD\ *64) download fzshell-$version-openbsd-amd64.tar.gz ;;
CYGWIN*\ *64) download fzshell-$version-windows-amd64.zip ;;
MINGW*\ *64) download fzshell-$version-windows-amd64.zip ;;
MSYS*\ *64) download fzshell-$version-windows-amd64.zip ;;
Windows*\ *64) download fzshell-$version-windows-amd64.zip ;;
*) binary_available=0 binary_error=1 ;;
esac

if [ -n "$binary_error" ]; then
  if [ $binary_available -eq 0 ]; then
    echo "No prebuilt binary for $archi ..."
  else
    echo "  - $binary_error !!!"
  fi
  if command -v go > /dev/null; then
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

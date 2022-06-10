# fzshell

[![Test](https://github.com/mnowotnik/fzshell/actions/workflows/test.yml/badge.svg)](https://github.com/mnowotnik/fzshell/actions/workflows/test.yml)

fzshell is a fuzzy command line completer that fetches completions from sources
predefined by a user. What does it mean? It means that now you can create custom completions for anything you want. All fzshell needs is a pattern to match and
command to generate completion list. It can even insert a completion at any point in a line, not just at the end. See for yourself:

https://user-images.githubusercontent.com/8244123/173105717-e22d4c82-7d38-4c9f-8da3-8b3d75b111d3.mov

This can be accomplished with a few lines:

```yml
completions:
  - pattern: "jq '?(\\.[^']*)'? (\\w+.json)"
    replacement: jq '{{ .item }}' {{ ._2 }}
    cmd: 'jq $1 $2 | jq keys | jq  ". []"'
    map: '{{ ._1 }}{{ printf "[%s]" .item }}'
    preview: >
      {{ shell "jq '" .item "' " ._2 }}
```

## Why?

fzshell was born out of my frustration with performing the same manual tasks
over and over. Like removing obsolete docker containers, deleting kubernetes pods with kubectl or browsing their logs and even copy pasting ticket id from a branch name to a commit message. All of these tasks can be automated at least partially by fzshell. See the [gallery](https://github.com/mnowotnik/fzshell/wiki/Gallery) for use cases.

## Quickstart

### Installation

#### using git

Execute these line to install fzshell on your computer.

```bash
git clone https://github.com/mnowotnik/fzshell ~/.fzshell
cd ~/.fzshell/
./scripts/install.sh
```

Then add initialization to your `.zshrc` or `.bashrc`:

```bash
source ~/.fzshell/fzshell.bash # for bash
```

```bash
source ~/.fzshell/fzshell.plugin.zsh # for zsh
```

#### using plugin manager

If you use package manager like [zplug](https://github.com/zplug/zplug) you
just need to add the following line in your `.zshrc`:

```bash
zplug "mnowotnik/fzshell", hook-build:"./scripts/install.sh"
```

### Basic configuration

fzshell needs a configuration file to load completion definitions.
By default, it loads them from: **~/.config/fzshell/fzshell.yaml**

However, this can be changed by the variable `$FZSHELL_CONFIG` that should
point to a valid configuration.

Below you can see an example configuration file:

```yml
completions:
  - pattern: "docker rmi"
    cmd: "docker images --format '{{.Repository}}:{{.Tag}}\t{{.ID}}'"
    map: ' {{ .item | splitList "\t" | last }}'
    preview: '{{ shell "docker image inspect " .item }}'
```

As you can see the completion definition here has several attributes:

- `pattern`  – regular expression [parsable by Go](https://pkg.go.dev/regexp). It can contain subexpressions (`(xxx)`) and named subexpressions (`(?P<foo>xxx)`)
- `cmd` – Bash shell expression
- `map` – Go [template expression](https://pkg.go.dev/text/template) that has access to [sprig](functions) and subexpression matches:
  - `.item` – whole line returned by the command in `cmd`
  - `._1`,`._2`,... – variables that store non-named subexpression matches
  - `.foo` – named subexpression matches
- `preview` – Just like `map`, but returns a preview of a matched item

Visit [wiki](https://github.com/mnowotnik/fzshell/wiki/Configuration) for a complete configuration guide.

### Usage

## On the way to 1.0.0

fzshell is still in beta, however the specification is unlikely to change, unless
by popular demand. fzshell will advance to 1.0 after user feedback in the coming weeks.

## Disclaimer

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

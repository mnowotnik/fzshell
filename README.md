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

If you find fzshell useful, consider giving it a star. Appreciated!

## Why?

fzshell was born out of my frustration with performing the same manual tasks
over and over. Like removing obsolete docker containers, deleting kubernetes
pods with kubectl or browsing their logs and even copy pasting ticket id from a
branch name to a commit message. 

I tried to to solve this problem in the past using only shell scripts and the result was [docker-fzf-completion](https://github.com/mnowotnik/docker-fzf-completion). However, it was not extensible at all and I had to write a lot of unreadable bash scripts to make it work for any extra command. Additionally, it required from a user more keystrokes than one.

Enter fzshell. All of these tasks I mentioned can be automated at least
partially by fzshell. It divides completions generation into familiar steps,
namely: matching, mapping and filtering. A user only has to provide logic for
those steps and doesn't have to worry about wiring it all together and edge
cases.  Check out the
[gallery of examples](https://github.com/mnowotnik/fzshell/wiki/Examples) to get ideas on how fzshell can help you.

## Want to show your completion definitions 🦚? Need help ❓

Visit [🗨️ Discussions](https://github.com/mnowotnik/fzshell/discussions)!

You can get your questions (probably) answered in [🙏 Questions & Answers ](https://github.com/mnowotnik/fzshell/discussions/categories/questions-answers).

Show us your completion definitions in [🦾 Completions Expo](https://github.com/mnowotnik/fzshell/discussions/categories/completions-expo).

## Quickstart

### Installation

#### using git

Execute these lines to install fzshell on your computer.

```bash
git clone https://github.com/mnowotnik/fzshell ~/.fzshell
cd ~/.fzshell/
./scripts/install.sh
```

Then follow printed instructions.

#### using plugin manager

If you use package manager like [zplug](https://github.com/zplug/zplug) you
just need to add the following line in your `.zshrc`:

```bash
zplug "mnowotnik/fzshell", hook-build:"./scripts/install.sh"
```

### Basic configuration

fzshell needs a configuration file to load completion definitions.
By default, it loads them from: 

**~/.config/fzshell/fzshell.yaml**

However, this can be changed by the variable `$FZSHELL_CONFIG` that should
point to a valid configuration.

Below, you can see an example configuration file:

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

## Usage

The hardest part of using fzshell is writing a correct configuration.
In that the case, all you need to do is press `Ctrl-n` when a cursor is just
after a matching pattern.

Let's consider the example above. Assume the command line looks like this:

```bash
jq . pets.json▉
```

You just need to press `Ctrl-n` to activate fzshell and get a match.
However, if the line deviates from a pattern even a bit nothing will happen.
No match for this line:

```bash
jq . pets.json ▉
```

You would need to modify the pattern a bit to handle extra spaces at the end:

```yml
pattern: "jq '?(\\.[^']*)'? (\\w+.json) *"
```

By default the completion will be inserted at the cursor position, however you
can have complete control over the insertion by defining the `replacement` template. It *replaces* the left part of the line buffer (meaning: to the
left of the cursor). Check [wiki](https://github.com/mnowotnik/fzshell/wiki/Configuration) for more details.

## Development

Development setup is smooth and easy thanks to Go modules.

- requirement: go version 1.18+
- command: `go build`
- manual testing: `source fzshell.bash/fzshell.fish/fzshell.plugin.zsh`
- automatic tests: `make test`
- linting: `make lint` (requires `staticcheck`)
- coverage report: `make cover-report`

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

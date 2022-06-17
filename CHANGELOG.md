# CHANGELOG

## 0.4.2

Fixed bugs in `preview`:

- operations on .item in preview work like in other templates
- listGet actually works
- better errors for listGet and mapGet instead of panic

## 0.4.1

Fixed fzf quirk that forced a user to take into account that preview items are always single quoted.
Now `preview` does not differ from other fields in that regard.

## 0.4.0

Since this is still beta officially this release introduces breaking changes to
the configuration file schema.  All because the finder backend has been swapped
to (slightly modified version of) fzf.  This change automatically adds a couple
of features that were sorely needed like colored output and header parsing.

### Added

- `header` option in completion definition to define a custom sticky header string

### Removed

- `sources` field
- `filter` field
- `view` field

Removed `sources` because multiple sources can  be defined via [process
substitution](https://tldp.org/LDP/abs/html/process-sub.html) in a single
command.

`view` was removed because there is no easy way to transform entries in fzf except for the
`--nth` flag.

`filter` was kind of superfluous and also hard to integrate with fzf.

### Changed

⚠️ Important! ⚠️
`preview` now evaluates to command that gets piped to fzf and executed in
bash. Not to the preview string itself like so far. You need to change it a
bit, but all in all, it's now easier to create previews.

### Fixed

- Cancelling the finder no longer modifies the line buffer

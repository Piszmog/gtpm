# TPM Plugin Manager (in Go)

This is [tpm](https://github.com/tmux-plugins/tpm) reenvisioned in Golang.

## Why?

I spent hours trying to get `tpm` to work but I kept running into errors.

```text
'~/.tmux/plugins/tpm/tpm' returned 127

and

'~/.tmux/plugins/tpm/tpm' returned 1
```

Eventually, I had to do some hacky work arounds to get `tpm` to be happy enough to load my 
plugins.

I was not satified with this and thought I could rewrite it into Go that had better error handling 
and debugging.

### Should I use this?

If `tpm` is working for you, **no**.

If `tpm` is giving you a hard time and you are at the end of your rope, sure give it a spin.

## Installation

TODO

## Installing Plugins

1. Add new plugin to `~/.tmux.conf` with `set -g @plugin '...'`
2. Press `prefix` + <kbd>I</kbd> (capital i, as in **I**nstall) to fetch the plugin.

You're good to go! The plugin was cloned to `~/.tmux/plugins/` dir and sourced.

## Uninstalling Plugins

TODO

## Key Bindings

`prefix` + <kbd>I</kbd>
- Installs new plugins from GitHub or any other git repository
- Refreshes TMUX environment

`prefix` + <kbd>U</kbd>
- updates plugin(s)

`prefix` + <kbd>alt</kbd> + <kbd>u</kbd>
- remove/uninstall plugins not on the plugin list


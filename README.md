# Tmux Plugin Manager (in Go)

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

If `tpm` is giving you a hard time and you are at the end of your rope, sure give it a spin. I have fully migrated to this, and it works great for me.

## Installation

### Requirements

- `tmux` 1.9 or higher
- `git`

### Download and Install

Download `gtpm` for your system by heading over to [Releases](https://github.com/Piszmog/gtpm/releases) and download the artifact for your architecture.

Or you can use [gh](https://cli.github.com/) to download the artifact.

```shell
# Download the latest 64-bit version for Linux
gh release download -R Piszmog/gtpm -p '*Linux_x86_64*'
# Download the latest 64-bit Intel version for macOS
gh release download -R Piszmog/gtpm -p '*Darwin_x86_64*'
# Download the latest Silicon for macOS
gh release download -R Piszmog/gtpm -p '*Darwin_arm64*'

# Untar the artifact
tar -xf gtpm_0.1.0_Linux_x86_64.tar.gz
# Delete the artifact
rm gtpm_0.1.0_Linux_x86_64.tar.gz   
# Move the binary to a directory on your PATH
mv gtpm /some/directory/that/is/in/your/path
```

Ensure `gtpm` is place on your `$PATH` so it can be accessed anywhere.

### Configure tmux

Locate your tmux conf file either at

- `~/.tmux.conf`
- `$XDG_CONFIG_HOME/tmux/tmux.conf`

And add the following to the bottom of the file

```text
# List of plugins
set -g @plugin 'tmux-plugins/tmux-sensible'

# Other examples:
# set -g @plugin 'github_username/plugin_name'
# set -g @plugin 'github_username/plugin_name#branch'
# set -g @plugin 'git@github.com:user/plugin'
# set -g @plugin 'git@bitbucket.com:user/plugin'

# Initialize TMUX plugin manager (keep this line at the very bottom of tmux.conf)
set-environment -g PATH "$PATH:<directory where tmux is installed>:<directory that contains gtpm executable>"
run 'gtpm source'
```

Note: it is important to include where `tmux` and `gtpm` are installed to to `PATH`. If you do not set, you will run into `127` errors 😔.

- `tmux` can be installed to
  - `/run/current-system/sw/bin` if using Nix
  - `/opt/homebrew/bin` if using homebrew

Reload TMUX environment so TPM is sourced:

```bash
# type this in terminal if tmux is already running
tmux source ~/.tmux.conf
# or
tmux source $XDG_CONFIG_HOME/tmux/tmux.conf
```

## Installing Plugins

1. Add new plugin to `~/.tmux.conf` with `set -g @plugin '...'`
2. Press `prefix` + <kbd>I</kbd> (capital i, as in **I**nstall) to fetch the plugin.

You're good to go! The plugin was cloned to `~/.tmux/plugins/` dir and sourced.

## Uninstalling Plugins

1. Remove (or comment out) plugin from the list.
2. Press `prefix` + <kbd>alt</kbd> + <kbd>u</kbd> (lowercase u as in **u**ninstall) to remove the plugin.

All the plugins are installed to `~/.tmux/plugins/` so alternatively you can
find plugin directory there and remove it.

## Key Bindings

`prefix` + <kbd>I</kbd>
- Installs new plugins from GitHub or any other git repository
- Refreshes TMUX environment

`prefix` + <kbd>U</kbd>
- updates plugin(s)

`prefix` + <kbd>alt</kbd> + <kbd>u</kbd>
- remove/uninstall plugins not on the plugin list

## Usage

```shell
$ ./gtpm
```

### Options

| Option                 | Default | Required  | Description                                                                                                                  |
|:-----------------------|:-------:|:---------:|:-----------------------------------------------------------------------------------------------------------------------------|
| `--level`              | `info`  | **False** | Set the logging level. Use `debug` to get more detailed logs.                                                                |
| `--help`, `-h`         | `false` | **False** | Shows help                                                                                                                   |

### Commands

| Command        | Description                                      | Options                                              |
|:---------------|:-------------------------------------------------|:-----------------------------------------------------|
| `clean`, `c`   | Cleans plugins no longer in `tmux` conf file     | N/A                                                  |
| `update`, `u`  | Update plugins                                   | `--plugin value` (repeat) to update specific plugins |
| `install`, `i` | Installs plugins                                 | N/A                                                  |
| `source`, `s`  | Sources plugins and configured key bindings      | N/A                                                  |
| `help`, `h`    | Shows a list of commands or help for one command | N/A                                                  |

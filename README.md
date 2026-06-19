# tmux-pane-status

Tmux pane border status line showing cwd, git info, and running command.

## Installation

    $ go install github.com/yhiraki/tmux-pane-status/cmd/tmux-pane-status@latest

Or download a binary from the [releases page](https://github.com/yhiraki/tmux-pane-status/releases).

## Quick start

Add to `~/.tmux.conf`:

    set -g pane-border-format '#(tmux-pane-status #{pane_current_path} #{pane_pid})'
    set -g pane-border-status top

## Display

| Context | Example |
|---|---|
| Plain directory | ` ~/src ` |
| Git repository | ` 🐱 user/repo  main [M] ` |
| Running command | ` 🐱 user/repo  main [M]  ⟩ vim 0:42` |

The running command is appended as a suffix — git/cwd info is never replaced.

## Configuration

All settings are environment variables. Set them in `~/.tmux.conf` with `set-environment`:

    set-environment -g PANE_STATUS_FORMAT_GIT ' {git_remote_server} {git_repository_name}{git_cwd} {git_current_branch} {git_status_icons} {project_python} '

### Format templates

| Variable | Default |
|---|---|
| `PANE_STATUS_FORMAT_DEFAULT` | ` {cwd} ` |
| `PANE_STATUS_FORMAT_GIT` | ` {git_remote_server} {git_repository_name}{git_cwd} {git_current_branch} {git_status_icons} {project_python} ` |
| `PANE_STATUS_FORMAT_COMMAND` | ` ⟩ {current_command} {current_command_elapsed}` |

Set `PANE_STATUS_FORMAT_COMMAND=''` to disable running-command detection entirely.

#### Available fields

| Field | Description |
|---|---|
| `{cwd}` | Current directory (`~` for home) |
| `{git_remote_server}` | Remote host icon (🐱 github, 🥛 bitbucket, or hostname) |
| `{git_repository_name}` | `user/repo` from remote |
| `{git_cwd}` | Subdirectory relative to git root |
| `{git_current_branch}` | Current branch name |
| `{git_status_icons}` | Changed file indicators, e.g. `[M?]` |
| `{project_python}` | 🐍 when inside a Python project |
| `{current_command}` | Foreground process name |
| `{current_command_elapsed}` | Elapsed time of foreground process |

### Styles (tmux attributes)

| Variable | Default |
|---|---|
| `PANE_STATUS_STYLE_CWD` | `fg=blue` |
| `PANE_STATUS_STYLE_GIT_REMOTE_SERVER` | `bold,fg=blue` |
| `PANE_STATUS_STYLE_GIT_REPOSITORY_NAME` | `bold,fg=blue` |
| `PANE_STATUS_STYLE_GIT_CURRENT_BRANCH` | `bold,fg=magenta` |
| `PANE_STATUS_STYLE_GIT_STATUS_ICONS` | `` |
| `PANE_STATUS_STYLE_GIT_CWD` | `` |
| `PANE_STATUS_STYLE_PROJECT_PYTHON` | `` |
| `PANE_STATUS_STYLE_CURRENT_COMMAND` | `` |
| `PANE_STATUS_STYLE_CURRENT_COMMAND_ELAPSED` | `bold,bg=green,fg=black` |

Set a style to `''` to remove all decoration from that field.

### Icons

| Variable | Default |
|---|---|
| `PANE_STATUS_ICON_GITHUB` | `🐱` |
| `PANE_STATUS_ICON_BITBUCKET` | `🥛` |
| `PANE_STATUS_ICON_PYTHON` | `🐍` |
| `PANE_STATUS_ICON_BRANCH` | `` |

### Other

| Variable | Description |
|---|---|
| `PANE_STATUS_NO_DEFAULTS` | Set to any non-empty value to clear all defaults (useful for building a config from scratch) |

## Migration from v0.x (Python)

The Python-based implementation and its `PY_TMUX_PANE_*` environment variables are no longer supported. Rename your variables:

| Old | New |
|---|---|
| `PY_TMUX_PANE_FORMAT__DEFAULT` | `PANE_STATUS_FORMAT_DEFAULT` |
| `PY_TMUX_PANE_FORMAT__COMMAND` | `PANE_STATUS_FORMAT_COMMAND` |
| `PY_TMUX_PANE_FORMAT__GIT` | `PANE_STATUS_FORMAT_GIT` |
| `PY_TMUX_PANE_OPTIONS__*` | `PANE_STATUS_STYLE_*` |
| `PY_TMUX_PANE_ICON__*` | `PANE_STATUS_ICON_*` |
| `PY_TMUX_PANE_OVERRIDE_DEFAULTS` | `PANE_STATUS_NO_DEFAULTS` |

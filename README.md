# tmux-pane-status

## Installation

    $ pip install git+https://github.com/yhiraki/tmux-pane-status

## Quick start

Add below to your `~/.tmux.conf`:

    set -g pane-border-format '#(tmux-pane-status #{pane_current_path} #{pane_pid})'

## Configuration

    PY_TMUX_PANE_FORMAT__DEFAULT=' {cwd} '
    PY_TMUX_PANE_FORMAT__COMMAND=' {current_command_elapsed} {current_command} '
    PY_TMUX_PANE_FORMAT__GIT=' {git_remote_server} {git_repository_name}{git_cwd} {git_current_branch} {git_status_icons} {project_python} '

# tmux-pane-status

## Installation

    $ pip install git+https://github.com/yhiraki/tmux-pane-status

## Configuration

Add below to your `~/.tmux.conf`:

    set -g pane-border-format '#(tmux-pane-status #{pane_current_path} #{pane_pid})'

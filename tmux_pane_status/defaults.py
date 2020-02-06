#!/usr/bin/env python3
# -*- coding: utf-8 -*-


PY_TMUX_PANE_FORMAT__DEFAULT = ' {cwd} '
PY_TMUX_PANE_FORMAT__GIT = (
    ' {git_remote_server}'
    ' {git_repository_name}'
    ' {git_current_branch}'
    ' {git_status_icons}'
    ' {project_python}'
    ' '
)

PY_TMUX_PANE_OPTIONS__CWD = 'fg=blue'
PY_TMUX_PANE_OPTIONS__GIT_REMOTE_SERVER = 'bold,fg=blue'
PY_TMUX_PANE_OPTIONS__GIT_REPOSITORY_NAME = 'bold,fg=blue'
PY_TMUX_PANE_OPTIONS__GIT_CURRENT_BRANCH = 'bold,fg=magenta'

PY_TMUX_PANE_ICON__PYTHON = '🐍'
PY_TMUX_PANE_ICON__GITHUB = '🐱'
PY_TMUX_PANE_ICON__BITBUCKET = '🥛'

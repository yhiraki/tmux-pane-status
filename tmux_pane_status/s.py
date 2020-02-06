#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import os
import sys
from pathlib import Path

from . import directory
from .value import make_values

PY_TMUX_PANE_FORMAT_DEFAULT = ' {cwd} '
PY_TMUX_PANE_FORMAT_GIT = (
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

options = {
    name.replace('PY_TMUX_PANE_OPTIONS__', '').lower(): [opt.strip() for opt in value.split(',')]
    for name, value in locals().items()
    if name.startswith('PY_TMUX_PANE_OPTIONS__')
}

icons = {
    name.replace('PY_TMUX_PANE_ICON__', '').lower(): value
    for name, value in locals().items()
    if name.startswith('PY_TMUX_PANE_ICON__')
}


def main(cwd, pid):
    s = PY_TMUX_PANE_FORMAT_DEFAULT

    if directory.is_git(cwd):
        s = PY_TMUX_PANE_FORMAT_GIT

    v = make_values(cwd=cwd, options=options, icons=icons)

    # 15% faster
    # from concurrent.futures import ThreadPoolExecutor
    # with ThreadPoolExecutor(max_workers=4) as e:
    #     for key in v:
    #         e.submit(v[key].__str__)

    sys.stdout.write(s.format(**v))


def cli():
    from argparse import ArgumentParser
    parser = ArgumentParser()
    parser.add_argument('cwd', type=Path)
    parser.add_argument('pid', type=int)
    args = parser.parse_args()

    os.chdir(str(args.cwd))
    main(args.cwd.resolve(), args.pid)


if __name__ == '__main__':
    cli()

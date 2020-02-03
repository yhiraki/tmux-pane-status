#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import sys
from pathlib import Path

from . import directory
from .value import make_values

PY_TMUX_PANE_FORMAT_DEFAULT = '{cwd}'
PY_TMUX_PANE_FORMAT_GIT = (
    '{git_remote_server} '
    '{git_repository_name} '
    '{git_current_branch} '
    '{git_status_icons} '
    '{project_python}'
)

PY_TMUX_PANE_OPTIONS__CWD = 'fg=blue'


def main():
    from argparse import ArgumentParser
    parser = ArgumentParser()
    parser.add_argument('cwd', type=Path)
    parser.add_argument('pid', type=int)
    args = parser.parse_args()

    cwd = args.cwd.resolve()
    s = PY_TMUX_PANE_FORMAT_DEFAULT

    if directory.is_git(cwd):
        s = PY_TMUX_PANE_FORMAT_GIT

    v = make_values(cwd=cwd)

    # 15% faster
    # from concurrent.futures import ThreadPoolExecutor
    # with ThreadPoolExecutor(max_workers=4) as e:
    #     for key in v:
    #         e.submit(v[key].__str__)

    sys.stdout.write(s.format(**v))


if __name__ == '__main__':
    main()

#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import os
import sys
from pathlib import Path

from . import pgrep
from . import directory
from .value import make_values

from . import defaults

formats = {}
options = {}
icons = {}
for name, value in vars(defaults).items():
    if os.environ.get(name):
        value = os.environ.get(name)
    if name.startswith('PY_TMUX_PANE_FORMAT__'):
        formats[name.replace('PY_TMUX_PANE_FORMAT__', '').lower()] = value
        continue
    if name.startswith('PY_TMUX_PANE_OPTIONS__'):
        options[name.replace('PY_TMUX_PANE_OPTIONS__', '').lower()] = [
            opt.strip() for opt in value.split(',')]
        continue
    if name.startswith('PY_TMUX_PANE_ICON__'):
        icons[name.replace('PY_TMUX_PANE_ICON__', '').lower()] = value
        continue


def main(cwd, pid):
    s = formats['default']

    pgrep_out = pgrep.p(pid).strip()
    child_pid = 0

    if pgrep_out:
        child_pid = sorted(pgrep_out.split('\n'), reverse=True)[0]
        s = formats['command']

    elif directory.is_git(cwd):
        s = formats['git']

    v = make_values(cwd=cwd, child_pid=child_pid, options=options, icons=icons)

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

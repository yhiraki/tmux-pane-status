#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import sys
from pathlib import Path
from .value import make_values


def is_git(path: Path):
    return git_root(path) is not None


def git_root(path: Path):
    path = Path(path).resolve()
    if not path.is_dir():
        path = path.parent
    while not path == path.parent:
        if (path / '.git').is_dir():
            return path
        path = path.parent


def main():
    from argparse import ArgumentParser
    parser = ArgumentParser()
    parser.add_argument('cwd', type=Path)
    parser.add_argument('pid', type=int)
    args = parser.parse_args()

    v = make_values(cwd=args.cwd)
    out = ''

    if is_git(args.cwd):
        out += '{git_repository_name} {git_current_branch} {git_status_icons}'.format(**v)
    else:
        out += '{cwd}'.format(**v)

    sys.stdout.write(out)


if __name__ == '__main__':
    main()

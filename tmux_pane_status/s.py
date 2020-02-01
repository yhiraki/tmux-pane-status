#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import sys
from pathlib import Path

from .directory import is_git
from .value import make_values


def main():
    from argparse import ArgumentParser
    parser = ArgumentParser()
    parser.add_argument('cwd', type=Path)
    parser.add_argument('pid', type=int)
    args = parser.parse_args()

    v = make_values(cwd=args.cwd)
    out = ''

    if is_git(args.cwd):
        out += (
            '{git_remote_server} '
            '{git_repository_name} '
            '{git_current_branch} '
            '{git_status_icons}'
        ).format(**v)
    else:
        out += '{cwd}'.format(**v)

    sys.stdout.write(out)


if __name__ == '__main__':
    main()

#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import sys
from pathlib import Path

from . import directory
from .value import make_values


def main():
    from argparse import ArgumentParser
    parser = ArgumentParser()
    parser.add_argument('cwd', type=Path)
    parser.add_argument('pid', type=int)
    args = parser.parse_args()

    cwd = args.cwd.resolve()

    out = ''

    if directory.is_git(cwd):
        out += (
            '{git_remote_server} '
            '{git_repository_name} '
            '{git_current_branch} '
            '{git_status_icons}'
        )
    else:
        out += '{cwd}'

    out += ' {project_python}'

    v = make_values(cwd=cwd)
    sys.stdout.write(out.format(**v))


if __name__ == '__main__':
    main()

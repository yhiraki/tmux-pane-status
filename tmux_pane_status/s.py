#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import re
import sys
from os import environ
from pathlib import Path

from . import git


class LazyString:
    def __init__(self, func):
        self.func = func

    def __str__(self):
        return str(self.func())


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
    parser.add_argument('pwd', type=Path)
    parser.add_argument('pid', type=int)
    args = parser.parse_args()

    out = ''

    if is_git(args.pwd):

        for remote in git.remote():
            url = remote[1].split('@')[1]
            m = re.match(r'(.*)[:/](.*)/(.*)', url)
            host, user, name = m.groups()
            out += f'{host}/{user}/{name}'
            break

        out += ' ' + git.current_branch()

        out += f" [{''.join(set(''.join([s[0] for s in git.status()])))}]"

    else:

        pwd = str(args.pwd.resolve())
        home = environ.get('HOME')
        if home:
            pwd = re.sub(f'^{home}', '~', pwd)
        out += f'{pwd}'

    sys.stdout.write(out)


if __name__ == '__main__':
    main()

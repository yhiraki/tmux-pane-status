#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import re
from configparser import ConfigParser
from io import StringIO
from pathlib import Path


class LazyString:
    def __init__(self, func):
        self.func = func

    def __str__(self):
        return str(self.func())


def git_root(path: Path):
    path = Path(path).resolve()
    if not path.is_dir():
        path = path.parent
    while not path == path.parent:
        if (path / '.git').is_dir():
            return path
        path = path.parent


def read_git_config(path):
    path = Path(path)
    with path.open() as f:
        fp = f.readlines()
    config = ConfigParser()
    config.read_file(StringIO(''.join([l.lstrip() for l in fp])))
    return config


def main():
    from argparse import ArgumentParser
    parser = ArgumentParser()
    parser.add_argument('pwd', type=Path)
    parser.add_argument('pid', type=int)
    args = parser.parse_args()

    root = git_root(args.pwd)
    if root:
        config = read_git_config(root / '.git/config')
        for section in config.sections():
            if section.startswith('remote'):
                url = config[section]['url'].split('@')[1]
                m = re.match(r'(.*)[:/](.*)/(.*)', url)
                host, user, name = m.groups()
                print(f'{host}/{user}/{name}')


if __name__ == '__main__':
    main()

#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from pathlib import Path


def _walk_up(path: Path):
    path = Path(path).resolve()
    if not path.is_dir():
        path = path.parent
    while not path == path.parent:
        yield path
        path = path.parent


def is_git(path: Path):
    return git_root(path) is not None


def git_root(path: Path):
    for p in _walk_up(path):
        if (p / '.git').is_dir():
            return p

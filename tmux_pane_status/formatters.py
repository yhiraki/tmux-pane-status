#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from .abc import Formatter
import re
from . import directory
from pathlib import Path


def git_parse_remote(s):
    """Parse git remote

    Args:
        remotes: result of git.remote()
    Returns:
        (name, server, user, name, method) (Tuple[str]):
    """
    for remote in s.split('\n'):
        if not remote.strip():
            continue
        r = remote.split()
        url = r[1].split('@')[1]
        m = re.match(r'(.*)[:/](.*)/(.*)', url)
        yield (r[0], *m.groups(), r[2])


def git_parse_status(s):
    for status in s.split('\n'):
        status = status.strip()
        if not status:
            continue
        yield status.split()


class GitRemoteServer(Formatter):
    def extract_data(self, s):
        for r in git_parse_remote(s):
            return r[1]


class GitRepositoryName(Formatter):
    def extract_data(self, s):
        for r in git_parse_remote(s):
            user, name = r[2], r[3]
            if name.endswith('.git'):
                name = name[:-4]
            return f'{user}/{name}'


class GitCurrentBranch(Formatter):
    def extract_data(self, s):
        return s.strip()


class GitStatusIcons(Formatter):
    def extract_data(self, s):
        v = set(''.join([s[0] for s in git_parse_status(s)]))
        v = ''.join(sorted(v))
        if v:
            return f"[{v}]"
        return ''


class Cwd(Formatter):
    def extract_data(self, s):
        return s


class ProjectPython(Formatter):
    def extract_data(self, s):
        if directory.is_python(Path(s)):
            return 'py'
        return ''

    def set_icons(self, s, icons):
        if s:
            return icons['python']
        return s

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
    def extract_data(self, *ss):
        for r in git_parse_remote(ss[0]):
            return r[1]

    def set_icons(self, s, icons):
        if s == 'github.com':
            return icons['github']
        if s == 'bitbucket.org':
            return icons['bitbucket']
        return s


class GitRepositoryName(Formatter):
    def extract_data(self, *ss):
        for r in git_parse_remote(ss[0]):
            user, name = r[2], r[3]
            if name.endswith('.git'):
                name = name[:-4]
            return f'{user}/{name}'


class GitCurrentBranch(Formatter):
    def extract_data(self, *ss):
        return ss[0].strip()


class GitStatusIcons(Formatter):
    def extract_data(self, *ss):
        v = set(''.join([s[0] for s in git_parse_status(ss[0])]))
        v = ''.join(sorted(v))
        if v:
            return f"[{v}]"
        return ''


class GitCwd(Formatter):
    def extract_data(self, *ss):
        gitroot = ss[0]
        gcwd = ss[1].split(gitroot)
        print(gitroot, gcwd)
        if len(gcwd) < 2:
            return ''
        return gcwd[1]


class Cwd(Formatter):
    def extract_data(self, *ss):
        return ss[0]


class ProjectPython(Formatter):
    def extract_data(self, *ss):
        if directory.is_python(Path(ss[0])):
            return 'py'
        return ''

    def set_icons(self, s, icons):
        if s:
            return icons['python']
        return s

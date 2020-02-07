#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import re
from os import environ
from pathlib import Path

from . import directory
from .abc import Formatter


def git_parse_remote(s):
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

    def set_icons(self, s, icons):
        if icons:
            return f"{icons['branch']} {s}"
        return s


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
        if len(gcwd) < 2:
            return ''
        return gcwd[1]


class Cwd(Formatter):
    def extract_data(self, *ss):
        home = environ.get('HOME')
        s = ss[0]
        if s.startswith(home):
            return s.replace(home, '~')
        if s.startswith('/private'):
            return s.replace('/private', '')
        return s


class ProjectPython(Formatter):
    def extract_data(self, *ss):
        if directory.is_python(Path(ss[0])):
            return 'py'
        return ''

    def set_icons(self, s, icons):
        if s:
            return icons['python']
        return s


class CurrentCommand(Formatter):
    def extract_data(self, *ss):
        cmd = ss[0].split('\n')[1].strip().split()

        if cmd[0] == 'ssh':
            i = 0
            length = len(cmd)
            _cmd = cmd[:]
            cmd = []
            while i < length:
                if _cmd[i] == 'ssh':
                    cmd.append('SSH ->')
                elif _cmd[i] == '-p':
                    i += 1
                elif _cmd[i].startswith('-'):
                    pass
                else:
                    cmd.append(_cmd[i])
                i += 1

        return ' '.join([
            c.split('/')[-1]
            for c in cmd
        ])


class CurrentCommandElapsed(Formatter):
    def extract_data(self, *ss):
        return ss[0].split('\n')[1].strip()

#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from functools import partial

from . import formatters as fmt
from .shell import git, ps
from . import directory


class LazyString:
    def __init__(self, func):
        self.func = func
        self.cache = None

    def __str__(self):
        if self.cache is None:
            self.cache = str(self.func())
        return self.cache


def make_values(*, cwd, child_pid, options, icons):
    def f(formatter, cmds, options=None, icons=None):
        res = [cmd() for cmd in cmds]
        return formatter().format(*res, options=options, icons=icons)

    m = (
        ('git_remote_server', fmt.GitRemoteServer, (git.remote, )),
        ('git_repository_name', fmt.GitRepositoryName, (git.remote,)),
        ('git_current_branch', fmt.GitCurrentBranch, (git.branch_current,)),
        ('git_status_icons', fmt.GitStatusIcons, (git.status,)),
        ('git_cwd', fmt.GitCwd, (partial(str, directory.git_root(cwd)), partial(str, cwd))),
        ('cwd', fmt.Cwd, (partial(str, cwd),)),
        ('project_python', fmt.ProjectPython, (partial(str, cwd), )),
        ('current_command', fmt.CurrentCommand, (partial(ps.p, child_pid, ('command', )),)),
        ('current_command_elapsed', fmt.CurrentCommandElapsed,
         (partial(ps.p, child_pid, ('etime', )),)),
    )
    v = {}
    for name, formatter, cmds in m:
        v[name] = LazyString(
            partial(f, formatter, cmds, options.get(name), icons))
    return v

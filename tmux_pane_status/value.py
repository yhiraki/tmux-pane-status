#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from functools import partial
from . import directory, format, git


class LazyString:
    def __init__(self, func):
        self.func = func
        self.cache = None

    def __str__(self):
        if self.cache is None:
            self.cache = str(self.func())
        return self.cache


def make_values(*, cwd, options):
    def f(fmt, c, opt=None):
        return fmt(c(), options=opt)
    m = (
        ('git_remote_server', git.remote),
        ('git_repository_name', git.remote),
        ('git_current_branch', git.branch_current),
        ('git_status_icons', git.status),
        ('cwd', partial(str, cwd)),
        ('project_python', partial(directory.is_python, cwd)),
    )
    v = {}
    for name, cmd in m:
        formatter = getattr(format, name)
        v[name] = LazyString(partial(f, formatter, cmd, opt=options.get(name)))
    return v

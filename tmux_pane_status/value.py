#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from . import format
from . import git


class LazyString:
    def __init__(self, func):
        self.func = func

    def __str__(self):
        return str(self.func())


def make_values(*, cwd):
    return {
        'git_remote_server': LazyString(lambda: format.git_remote_server(git.remote())),
        'git_repository_name': LazyString(lambda: format.git_repository_name(git.remote())),
        'git_current_branch': LazyString(lambda: format.git_current_branch(git.branch_current())),
        'git_status_icons': LazyString(lambda: format.git_status_icons(git.status())),
        'cwd': LazyString(lambda: format.cwd(str(cwd.resolve())))
    }

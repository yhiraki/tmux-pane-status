#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from . import directory, format, git


class LazyString:
    def __init__(self, func):
        self.func = func
        self.cache = None

    def __str__(self):
        if self.cache is None:
            self.cache = str(self.func())
        return self.cache


def make_values(*, cwd):
    return {
        'git_remote_server': LazyString(lambda: format.git_remote_server(git.remote())),
        'git_repository_name': LazyString(lambda: format.git_repository_name(git.remote())),
        'git_current_branch': LazyString(lambda: format.git_current_branch(git.branch_current())),
        'git_status_icons': LazyString(lambda: format.git_status_icons(git.status())),
        'cwd': LazyString(lambda: format.cwd(str(cwd))),
        'project_python': LazyString(lambda: format.project_python(directory.is_python(cwd)))
    }

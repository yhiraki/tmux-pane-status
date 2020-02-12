#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import subprocess
from typing import Tuple


def run(*cmd: Tuple[str]):
    r = subprocess.run(
        list(cmd),
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE)
    return r.stdout.decode()


class NameSpace:
    def add_func(self, func, name=None):
        if name is None:
            name = func.__name__
        setattr(self, name, func)

    def __call__(self, name=None):
        def _wrapper(func):
            self.add_func(func, name)
        return _wrapper


git = NameSpace()


@git()
def remote():
    return run('git', 'remote', '-v')


@git()
def branch_current():
    return run('git', 'rev-parse', '--abbrev-ref', 'HEAD')


@git()
def status():
    return run('git', 'status', '-s')


ps = NameSpace()


@ps()
def p(pid, headers: Tuple[str]):
    return run('ps', '-p', str(pid), '-o', ','.join(headers))


pgrep = NameSpace()


@pgrep('P')
def _p(pid):
    return run('pgrep', '-P', str(pid))

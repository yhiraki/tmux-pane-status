#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from typing import Tuple
import subprocess


GIT = '/usr/bin/git'


def git(command: str, *options: Tuple[str]):
    r = subprocess.run([GIT, command] + list(options), capture_output=True)
    for line in r.stdout.decode().split('\n'):
        sl = line.strip()
        if sl:
            yield sl


def remote():
    for line in git('remote', '-v'):
        yield line.split()


def branch_current():
    return next(git('rev-parse', '--abbrev-ref', 'HEAD'))


def status():
    for line in git('status', '-s'):
        s = line.split()
        yield s[0], ' '.join(s[1:])

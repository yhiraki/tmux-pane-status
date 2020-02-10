#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from typing import Tuple
import subprocess


GIT = '/usr/bin/git'


def git(command: str, *options: Tuple[str]):
    r = subprocess.run([GIT, command] + list(options),
                       stdout=subprocess.PIPE,
                       stderr=subprocess.PIPE)
    return r.stdout.decode()


def remote():
    return git('remote', '-v')


def branch_current():
    return git('rev-parse', '--abbrev-ref', 'HEAD')


def status():
    return git('status', '-s')

#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import subprocess
from typing import Tuple

PS = '/bin/ps'


def ps(*options: Tuple[str]):
    r = subprocess.run([PS] + list(options), capture_output=True)
    return r.stdout.decode()


def p(pid, headers: Tuple[str]):
    return ps('-p', str(pid), '-o', ','.join(headers))

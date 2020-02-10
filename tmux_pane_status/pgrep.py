#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import subprocess
from typing import Tuple

PGREP = '/usr/bin/pgrep'


def pgrep(*options: Tuple[str]):
    r = subprocess.run([PGREP] + list(options),
                       stdout=subprocess.PIPE,
                       stderr=subprocess.PIPE)
    return r.stdout.decode()


def p(pid):
    return pgrep('-P', str(pid))


if __name__ == '__main__':
    print(p(52250))

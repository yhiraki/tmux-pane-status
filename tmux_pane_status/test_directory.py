# -*- coding: utf-8 -*-

import os
import tempfile
import unittest
from pathlib import Path

from . import directory


class TestStatus(unittest.TestCase):

    def test_git_root(self):
        def helper(files, pwd, gitroot):
            with tempfile.TemporaryDirectory(prefix='/tmp/tmux_pane_status_test_is_git') as d:
                os.chdir(d)
                for file in files:
                    f = Path(file)
                    if file.endswith('/'):
                        f.mkdir()
                        continue
                    f.touch()
                if gitroot is None:
                    self.assertEqual(directory.git_root(pwd), gitroot, (files, pwd, gitroot))
                    return
                self.assertEqual(directory.git_root(pwd), Path(gitroot).resolve(), (files, pwd, gitroot))

        tests = (
            (
                (
                    '.git/',
                ),
                '.',
                '.'
            ),
            (
                (
                    '.git/',
                    'hoge/',
                    'hoge/fuga'
                ),
                'hoge/',
                '.'
            ),
            (
                (
                    'a/',
                    'a/.git/',
                    'a/a/',
                    'a/a/a/',
                    'a/a/a/a/',
                    'a/a/a/a/a/',
                ),
                'a/a/a/a/a/',
                'a/'
            ),
            (
                (
                    'a/',
                    'a/.git/',
                    'a/a/',
                    'a/a/a/',
                    'a/a/a/.git/',
                    'a/a/a/a/',
                    'a/a/a/a/a/',
                ),
                'a/a/a/a/a/',
                'a/a/a/'
            ),
            (
                (
                    'a/',
                    'a/.git'
                ),
                'a/',
                None
            ),
            (
                (
                    '.git'
                ),
                '.',
                None
            ),
            (
                (),
                '.',
                None
            ),
        )
        for files, pwd, gitroot in tests:
            helper(files, pwd, gitroot)


if __name__ == '__main__':
    unittest.main()

#!/usr/bin/env python
# -*- coding: utf-8 -*-

import unittest
from . import formatters


class TestFormatter(unittest.TestCase):

    def setUp(self):
        pass

    def tearDown(self):
        pass

    def test_git_parse_remote(self):
        s = '''\
origin  git@github.com:yhiraki/tmux-pane-status.git (fetch)
origin  git@github.com:yhiraki/tmux-pane-status.git (push)
'''
        gen = formatters.git_parse_remote(s)
        res = next(gen)
        self.assertEqual(res, ('origin', 'github.com', 'yhiraki',
                               'tmux-pane-status.git', '(fetch)'))
        res = next(gen)
        self.assertEqual(res, ('origin', 'github.com', 'yhiraki',
                               'tmux-pane-status.git', '(push)'))
        with self.assertRaises(StopIteration):
            res = next(gen)

    def test_git_parse_status(self):
        s = '''\
 A format.py
 M git.py
 M s.py
 M value.py
?? .envrc
'''
        want = (
            ['A', 'format.py'],
            ['M', 'git.py'],
            ['M', 's.py'],
            ['M', 'value.py'],
            ['??', '.envrc'],
        )
        for i, got in enumerate(formatters.git_parse_status(s)):
            self.assertEqual(want[i], got)


if __name__ == '__main__':
    unittest.main()

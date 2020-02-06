#!/usr/bin/env python3
# -*- coding: utf-8 -*-


class Formatter:

    def set_options(self, s, options):
        if not options:
            return s
        opt = ''.join([f'#[{i}]' for i in options])
        end = '#[default]'
        return f'{opt}{s}{end}'

    def set_icons(self, s, icons):
        return s

    def extract_data(self, s):
        raise NotImplementedError()

    def format(self, s, options, icons):
        s = self.extract_data(s)
        s = self.set_icons(s, icons)
        s = self.set_options(s, options)
        return s

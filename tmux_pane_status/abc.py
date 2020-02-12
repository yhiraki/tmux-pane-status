#!/usr/bin/env python3
# -*- coding: utf-8 -*-


class Formatter:

    def _set_options(self, s, options):
        if not options:
            return s
        opt = ''.join([f'#[{i}]' for i in options])
        end = '#[default]'
        return f'{opt}{s}{end}'

    def _set_icons(self, s, icons):
        return s

    def _extract_data(self, *ss):
        raise NotImplementedError()

    def format(self, *ss, options, icons):
        s = self._extract_data(*ss)
        s = self._set_icons(s, icons)
        s = self._set_options(s, options)
        return s

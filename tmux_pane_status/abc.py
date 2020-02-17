#!/usr/bin/env python3
# -*- coding: utf-8 -*-


class Formatter:
    cache_enabled = False

    def _set_options(self, s: str, options: str, cache=None):
        if not options:
            return s
        opts = [o.strip() for o in options.split(',')]
        self.cache = cache
        opt = ''.join([f'#[{i}]' for i in opts])
        end = '#[default]'
        return f'{opt}{s}{end}'

    def _set_icons(self, s, icons):
        return s

    def _extract_data(self, *ss):
        raise NotImplementedError()

    def _load_cache(self, c: dict):
        raise NotImplementedError()

    def _store_cache(self, c: dict):
        raise NotImplementedError()

    def format(self, *ss, options, icons):
        s = ''
        if self.cache_enabled and self.cache:
            s = self._load_cache()
        else:
            s = self._extract_data(*ss)
            s = self._set_icons(s, icons)
            s = self._set_options(s, options)
        if self.cache_enabled:
            self._store_cache()
        return s

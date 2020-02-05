#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import re
from functools import wraps
from os import environ


def set_options(func):
    @wraps(func)
    def _wrapper(*args, options=None, **kwargs):
        s = func(*args, **kwargs)
        if options is None:
            return s
        opt = ''.join([f'#[{k}={v}]' for k, v in options])
        end = '#[none]'
        return f'{opt}{s}{end}'
    return _wrapper


@set_options
def _git_parse_remote(remotes):
    """Parse git remote

    Args:
        remotes: result of git.remote()
    Returns:
        server, user, name (Tuple[str]):
    """
    err = None
    for remote in remotes:
        try:
            url = remote[1].split('@')[1]
            m = re.match(r'(.*)[:/](.*)/(.*)', url)
            return m.groups()
        except Exception as e:
            err = e
    raise err


@set_options
def git_remote_server(remotes):
    """Git remote server name

    Args:
        remotes: result of git.remote()
    Returns:
        str
    """
    return _git_parse_remote(remotes)[0]


@set_options
def git_repository_name(remotes):
    """Git repository name

    Args:
        remotes: result of git.remote()
    Returns:
        str
    """
    _, user, name = _git_parse_remote(remotes)
    if name.endswith('.git'):
        name = name[:-4]
    return f'{user}/{name}'


@set_options
def git_current_branch(branch_name):
    """Git current branch name

    Args:
        current_branch: result of git.branch_current()
    Returns:
        str
    """
    return branch_name


@set_options
def git_status_icons(statuses):
    """Git status icons

    Args:
        statuses: result of git.status()
    Returns:
        str
    """
    v = set(''.join([s[0] for s in statuses]))
    v = ''.join(sorted(v))
    if v:
        return f"[{v}]"
    return ''


@set_options
def cwd(working_dir):
    """Current working direcgtory

    Args:
        working_dir (str): Full path of cwd
    Returns:
        str
    """
    home = environ.get('HOME')
    if home:
        working_dir = re.sub(f'^{home}', '~', working_dir)
    return working_dir


@set_options
def project_python(yes):
    if yes:
        return 'py'
    return ''

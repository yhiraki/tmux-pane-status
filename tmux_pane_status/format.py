#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import re
from os import environ


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


def git_current_branch(branch_name):
    """Git current branch name

    Args:
        current_branch: result of git.branch_current()
    Returns:
        str
    """
    return branch_name


def git_status_icons(statuses):
    """Git status icons

    Args:
        statuses: result of git.status()
    Returns:
        str
    """
    return f"[{''.join(set(''.join([s[0] for s in statuses])))}]"


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
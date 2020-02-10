#!/bin/bash

set -e

export PY_TMUX_PANE_FORMAT__DEFAULT=' {cwd} '
export PY_TMUX_PANE_FORMAT__COMMAND=' {current_command_elapsed} {current_command} '
export PY_TMUX_PANE_FORMAT__GIT=' {git_remote_server} {git_repository_name}{git_cwd} {git_current_branch} {git_status_icons} {project_python}'
 
export PY_TMUX_PANE_ICON__BITBUCKET='bb'
export PY_TMUX_PANE_ICON__BRANCH=''
export PY_TMUX_PANE_ICON__BRANCH=''
export PY_TMUX_PANE_ICON__GITHUB='gh'
export PY_TMUX_PANE_ICON__PYTHON='py'
export PY_TMUX_PANE_OPTIONS__CURRENT_COMMAND=''
export PY_TMUX_PANE_OPTIONS__CURRENT_COMMAND_ELAPSED=''
export PY_TMUX_PANE_OPTIONS__CWD=''
export PY_TMUX_PANE_OPTIONS__GIT_CURRENT_BRANCH=''
export PY_TMUX_PANE_OPTIONS__GIT_CWD=''
export PY_TMUX_PANE_OPTIONS__GIT_REMOTE_SERVER=''
export PY_TMUX_PANE_OPTIONS__GIT_REPOSITORY_NAME=''
export PY_TMUX_PANE_OPTIONS__GIT_STATUS_ICONS=''
export PY_TMUX_PANE_OPTIONS__PROJECT_PYTHON=''
export PY_TMUX_PANE_OVERRIDE_DEFAULTS=1

function assert () {
  local res=$(tmux-pane-status $(pwd) $$)
  diff <(echo -n "$res") <(echo -n "$1")
}

cd /tmp

REPODIR=.tmux-pane-status/repo
rm -rf $REPODIR
mkdir -p $REPODIR
cd $REPODIR

echo pwd
assert ' /tmp/.tmux-pane-status/repo '

git init > /dev/null
git remote add origin git@github.com:yhiraki/tmux-pane-status
:> README
git add README
git commit -m 'first' > /dev/null

echo github
assert ' gh yhiraki/tmux-pane-status  master  '

git remote remove origin
git remote add origin git@bitbucket.org:yhiraki/tmux-pane-status

echo bitbucket
assert ' bb yhiraki/tmux-pane-status  master  '

sleep 2 &

echo command
assert ' 00:00 sleep 2 '

cd /tmp
rm -rf $REPODIR

echo ---
echo ok

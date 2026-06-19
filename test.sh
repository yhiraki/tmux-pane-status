#!/bin/bash
set -e

BINARY=$(realpath "${1:-tmux-pane-status}")

assert() {
  local got
  got=$(PANE_STATUS_NO_DEFAULTS=1 \
    PANE_STATUS_FORMAT_DEFAULT=' {cwd} ' \
    PANE_STATUS_FORMAT_GIT=' {git_remote_server} {git_repository_name}{git_cwd} {git_current_branch} {git_status_icons} {project_python}' \
    PANE_STATUS_FORMAT_COMMAND=' {current_command_elapsed} {current_command}' \
    PANE_STATUS_ICON_GITHUB='gh' \
    PANE_STATUS_ICON_BITBUCKET='bb' \
    PANE_STATUS_ICON_PYTHON='py' \
    PANE_STATUS_ICON_BRANCH='' \
    "$BINARY" "$(pwd)" $$)
  if ! diff <(printf '%s' "$got") <(printf '%s' "$1"); then
    echo "FAIL: expected '$1', got '$got'" >&2
    exit 1
  fi
}

REPODIR=$(mktemp -d /tmp/tmux-pane-status-test-XXXX)
trap 'rm -rf "$REPODIR"' EXIT

cd "$REPODIR"

echo "==> plain cwd"
assert " $REPODIR "

git -c user.email=t@t.com -c user.name=t init -b master > /dev/null
git remote add origin git@github.com:yhiraki/tmux-pane-status
: > README
git -c user.email=t@t.com -c user.name=t add README
git -c user.email=t@t.com -c user.name=t commit -m first > /dev/null

echo "==> github remote"
assert " gh yhiraki/tmux-pane-status master  "

git remote remove origin
git remote add origin git@bitbucket.org:yhiraki/tmux-pane-status

echo "==> bitbucket remote"
assert " bb yhiraki/tmux-pane-status master  "

echo "==> command suffix"
sleep 2 &
SLEEP_PID=$!
got=$(PANE_STATUS_NO_DEFAULTS=1 \
  PANE_STATUS_FORMAT_DEFAULT=' {cwd} ' \
  PANE_STATUS_FORMAT_GIT=' {git_remote_server}{git_current_branch}' \
  PANE_STATUS_FORMAT_COMMAND=' {current_command_elapsed} {current_command}' \
  PANE_STATUS_ICON_GITHUB='gh' PANE_STATUS_ICON_BRANCH='' \
  "$BINARY" "$(pwd)" $$)
kill "$SLEEP_PID" 2>/dev/null || true
if [[ "$got" != *"sleep"* ]]; then
  echo "FAIL: expected 'sleep' in command output, got '$got'" >&2
  exit 1
fi

echo "==> all ok"

// Package collect runs the external commands (git, ps, pgrep) that feed the
// status fields, memoizing each distinct invocation so it runs at most once.
package collect

import (
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// RunFunc runs a command and returns its stdout. Errors yield "".
type RunFunc func(name string, args ...string) string

func execRun(name string, args ...string) string {
	out, _ := exec.Command(name, args...).Output()
	return string(out)
}

// Collector gathers raw command output for a pane. Commands inherit the
// process working directory, so callers must chdir to the pane cwd first.
type Collector struct {
	Pid int

	run    RunFunc
	getpid func() int

	mu    sync.Mutex
	onces map[string]*sync.Once
	vals  map[string]string
}

// New returns a Collector for the given pane pid.
func New(pid int) *Collector {
	return &Collector{
		Pid:    pid,
		run:    execRun,
		getpid: os.Getpid,
		onces:  map[string]*sync.Once{},
		vals:   map[string]string{},
	}
}

// NewWithRunner returns a Collector with an injected run function (for testing).
func NewWithRunner(pid int, run RunFunc) *Collector {
	return &Collector{
		Pid:    pid,
		run:    run,
		getpid: os.Getpid,
		onces:  map[string]*sync.Once{},
		vals:   map[string]string{},
	}
}

// cached runs f at most once per key. Distinct keys may run concurrently.
func (c *Collector) cached(key string, f func() string) string {
	c.mu.Lock()
	o, ok := c.onces[key]
	if !ok {
		o = &sync.Once{}
		c.onces[key] = o
	}
	c.mu.Unlock()

	o.Do(func() {
		v := f()
		c.mu.Lock()
		c.vals[key] = v
		c.mu.Unlock()
	})

	c.mu.Lock()
	defer c.mu.Unlock()
	return c.vals[key]
}

// GitRemote returns `git remote -v`.
func (c *Collector) GitRemote() string {
	return c.cached("git_remote", func() string { return c.run("git", "remote", "-v") })
}

// GitStatus returns `git status -s`. --no-optional-locks keeps status read-only
// (it does not write back the refreshed index), so this frequent per-pane poll
// never takes .git/index.lock and cannot collide with a concurrent commit/add.
func (c *Collector) GitStatus() string {
	return c.cached("git_status", func() string { return c.run("git", "--no-optional-locks", "status", "-s") })
}

// GitBranch returns the current branch name.
func (c *Collector) GitBranch() string {
	return c.cached("git_branch", func() string {
		return c.run("git", "rev-parse", "--abbrev-ref", "HEAD")
	})
}

// Ps returns `ps -p pid -o headers...` for the given pid.
func (c *Collector) Ps(pid string, headers ...string) string {
	key := "ps:" + pid + ":" + strings.Join(headers, ",")
	return c.cached(key, func() string {
		return c.run("ps", "-p", pid, "-o", strings.Join(headers, ","))
	})
}

// ChildPID returns the foreground child of the pane shell (the highest pid
// among pgrep -P children, excluding this status process itself), or "" when
// there is none.
func (c *Collector) ChildPID() string {
	return c.cached("child_pid", func() string {
		out := strings.TrimSpace(c.run("pgrep", "-P", strconv.Itoa(c.Pid)))
		if out == "" {
			return ""
		}
		pids := strings.Split(out, "\n")
		sort.Sort(sort.Reverse(sort.StringSlice(pids)))
		if pids[0] == strconv.Itoa(c.getpid()) {
			pids = pids[1:]
		}
		if len(pids) == 0 {
			return ""
		}
		return pids[0]
	})
}

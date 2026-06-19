package render

import (
	"regexp"
	"strings"
)

// remote is a parsed `git remote -v` entry.
type remote struct {
	name   string // e.g. "origin"
	server string // e.g. "github.com"
	user   string // e.g. "yhiraki"
	repo   string // e.g. "tmux-pane-status.git"
	kind   string // e.g. "(fetch)"
}

var (
	reRemoteSSH   = regexp.MustCompile(`(.*)[:/](.*)/(.*)`)
	reRemoteHTTPS = regexp.MustCompile(`https://(.*)/(.*)/(.*)`)
)

// gitParseRemote parses the output of `git remote -v` into entries, handling
// both ssh (git@host:user/repo.git) and https (https://host/user/repo.git).
func gitParseRemote(s string) []remote {
	var out []remote
	for _, line := range strings.Split(s, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		r := strings.Fields(line)
		if len(r) < 3 {
			continue
		}
		if strings.Contains(r[1], "@") {
			url := strings.SplitN(r[1], "@", 2)[1]
			if m := reRemoteSSH.FindStringSubmatch(url); m != nil {
				out = append(out, remote{r[0], m[1], m[2], m[3], r[2]})
			}
		}
		if strings.HasPrefix(r[1], "https://") {
			if m := reRemoteHTTPS.FindStringSubmatch(r[1]); m != nil {
				out = append(out, remote{r[0], m[1], m[2], m[3], r[2]})
			}
		}
	}
	return out
}

// gitGetRemote returns the "origin" remote if present, otherwise the first one.
func gitGetRemote(s string) (remote, bool) {
	rs := gitParseRemote(s)
	if len(rs) == 0 {
		return remote{}, false
	}
	for _, r := range rs {
		if r.name == "origin" {
			return r, true
		}
	}
	return rs[0], true
}

// gitParseStatus splits each non-empty line of `git status -s` into fields.
func gitParseStatus(s string) [][]string {
	var out [][]string
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		out = append(out, strings.Fields(line))
	}
	return out
}

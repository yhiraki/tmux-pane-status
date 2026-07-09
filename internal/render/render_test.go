package render

import (
	"strings"
	"testing"

	"github.com/yhiraki/tmux-pane-status/internal/collect"
	"github.com/yhiraki/tmux-pane-status/internal/config"
)

// mockCollector builds a Collector whose run func is backed by a map of
// "cmd arg arg..." → output.
func mockCollector(responses map[string]string) *collect.Collector {
	c := collect.NewWithRunner(1, func(name string, args ...string) string {
		key := name + " " + strings.Join(args, " ")
		return responses[key]
	})
	return c
}

var defaultCfg = &config.Config{
	Formats: map[string]string{
		"default": " {cwd} ",
		"command": " > {current_command} {current_command_elapsed}",
		"git":     " {git_remote_server} {git_current_branch} {git_status_icons} ",
	},
	Styles: map[string]string{},
	Icons:  map[string]string{"github": "GH", "bitbucket": "BB", "branch": ""},
}

func TestRenderDefault(t *testing.T) {
	col := mockCollector(map[string]string{})
	pane := Pane{Cwd: "/home/user/projects", Home: "/home/user"}
	out := Render(pane, col, defaultCfg)
	if out != " ~/projects " {
		t.Errorf("Render = %q, want \" ~/projects \"", out)
	}
}

func TestRenderGit(t *testing.T) {
	col := mockCollector(map[string]string{
		"git remote -v":                   "origin\tgit@github.com:user/repo.git (fetch)\n",
		"git rev-parse --abbrev-ref HEAD": "main\n",
		"git --no-optional-locks status -s": " M foo.go\n",
	})
	pane := Pane{Cwd: "/repo", Home: "/home/user", GitRoot: "/repo"}
	out := Render(pane, col, defaultCfg)
	if !strings.Contains(out, "GH") {
		t.Errorf("Render git missing remote server icon, got %q", out)
	}
	if !strings.Contains(out, "main") {
		t.Errorf("Render git missing branch, got %q", out)
	}
	if !strings.Contains(out, "[M]") {
		t.Errorf("Render git missing status icons, got %q", out)
	}
}

func TestRenderCommandSuffix(t *testing.T) {
	col := mockCollector(map[string]string{
		"ps -p 42 -o command": "COMMAND\nvim file.txt\n",
		"ps -p 42 -o etime":   "ELAPSED\n  0:30\n",
	})
	pane := Pane{Cwd: "/home/user", Home: "/home/user", ChildPID: "42"}
	out := Render(pane, col, defaultCfg)
	if !strings.Contains(out, "vim") {
		t.Errorf("Render command suffix missing command, got %q", out)
	}
	if !strings.Contains(out, "0:30") {
		t.Errorf("Render command suffix missing elapsed, got %q", out)
	}
}

func TestRenderEmptyCommandFormatSkipsPs(t *testing.T) {
	called := false
	col := collect.NewWithRunner(1, func(name string, args ...string) string {
		if name == "ps" || name == "pgrep" {
			called = true
		}
		return ""
	})
	cfg := &config.Config{
		Formats: map[string]string{
			"default": " {cwd} ",
			"command": "",
			"git":     " {cwd} ",
		},
		Styles: map[string]string{},
		Icons:  map[string]string{},
	}
	pane := Pane{Cwd: "/tmp", Home: "/home/user"}
	Render(pane, col, cfg)
	if called {
		t.Error("ps/pgrep was called despite empty FORMAT_COMMAND")
	}
}

func TestExtractFields(t *testing.T) {
	fields := extractFields("{cwd} {git_current_branch} {cwd}")
	if len(fields) != 2 {
		t.Errorf("extractFields len = %d, want 2: %v", len(fields), fields)
	}
	if !fields["cwd"] || !fields["git_current_branch"] {
		t.Errorf("extractFields missing expected keys: %v", fields)
	}
}

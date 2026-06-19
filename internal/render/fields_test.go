package render

import "testing"

func TestGitRemoteServer(t *testing.T) {
	icons := map[string]string{"github": "🐱", "bitbucket": "🥛"}
	tests := []struct {
		raw  string
		want string
	}{
		{"origin\tgit@github.com:user/repo.git (fetch)\n", "🐱"},
		{"origin\tgit@bitbucket.org:user/repo.git (fetch)\n", "🥛"},
		{"origin\tgit@gitlab.com:user/repo.git (fetch)\n", "gitlab.com"},
		{"", ""},
	}
	for _, tt := range tests {
		got := gitRemoteServer(tt.raw, icons)
		if got != tt.want {
			t.Errorf("gitRemoteServer(%q) = %q, want %q", tt.raw, got, tt.want)
		}
	}
}

func TestGitRepositoryName(t *testing.T) {
	tests := []struct {
		raw  string
		want string
	}{
		{"origin\tgit@github.com:yhiraki/tmux-pane-status.git (fetch)\n", "yhiraki/tmux-pane-status"},
		{"origin\thttps://github.com/yhiraki/repo (fetch)\n", "yhiraki/repo"},
		{"", ""},
	}
	for _, tt := range tests {
		got := gitRepositoryName(tt.raw)
		if got != tt.want {
			t.Errorf("gitRepositoryName(%q) = %q, want %q", tt.raw, got, tt.want)
		}
	}
}

func TestGitCurrentBranch(t *testing.T) {
	// non-empty icon: prepend icon + space
	got := gitCurrentBranch("main\n", map[string]string{"branch": "B"})
	if got != "B main" {
		t.Errorf("gitCurrentBranch (with icon) = %q, want \"B main\"", got)
	}

	// empty icon: return branch name only
	got2 := gitCurrentBranch("main\n", map[string]string{"branch": ""})
	if got2 != "main" {
		t.Errorf("gitCurrentBranch (no icon) = %q, want \"main\"", got2)
	}
}

func TestGitStatusIcons(t *testing.T) {
	tests := []struct {
		raw  string
		want string
	}{
		{" A fmt.go\n M git.go\n M s.go\n?? .envrc\n", "[?AM]"},
		{"", ""},
	}
	for _, tt := range tests {
		got := gitStatusIcons(tt.raw)
		if got != tt.want {
			t.Errorf("gitStatusIcons(%q) = %q, want %q", tt.raw, got, tt.want)
		}
	}
}

func TestGitCwd(t *testing.T) {
	tests := []struct {
		root string
		cwd  string
		want string
	}{
		{"/repo", "/repo/sub/dir", "/sub/dir"},
		{"/repo", "/repo", ""},
		{"/repo", "/other", ""},
	}
	for _, tt := range tests {
		got := gitCwd(tt.root, tt.cwd)
		if got != tt.want {
			t.Errorf("gitCwd(%q, %q) = %q, want %q", tt.root, tt.cwd, got, tt.want)
		}
	}
}

func TestCwd(t *testing.T) {
	tests := []struct {
		cwd  string
		home string
		want string
	}{
		{"/home/user/src", "/home/user", "~/src"},
		{"/private/tmp/foo", "/root", "/tmp/foo"},
		{"/var/log", "/home/user", "/var/log"},
	}
	for _, tt := range tests {
		got := cwdField(tt.cwd, tt.home)
		if got != tt.want {
			t.Errorf("cwdField(%q, %q) = %q, want %q", tt.cwd, tt.home, got, tt.want)
		}
	}
}

func TestCurrentCommand(t *testing.T) {
	tests := []struct {
		psOut string
		want  string
	}{
		{"COMMAND\nvim /home/user/file.txt\n", "vim file.txt"},
		{"COMMAND\nssh -p 2222 user@server.com\n", "SSH -> user@server.com"},
		{"COMMAND\n", ""},
		{"", ""},
	}
	for _, tt := range tests {
		got := currentCommand(tt.psOut)
		if got != tt.want {
			t.Errorf("currentCommand(%q) = %q, want %q", tt.psOut, got, tt.want)
		}
	}
}

func TestCurrentCommandElapsed(t *testing.T) {
	got := currentCommandElapsed("ELAPSED\n  1:23\n")
	if got != "1:23" {
		t.Errorf("currentCommandElapsed = %q, want \"1:23\"", got)
	}
}

func TestApplyStyle(t *testing.T) {
	got := applyStyle("main", "bold,fg=magenta")
	want := "#[bold]#[fg=magenta]main#[default]"
	if got != want {
		t.Errorf("applyStyle = %q, want %q", got, want)
	}
	if applyStyle("main", "") != "main" {
		t.Error("applyStyle with empty style should return value unchanged")
	}
}

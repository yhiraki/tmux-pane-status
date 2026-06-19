package render

import (
	"reflect"
	"testing"
)

func TestGitParseRemoteSSH(t *testing.T) {
	s := "origin\tgit@github.com:yhiraki/tmux-pane-status.git (fetch)\n" +
		"origin\tgit@github.com:yhiraki/tmux-pane-status.git (push)\n"
	got := gitParseRemote(s)
	want := []remote{
		{"origin", "github.com", "yhiraki", "tmux-pane-status.git", "(fetch)"},
		{"origin", "github.com", "yhiraki", "tmux-pane-status.git", "(push)"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("gitParseRemote = %+v, want %+v", got, want)
	}
}

func TestGitParseRemoteHTTPS(t *testing.T) {
	s := "origin\thttps://github.com/yhiraki/tmux-pane-status.git (fetch)\n" +
		"origin\thttps://github.com/yhiraki/tmux-pane-status.git (push)\n"
	got := gitParseRemote(s)
	want := []remote{
		{"origin", "github.com", "yhiraki", "tmux-pane-status.git", "(fetch)"},
		{"origin", "github.com", "yhiraki", "tmux-pane-status.git", "(push)"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("gitParseRemote = %+v, want %+v", got, want)
	}
}

func TestGitGetRemotePrefersOrigin(t *testing.T) {
	s := "upstream\tgit@github.com:other/repo.git (fetch)\n" +
		"origin\tgit@github.com:yhiraki/repo.git (fetch)\n"
	got, ok := gitGetRemote(s)
	if !ok {
		t.Fatal("gitGetRemote ok = false, want true")
	}
	if got.name != "origin" || got.user != "yhiraki" {
		t.Errorf("gitGetRemote = %+v, want origin/yhiraki", got)
	}
}

func TestGitGetRemoteFallsBackToFirst(t *testing.T) {
	s := "upstream\tgit@github.com:other/repo.git (fetch)\n"
	got, ok := gitGetRemote(s)
	if !ok || got.name != "upstream" {
		t.Errorf("gitGetRemote = %+v ok=%v, want upstream", got, ok)
	}
}

func TestGitGetRemoteEmpty(t *testing.T) {
	if _, ok := gitGetRemote(""); ok {
		t.Error("gitGetRemote ok = true, want false")
	}
}

func TestGitParseStatus(t *testing.T) {
	s := " A format.py\n M git.py\n M s.py\n M value.py\n?? .envrc\n"
	want := [][]string{
		{"A", "format.py"},
		{"M", "git.py"},
		{"M", "s.py"},
		{"M", "value.py"},
		{"??", ".envrc"},
	}
	got := gitParseStatus(s)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("gitParseStatus = %v, want %v", got, want)
	}
}

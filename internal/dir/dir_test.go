package dir

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

// build creates the given entries under root. Entries ending in "/" are
// directories, others are files (their parent dirs are created as needed).
func build(t *testing.T, root string, files []string) {
	t.Helper()
	for _, f := range files {
		if strings.HasSuffix(f, "/") {
			must(t, os.MkdirAll(filepath.Join(root, f), 0o755))
			continue
		}
		p := filepath.Join(root, f)
		must(t, os.MkdirAll(filepath.Dir(p), 0o755))
		must(t, os.WriteFile(p, nil, 0o644))
	}
}

func TestGitRoot(t *testing.T) {
	tests := []struct {
		name    string
		files   []string
		pwd     string
		gitroot string // "" means none
	}{
		{"git at cwd", []string{".git/"}, ".", "."},
		{"git above pwd", []string{".git/", "hoge/", "hoge/fuga"}, "hoge/", "."},
		{"nearest git deep", []string{"a/", "a/.git/", "a/a/", "a/a/a/", "a/a/a/a/", "a/a/a/a/a/"}, "a/a/a/a/a/", "a/"},
		{"nearest of multiple", []string{"a/", "a/.git/", "a/a/", "a/a/a/", "a/a/a/.git/", "a/a/a/a/", "a/a/a/a/a/"}, "a/a/a/a/a/", "a/a/a/"},
		{"git is a file", []string{"a/", "a/.git"}, "a/", ""},
		{"dotgit file at root", []string{".git"}, ".", ""},
		{"empty", []string{}, ".", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmp := t.TempDir()
			build(t, tmp, tt.files)

			got := GitRoot(filepath.Join(tmp, tt.pwd))
			if tt.gitroot == "" {
				if got != "" {
					t.Errorf("GitRoot = %q, want empty", got)
				}
				return
			}
			want, err := filepath.EvalSymlinks(filepath.Join(tmp, tt.gitroot))
			must(t, err)
			if got != want {
				t.Errorf("GitRoot = %q, want %q", got, want)
			}
		})
	}
}

func TestIsGit(t *testing.T) {
	tmp := t.TempDir()
	build(t, tmp, []string{".git/", "sub/"})
	if !IsGit(filepath.Join(tmp, "sub")) {
		t.Error("IsGit = false, want true")
	}

	tmp2 := t.TempDir()
	build(t, tmp2, []string{"sub/"})
	if IsGit(filepath.Join(tmp2, "sub")) {
		t.Error("IsGit = true, want false")
	}
}

func TestIsPython(t *testing.T) {
	markers := []string{"__init__.py", "setup.py", "setup.cfg", "pytest.ini", "manage.py"}
	for _, m := range markers {
		tmp := t.TempDir()
		build(t, tmp, []string{m, "sub/"})
		if !IsPython(filepath.Join(tmp, "sub")) {
			t.Errorf("IsPython with %s = false, want true", m)
		}
	}

	tmp := t.TempDir()
	build(t, tmp, []string{"main.go", "sub/"})
	if IsPython(filepath.Join(tmp, "sub")) {
		t.Error("IsPython = true, want false")
	}
}

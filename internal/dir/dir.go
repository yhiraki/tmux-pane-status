// Package dir detects project markers (git root, python project) by walking
// up the directory tree.
package dir

import (
	"os"
	"path/filepath"
)

// walkUp yields path and each of its ancestors up to the filesystem root.
// path is made absolute and symlink-resolved; a non-directory starts from its
// parent.
func walkUp(path string) []string {
	p, err := filepath.Abs(path)
	if err != nil {
		return nil
	}
	if resolved, err := filepath.EvalSymlinks(p); err == nil {
		p = resolved
	}
	if info, err := os.Stat(p); err != nil || !info.IsDir() {
		p = filepath.Dir(p)
	}
	var out []string
	for {
		out = append(out, p)
		parent := filepath.Dir(p)
		if parent == p {
			break
		}
		p = parent
	}
	return out
}

// GitRoot returns the nearest ancestor containing a .git directory, or "" if
// none. A .git file (not a directory) does not count.
func GitRoot(path string) string {
	for _, p := range walkUp(path) {
		if info, err := os.Stat(filepath.Join(p, ".git")); err == nil && info.IsDir() {
			return p
		}
	}
	return ""
}

// IsGit reports whether path is inside a git working tree.
func IsGit(path string) bool {
	return GitRoot(path) != ""
}

var pythonMarkers = []string{"__init__.py", "setup.py", "setup.cfg", "pytest.ini", "manage.py"}

// IsPython reports whether path is inside a python project.
func IsPython(path string) bool {
	for _, p := range walkUp(path) {
		for _, m := range pythonMarkers {
			if info, err := os.Stat(filepath.Join(p, m)); err == nil && !info.IsDir() {
				return true
			}
		}
	}
	return false
}

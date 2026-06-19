// Package render builds the tmux pane status string from collected data.
package render

import (
	"regexp"
	"sync"

	"github.com/yhiraki/tmux-pane-status/internal/collect"
	"github.com/yhiraki/tmux-pane-status/internal/config"
)

var reField = regexp.MustCompile(`\{(\w+)\}`)

// extractFields returns the set of unique field names referenced in tmpl.
func extractFields(tmpl string) map[string]bool {
	fields := map[string]bool{}
	for _, m := range reField.FindAllStringSubmatch(tmpl, -1) {
		fields[m[1]] = true
	}
	return fields
}

// Pane holds the pane-specific context resolved before rendering.
type Pane struct {
	Cwd      string // absolute, symlink-resolved
	Home     string // $HOME
	GitRoot  string // "" when not in a git repo
	IsPython bool
	ChildPID string // foreground non-shell child, "" when none
}

// IsGit reports whether the pane is inside a git working tree.
func (p Pane) IsGit() bool { return p.GitRoot != "" }

// Render produces the final pane status string.
func Render(pane Pane, col *collect.Collector, cfg *config.Config) string {
	cwd := pane.Cwd
	home := pane.Home
	isGit := pane.IsGit()
	childPID := pane.ChildPID
	// Choose base template.
	base := cfg.Formats["default"]
	if isGit {
		base = cfg.Formats["git"]
	}

	// Append command suffix when there is a foreground child and the format is
	// non-empty. Skipping this also skips the ps/pgrep calls.
	cmdFmt := cfg.Formats["command"]
	tmpl := base
	if childPID != "" && cmdFmt != "" {
		tmpl = base + cmdFmt
	}

	fields := extractFields(tmpl)

	// Determine which raw sources are needed and warm them in parallel.
	type warmup struct {
		needed bool
		fn     func()
	}
	warmups := []warmup{
		{
			needed: fields["git_remote_server"] || fields["git_repository_name"],
			fn:     func() { col.GitRemote() },
		},
		{
			needed: fields["git_status_icons"],
			fn:     func() { col.GitStatus() },
		},
		{
			needed: fields["git_current_branch"],
			fn:     func() { col.GitBranch() },
		},
		{
			needed: (fields["current_command"] || fields["current_command_elapsed"]) && childPID != "",
			fn:     func() { col.Ps(childPID, "command"); col.Ps(childPID, "etime") },
		},
	}
	var wg sync.WaitGroup
	for _, w := range warmups {
		if w.needed {
			wg.Add(1)
			fn := w.fn
			go func() { defer wg.Done(); fn() }()
		}
	}
	wg.Wait()

	gitRoot := pane.GitRoot

	// Resolve each referenced field to a display string.
	vals := map[string]string{}

	if fields["cwd"] {
		vals["cwd"] = applyStyle(cwdField(cwd, home), cfg.Styles["cwd"])
	}
	if fields["git_remote_server"] {
		v := gitRemoteServer(col.GitRemote(), cfg.Icons)
		vals["git_remote_server"] = applyStyle(v, cfg.Styles["git_remote_server"])
	}
	if fields["git_repository_name"] {
		v := gitRepositoryName(col.GitRemote())
		vals["git_repository_name"] = applyStyle(v, cfg.Styles["git_repository_name"])
	}
	if fields["git_current_branch"] {
		v := gitCurrentBranch(col.GitBranch(), cfg.Icons)
		vals["git_current_branch"] = applyStyle(v, cfg.Styles["git_current_branch"])
	}
	if fields["git_status_icons"] {
		v := gitStatusIcons(col.GitStatus())
		vals["git_status_icons"] = applyStyle(v, cfg.Styles["git_status_icons"])
	}
	if fields["git_cwd"] {
		v := gitCwd(gitRoot, cwd)
		vals["git_cwd"] = applyStyle(v, cfg.Styles["git_cwd"])
	}
	if fields["project_python"] {
		v := ""
		if pane.IsPython {
			icon := cfg.Icons["python"]
			if icon == "" {
				v = "py"
			} else {
				v = icon
			}
		}
		vals["project_python"] = applyStyle(v, cfg.Styles["project_python"])
	}
	if fields["current_command"] && childPID != "" {
		v := currentCommand(col.Ps(childPID, "command"))
		vals["current_command"] = applyStyle(v, cfg.Styles["current_command"])
	}
	if fields["current_command_elapsed"] && childPID != "" {
		v := currentCommandElapsed(col.Ps(childPID, "etime"))
		vals["current_command_elapsed"] = applyStyle(v, cfg.Styles["current_command_elapsed"])
	}

	// Replace {field} in template.
	return reField.ReplaceAllStringFunc(tmpl, func(m string) string {
		key := m[1 : len(m)-1]
		if v, ok := vals[key]; ok {
			return v
		}
		return m
	})
}

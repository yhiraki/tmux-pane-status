// Package config loads the status configuration from defaults overridden by
// PANE_STATUS_* environment variables.
package config

import (
	"os"
	"strings"
)

// Config holds the format templates, per-field tmux styles, and icons.
type Config struct {
	Formats map[string]string // "default", "command", "git"
	Styles  map[string]string // per-field tmux style, e.g. "cwd" -> "fg=blue"
	Icons   map[string]string // "python", "github", "bitbucket", "branch"
}

func defaultFormats() map[string]string {
	return map[string]string{
		"default": " {cwd} ",
		"command": " ⟩ {current_command} {current_command_elapsed}",
		"git":     " {git_remote_server} {git_repository_name}{git_cwd} {git_current_branch} {git_status_icons} {project_python} ",
	}
}

func defaultStyles() map[string]string {
	return map[string]string{
		"cwd":                     "fg=blue",
		"git_cwd":                 "",
		"git_remote_server":       "bold,fg=blue",
		"git_repository_name":     "bold,fg=blue",
		"git_current_branch":      "bold,fg=magenta",
		"git_status_icons":        "",
		"project_python":          "",
		"current_command":         "",
		"current_command_elapsed": "bold,bg=green,fg=black",
	}
}

func defaultIcons() map[string]string {
	return map[string]string{
		"python":    "🐍",
		"github":    "🐱",
		"bitbucket": "🥛",
		"branch":    "",
	}
}

// Load builds the config from defaults, then applies PANE_STATUS_* overrides
// read via getenv. A present env var overrides the default even when empty
// (this is how a section is disabled). PANE_STATUS_NO_DEFAULTS=<non-empty>
// blanks every default before overrides are applied.
func Load(getenv func(string) (string, bool)) *Config {
	c := &Config{
		Formats: defaultFormats(),
		Styles:  defaultStyles(),
		Icons:   defaultIcons(),
	}

	if v, ok := getenv("PANE_STATUS_NO_DEFAULTS"); ok && v != "" {
		blank(c.Formats)
		blank(c.Styles)
		blank(c.Icons)
	}

	apply(c.Formats, "PANE_STATUS_FORMAT_", getenv)
	apply(c.Styles, "PANE_STATUS_STYLE_", getenv)
	apply(c.Icons, "PANE_STATUS_ICON_", getenv)

	return c
}

// LoadEnv loads the config from the process environment.
func LoadEnv() *Config {
	return Load(os.LookupEnv)
}

func blank(m map[string]string) {
	for k := range m {
		m[k] = ""
	}
}

func apply(m map[string]string, prefix string, getenv func(string) (string, bool)) {
	for k := range m {
		if v, ok := getenv(prefix + strings.ToUpper(k)); ok {
			m[k] = v
		}
	}
}

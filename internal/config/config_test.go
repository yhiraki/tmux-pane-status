package config

import "testing"

// envMap builds a getenv func backed by a map.
func envMap(m map[string]string) func(string) (string, bool) {
	return func(k string) (string, bool) {
		v, ok := m[k]
		return v, ok
	}
}

func TestDefaults(t *testing.T) {
	c := Load(envMap(nil))
	if c.Formats["git"] == "" {
		t.Error("default git format is empty")
	}
	if c.Styles["cwd"] != "fg=blue" {
		t.Errorf("default cwd style = %q, want fg=blue", c.Styles["cwd"])
	}
	if c.Icons["python"] != "🐍" {
		t.Errorf("default python icon = %q", c.Icons["python"])
	}
}

func TestEnvOverride(t *testing.T) {
	c := Load(envMap(map[string]string{
		"PANE_STATUS_FORMAT_DEFAULT": " {cwd}>> ",
		"PANE_STATUS_STYLE_CWD":      "fg=red",
		"PANE_STATUS_ICON_PYTHON":    "PY",
	}))
	if c.Formats["default"] != " {cwd}>> " {
		t.Errorf("format default = %q", c.Formats["default"])
	}
	if c.Styles["cwd"] != "fg=red" {
		t.Errorf("style cwd = %q", c.Styles["cwd"])
	}
	if c.Icons["python"] != "PY" {
		t.Errorf("icon python = %q", c.Icons["python"])
	}
}

// An empty (but present) env var overrides the default — this is how the
// command suffix is disabled.
func TestEnvOverrideEmptyDisables(t *testing.T) {
	c := Load(envMap(map[string]string{
		"PANE_STATUS_FORMAT_COMMAND": "",
	}))
	if c.Formats["command"] != "" {
		t.Errorf("format command = %q, want empty", c.Formats["command"])
	}
}

func TestNoDefaults(t *testing.T) {
	c := Load(envMap(map[string]string{
		"PANE_STATUS_NO_DEFAULTS":    "1",
		"PANE_STATUS_FORMAT_DEFAULT": " {cwd} ",
	}))
	if c.Formats["git"] != "" {
		t.Errorf("git format = %q, want empty under NO_DEFAULTS", c.Formats["git"])
	}
	if c.Formats["default"] != " {cwd} " {
		t.Errorf("default format = %q, want override kept", c.Formats["default"])
	}
	if c.Styles["cwd"] != "" {
		t.Errorf("cwd style = %q, want empty under NO_DEFAULTS", c.Styles["cwd"])
	}
}

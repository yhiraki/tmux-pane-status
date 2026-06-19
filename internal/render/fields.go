package render

import (
	"fmt"
	"sort"
	"strings"
)

func gitRemoteServer(raw string, icons map[string]string) string {
	r, ok := gitGetRemote(raw)
	if !ok {
		return ""
	}
	switch r.server {
	case "github.com":
		return icons["github"]
	case "bitbucket.org":
		return icons["bitbucket"]
	}
	return r.server
}

func gitRepositoryName(raw string) string {
	r, ok := gitGetRemote(raw)
	if !ok {
		return ""
	}
	name := r.repo
	name = strings.TrimSuffix(name, ".git")
	return fmt.Sprintf("%s/%s", r.user, name)
}

func gitCurrentBranch(raw string, icons map[string]string) string {
	b := strings.TrimSpace(raw)
	if icon := icons["branch"]; icon != "" {
		return icon + " " + b
	}
	return b
}

func gitStatusIcons(raw string) string {
	entries := gitParseStatus(raw)
	if len(entries) == 0 {
		return ""
	}
	seen := map[rune]bool{}
	for _, e := range entries {
		for _, ch := range e[0] {
			seen[ch] = true
		}
	}
	runes := make([]rune, 0, len(seen))
	for ch := range seen {
		runes = append(runes, ch)
	}
	sort.Slice(runes, func(i, j int) bool { return runes[i] < runes[j] })
	return "[" + string(runes) + "]"
}

func gitCwd(root, cwd string) string {
	if !strings.HasPrefix(cwd, root) {
		return ""
	}
	rel := strings.TrimPrefix(cwd, root)
	if rel == "" {
		return ""
	}
	return rel
}

func cwdField(cwd, home string) string {
	if home != "" && strings.HasPrefix(cwd, home) {
		return "~" + cwd[len(home):]
	}
	if strings.HasPrefix(cwd, "/private") {
		return cwd[len("/private"):]
	}
	return cwd
}

func currentCommand(psOut string) string {
	lines := strings.Split(psOut, "\n")
	if len(lines) < 2 {
		return ""
	}
	cmd := strings.Fields(strings.TrimSpace(lines[1]))
	if len(cmd) == 0 {
		return ""
	}
	if cmd[0] == "ssh" {
		return parseSsh(cmd)
	}
	parts := make([]string, len(cmd))
	for i, c := range cmd {
		// basename only
		s := strings.Split(c, "/")
		parts[i] = s[len(s)-1]
	}
	return strings.Join(parts, " ")
}

func parseSsh(cmd []string) string {
	out := []string{"SSH ->"}
	i := 1
	for i < len(cmd) {
		if cmd[i] == "-p" {
			i += 2
			continue
		}
		if strings.HasPrefix(cmd[i], "-") {
			i++
			continue
		}
		out = append(out, cmd[i])
		i++
	}
	return strings.Join(out, " ")
}

func currentCommandElapsed(psOut string) string {
	lines := strings.Split(psOut, "\n")
	if len(lines) < 2 {
		return ""
	}
	return strings.TrimSpace(lines[1])
}

// applyStyle wraps s in tmux #[attr] decorators, one per comma-separated
// option in style. Returns s unchanged when style is empty.
func applyStyle(s, style string) string {
	if style == "" {
		return s
	}
	opts := strings.Split(style, ",")
	var b strings.Builder
	for _, o := range opts {
		fmt.Fprintf(&b, "#[%s]", strings.TrimSpace(o))
	}
	b.WriteString(s)
	b.WriteString("#[default]")
	return b.String()
}

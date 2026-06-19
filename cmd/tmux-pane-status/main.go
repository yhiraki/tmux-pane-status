package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/yhiraki/tmux-pane-status/internal/collect"
	"github.com/yhiraki/tmux-pane-status/internal/config"
	"github.com/yhiraki/tmux-pane-status/internal/dir"
	"github.com/yhiraki/tmux-pane-status/internal/render"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: tmux-pane-status <cwd> <pid>")
		os.Exit(1)
	}

	cwd, err := filepath.Abs(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if resolved, err := filepath.EvalSymlinks(cwd); err == nil {
		cwd = resolved
	}

	pid, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid pid: %s\n", args[1])
		os.Exit(1)
	}

	if err := os.Chdir(cwd); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cfg := config.LoadEnv()

	col := collect.New(pid)

	// Resolve child pid up front only when command format is enabled.
	childPID := ""
	if cfg.Formats["command"] != "" {
		childPID = col.ChildPID()
		// Skip if the child is a shell (zsh/bash/fish).
		if childPID != "" {
			psOut := col.Ps(childPID, "command")
			lines := strings.Split(psOut, "\n")
			if len(lines) >= 2 {
				cmd := strings.TrimSpace(lines[1])
				if isShell(cmd) {
					childPID = ""
				}
			}
		}
	}

	pane := render.Pane{
		Cwd:      cwd,
		Home:     os.Getenv("HOME"),
		GitRoot:  dir.GitRoot(cwd),
		IsPython: dir.IsPython(cwd),
		ChildPID: childPID,
	}

	fmt.Print(render.Render(pane, col, cfg))
}

func isShell(cmd string) bool {
	base := filepath.Base(cmd)
	switch base {
	case "zsh", "bash", "fish", "sh", "dash":
		return true
	}
	return false
}

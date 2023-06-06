package watcher

import (
	"context"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rjeczalik/notify"
)

func HotReload(command ...string) {
	prog := "go"
	args := []string{"run", ".", "hotreload"}
	if isStartedFromHotReloader() {
		return
	}
	if len(command) > 0 {
		prog = command[0]
		args = append(append(args, command[1:]...), "hotreload")
	}

	c := make(chan notify.EventInfo, 1)
	if err := notify.Watch("./...", c, notify.All); err != nil {
		panic(err)
	}
	defer notify.Stop(c)

	logger := log.New(os.Stdout, "HOTRELOAD: ", log.LstdFlags|log.Lmsgprefix)

	cmd := runCmd(prog, args...)
	for e := range c {
		ext := filepath.Ext(e.Path())
		switch ext {
		case ".go":
		default:
			continue
		}
		cmd.Cancel()
		println()
		logger.Println("go run .")
		cmd = runCmd(prog, args...)
	}
	os.Exit(0)
}

func runCmd(prog string, args ...string) *exec.Cmd {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, prog, args...)
	setCmdInAndOut(cmd)
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	return cmd
}

func isStartedFromHotReloader() bool {
	if len(os.Args) == 0 {
		return false
	}
	args := os.Args[1:]
	for i := range args {
		if args[i] == "hotreload" {
			return true
		}
	}
	return false
}

func setCmdInAndOut(cmd *exec.Cmd) {
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
}

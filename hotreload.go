package watcher

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/rjeczalik/notify"
)

// Call this function as the first function in your code.
// It recompiles your program after a change in a go file and starts it after stopping the previous application.
//
//	func main() {
//		watcher.HotReload()
//		// your code
//	}
func HotReload(command ...string) {
	prog := "go"
	args := []string{"run", "."}
	if isStartedFromHotReloader() {
		return
	}
	if len(command) > 0 {
		prog = command[0]
		args = command[1:]
	}
	args = append(args, "hotreload")

	c := make(chan notify.EventInfo, 1)
	if err := notify.Watch("./...", c, notify.All); err != nil {
		panic(err)
	}
	defer notify.Stop(c)

	logger := log.New(os.Stdout, "HOTRELOAD: ", log.LstdFlags|log.Lmsgprefix)

	controlC := catchControlC()
	cmd := runCmd(prog, args...)

	for {
		select {
		case e := <-c:
			ext := filepath.Ext(e.Path())
			switch ext {
			case ".go":
			default:
				continue
			}
			syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			cmd.Process.Wait()
			println()
			logger.Println("go run .")
			cmd = runCmd(prog, args...)
		case <-controlC:
			logger.Println("STOPPED")
			if cmd != nil {
				syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
				cmd.Process.Wait()
			}
			os.Exit(0)
		}
	}
}

func runCmd(prog string, args ...string) *exec.Cmd {
	cmd := exec.Command(prog, args...)
	setCmdProps(cmd)
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

func setCmdProps(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
}

func catchControlC() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	return c
}

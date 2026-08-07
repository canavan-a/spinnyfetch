package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type terminalModule struct{}

func (terminalModule) Key() string { return "Terminal" }

func (terminalModule) Fetch() (string, error) {
	if name := os.Getenv("TERM_PROGRAM"); name != "" {
		return name, nil
	}
	if name, err := parentProcessName(); err == nil && name != "" {
		return name, nil
	}
	if term := os.Getenv("TERM"); term != "" {
		return term, nil
	}
	return "", fmt.Errorf("no terminal detected")
}

// parentProcessName walks up from the current process to find the
// terminal emulator's process name via /proc, skipping the immediate
// shell parent.
func parentProcessName() (string, error) {
	pid := os.Getppid()
	for i := 0; i < 5 && pid > 1; i++ {
		comm, err := readFirstLine(fmt.Sprintf("/proc/%d/comm", pid))
		if err != nil {
			return "", err
		}
		if !isShell(comm) {
			return comm, nil
		}
		ppid, err := parentOf(pid)
		if err != nil {
			return "", err
		}
		pid = ppid
	}
	return "", fmt.Errorf("terminal process not found")
}

func isShell(name string) bool {
	switch filepath.Base(name) {
	case "bash", "zsh", "fish", "sh", "dash":
		return true
	}
	return false
}

func parentOf(pid int) (int, error) {
	kv, err := readKV(fmt.Sprintf("/proc/%d/status", pid), ":")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(kv["PPid"])
}

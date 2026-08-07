package modules

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type shellModule struct{}

func (shellModule) Key() string { return "Shell" }

func (shellModule) Fetch() (string, error) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "", fmt.Errorf("$SHELL not set")
	}
	name := filepath.Base(shell)

	out, err := exec.Command(shell, "--version").Output()
	if err != nil {
		return name, nil
	}
	version := firstVersionToken(string(out))
	if version == "" {
		return name, nil
	}
	return fmt.Sprintf("%s %s", name, version), nil
}

// firstVersionToken pulls the first dotted-number token out of a
// "--version" banner (e.g. "5.3.9" out of "GNU bash, version 5.3.9(1)...").
func firstVersionToken(s string) string {
	for _, tok := range strings.Fields(s) {
		tok = strings.Trim(tok, ",()")
		if tok == "" {
			continue
		}
		if tok[0] < '0' || tok[0] > '9' {
			continue
		}
		return tok
	}
	return ""
}

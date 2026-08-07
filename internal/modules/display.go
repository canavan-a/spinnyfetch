package modules

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type displayModule struct{}

func (displayModule) Key() string { return "Display" }

func (displayModule) Fetch() (string, error) {
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		if s, err := wlrRandrDisplay(); err == nil {
			return s, nil
		}
	}
	if os.Getenv("DISPLAY") != "" {
		if s, err := xrandrDisplay(); err == nil {
			return s, nil
		}
	}
	return "", fmt.Errorf("no display detected")
}

var wlrModeRe = regexp.MustCompile(`(\d+)x(\d+)[^@]*@\s*([\d.]+)\s*Hz`)

func wlrRandrDisplay() (string, error) {
	out, err := exec.Command("wlr-randr").Output()
	if err != nil {
		return "", err
	}
	m := wlrModeRe.FindStringSubmatch(string(out))
	if m == nil {
		return "", fmt.Errorf("no mode found in wlr-randr output")
	}
	return fmt.Sprintf("%sx%s, %s Hz", m[1], m[2], m[3]), nil
}

var xrandrModeRe = regexp.MustCompile(`(\d+)x(\d+)\s+([\d.]+)\*`)

func xrandrDisplay() (string, error) {
	out, err := exec.Command("xrandr").Output()
	if err != nil {
		return "", err
	}
	m := xrandrModeRe.FindStringSubmatch(string(out))
	if m == nil {
		return "", fmt.Errorf("no active mode found in xrandr output")
	}
	return fmt.Sprintf("%sx%s, %s Hz", m[1], m[2], strings.TrimSuffix(m[3], ".")), nil
}

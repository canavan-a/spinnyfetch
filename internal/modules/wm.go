package modules

import (
	"fmt"
	"os"
)

type wmModule struct{}

func (wmModule) Key() string { return "WM" }

func (wmModule) Fetch() (string, error) {
	name := os.Getenv("XDG_CURRENT_DESKTOP")
	if name == "" {
		name = os.Getenv("DESKTOP_SESSION")
	}
	if name == "" {
		if os.Getenv("HYPRLAND_INSTANCE_SIGNATURE") != "" {
			name = "Hyprland"
		} else if os.Getenv("SWAYSOCK") != "" {
			name = "sway"
		}
	}
	if name == "" {
		return "", fmt.Errorf("no WM/DE detected")
	}

	session := "X11"
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		session = "Wayland"
	}
	return fmt.Sprintf("%s (%s)", name, session), nil
}

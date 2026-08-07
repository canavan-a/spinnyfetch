package modules

import (
	"fmt"
	"path/filepath"
	"strings"
)

type batteryModule struct{}

func (batteryModule) Key() string { return "Battery" }

func (batteryModule) Fetch() (string, error) {
	matches, err := filepath.Glob("/sys/class/power_supply/BAT*")
	if err != nil || len(matches) == 0 {
		return "", fmt.Errorf("no battery found")
	}
	bat := matches[0]

	capacity, err := readFirstLine(filepath.Join(bat, "capacity"))
	if err != nil {
		return "", err
	}
	status, _ := readFirstLine(filepath.Join(bat, "status"))

	name := strings.TrimPrefix(filepath.Base(bat), "BAT")
	label := fmt.Sprintf("BAT%s", name)

	acState := "Battery"
	switch status {
	case "Charging", "Full":
		acState = "AC Connected"
	case "Discharging":
		acState = "Discharging"
	}
	return fmt.Sprintf("%s%% [%s] (%s)", capacity, acState, label), nil
}

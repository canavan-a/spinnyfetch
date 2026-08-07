package modules

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type uptimeModule struct{}

func (uptimeModule) Key() string { return "Uptime" }

func (uptimeModule) Fetch() (string, error) {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return "", err
	}
	fields := strings.Fields(string(data))
	if len(fields) < 1 {
		return "", fmt.Errorf("bad /proc/uptime")
	}
	secondsF, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return "", err
	}
	total := int(secondsF)
	days := total / 86400
	hours := (total % 86400) / 3600
	mins := (total % 3600) / 60

	var parts []string
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d day%s", days, plural(days)))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d hour%s", hours, plural(hours)))
	}
	parts = append(parts, fmt.Sprintf("%d min%s", mins, plural(mins)))
	return strings.Join(parts, ", "), nil
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

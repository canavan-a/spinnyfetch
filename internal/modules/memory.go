package modules

import (
	"fmt"
	"strconv"
	"strings"
)

type memoryModule struct{}

func (memoryModule) Key() string { return "Memory" }

func (memoryModule) Fetch() (string, error) {
	kv, err := readKV("/proc/meminfo", ":")
	if err != nil {
		return "", err
	}
	totalKB, err := meminfoKB(kv["MemTotal"])
	if err != nil {
		return "", err
	}
	availKB, err := meminfoKB(kv["MemAvailable"])
	if err != nil {
		return "", err
	}
	usedKB := totalKB - availKB
	pct := 0
	if totalKB > 0 {
		pct = int(float64(usedKB) / float64(totalKB) * 100)
	}
	return fmt.Sprintf("%s / %s (%d%%)", formatGiB(usedKB), formatGiB(totalKB), pct), nil
}

func meminfoKB(field string) (float64, error) {
	f := strings.Fields(field)
	if len(f) == 0 {
		return 0, fmt.Errorf("empty meminfo field")
	}
	return strconv.ParseFloat(f[0], 64)
}

func formatGiB(kb float64) string {
	return fmt.Sprintf("%.2f GiB", kb/1024/1024)
}

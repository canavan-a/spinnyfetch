package modules

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type swapModule struct{}

func (swapModule) Key() string { return "Swap" }

func (swapModule) Fetch() (string, error) {
	f, err := os.Open("/proc/swaps")
	if err != nil {
		return "", err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	sc.Scan() // header
	var totalKB, usedKB float64
	found := false
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) < 4 {
			continue
		}
		found = true
		t, _ := strconv.ParseFloat(fields[2], 64)
		u, _ := strconv.ParseFloat(fields[3], 64)
		totalKB += t
		usedKB += u
	}
	if !found || totalKB == 0 {
		return "Disabled", nil
	}
	return fmt.Sprintf("%s / %s", formatGiB(usedKB), formatGiB(totalKB)), nil
}

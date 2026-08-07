package modules

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type cpuModule struct{}

func (cpuModule) Key() string { return "CPU" }

func (cpuModule) Fetch() (string, error) {
	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return "", err
	}
	defer f.Close()

	var model string
	cores := 0
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "model name") {
			if model == "" {
				if i := strings.Index(line, ":"); i >= 0 {
					model = strings.TrimSpace(line[i+1:])
				}
			}
			cores++
		}
	}
	if model == "" {
		return "", fmt.Errorf("no cpu model found")
	}

	freq := currentFreqGHz()
	if freq != "" {
		return fmt.Sprintf("%s (%d) @ %s GHz", model, cores, freq), nil
	}
	return fmt.Sprintf("%s (%d)", model, cores), nil
}

func currentFreqGHz() string {
	data, err := os.ReadFile("/sys/devices/system/cpu/cpu0/cpufreq/scaling_cur_freq")
	if err != nil {
		return ""
	}
	khz, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
	if err != nil {
		return ""
	}
	return strconv.FormatFloat(khz/1e6, 'f', 2, 64)
}

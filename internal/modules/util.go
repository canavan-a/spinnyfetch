package modules

import (
	"bufio"
	"os"
	"strings"
)

// readKV reads a "KEY=value" or "key: value" style file and returns a map
// of trimmed, unquoted values. Used for /etc/os-release and /proc/meminfo.
func readKV(path string, sep string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	out := map[string]string{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		i := strings.Index(line, sep)
		if i < 0 {
			continue
		}
		key := strings.TrimSpace(line[:i])
		val := strings.TrimSpace(line[i+len(sep):])
		val = strings.Trim(val, `"`)
		out[key] = val
	}
	return out, sc.Err()
}

func readFirstLine(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	if sc.Scan() {
		return strings.TrimSpace(sc.Text()), nil
	}
	return "", sc.Err()
}

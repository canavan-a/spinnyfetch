package modules

import (
	"fmt"
	"path/filepath"
	"strings"
)

type gpuModule struct{}

func (gpuModule) Key() string { return "GPU" }

var pciVendorNames = map[string]string{
	"0x1002": "AMD",
	"0x10de": "NVIDIA",
	"0x8086": "Intel",
}

// Fetch reports all GPUs joined into one value; each is also independently
// tagged Discrete/Integrated by boot_vga (the primary/integrated adapter
// usually has boot_vga=1).
func (gpuModule) Fetch() (string, error) {
	cards, err := filepath.Glob("/sys/class/drm/card[0-9]*")
	if err != nil {
		return "", err
	}

	var lines []string
	seen := map[string]bool{}
	for _, card := range cards {
		devPath := filepath.Join(card, "device")
		if seen[devPath] {
			continue
		}
		seen[devPath] = true

		vendorID, err := readFirstLine(filepath.Join(devPath, "vendor"))
		if err != nil {
			continue
		}
		vendor := pciVendorNames[strings.ToLower(vendorID)]
		if vendor == "" {
			vendor = vendorID
		}

		kind := "Integrated"
		if bootVGA, err := readFirstLine(filepath.Join(devPath, "boot_vga")); err == nil && bootVGA != "1" {
			kind = "Discrete"
		}

		lines = append(lines, fmt.Sprintf("%s Graphics [%s]", vendor, kind))
	}

	if len(lines) == 0 {
		return "", fmt.Errorf("no GPU found")
	}
	return strings.Join(lines, ", "), nil
}

package modules

import (
	"fmt"
	"strings"
)

type hostModule struct{}

func (hostModule) Key() string { return "Host" }

func (hostModule) Fetch() (string, error) {
	name, err := readFirstLine("/sys/class/dmi/id/product_name")
	if err != nil || name == "" {
		return "", fmt.Errorf("no product name")
	}
	version, _ := readFirstLine("/sys/class/dmi/id/product_version")
	version = strings.TrimSpace(version)
	if version == "" || version == "None" {
		return name, nil
	}
	return fmt.Sprintf("%s (%s)", name, version), nil
}

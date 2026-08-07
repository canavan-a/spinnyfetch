package modules

import "fmt"

type osModule struct{}

func (osModule) Key() string { return "OS" }

func (osModule) Fetch() (string, error) {
	kv, err := readKV("/etc/os-release", "=")
	if err != nil {
		return "", err
	}
	name := kv["PRETTY_NAME"]
	if name == "" {
		name = kv["NAME"]
	}
	if name == "" {
		return "", fmt.Errorf("no os-release name")
	}
	arch, err := machine()
	if err != nil {
		return name, nil
	}
	return fmt.Sprintf("%s %s", name, arch), nil
}

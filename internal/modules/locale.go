package modules

import (
	"fmt"
	"os"
)

type localeModule struct{}

func (localeModule) Key() string { return "Locale" }

func (localeModule) Fetch() (string, error) {
	if l := os.Getenv("LANG"); l != "" {
		return l, nil
	}
	if l := os.Getenv("LC_ALL"); l != "" {
		return l, nil
	}
	return "", fmt.Errorf("no locale set")
}

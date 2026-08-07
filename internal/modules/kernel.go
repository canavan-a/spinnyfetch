package modules

import (
	"fmt"

	"golang.org/x/sys/unix"
)

type kernelModule struct{}

func (kernelModule) Key() string { return "Kernel" }

func (kernelModule) Fetch() (string, error) {
	var u unix.Utsname
	if err := unix.Uname(&u); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s %s", cstr(u.Sysname), cstr(u.Release)), nil
}

// machine returns the machine architecture string (e.g. "x86_64").
func machine() (string, error) {
	var u unix.Utsname
	if err := unix.Uname(&u); err != nil {
		return "", err
	}
	return cstr(u.Machine), nil
}

func cstr(b [65]byte) string {
	n := 0
	for n < len(b) && b[n] != 0 {
		n++
	}
	return string(b[:n])
}

package modules

import (
	"fmt"

	"golang.org/x/sys/unix"
)

type diskModule struct{}

func (diskModule) Key() string { return "Disk (/)" }

func (diskModule) Fetch() (string, error) {
	var st unix.Statfs_t
	if err := unix.Statfs("/", &st); err != nil {
		return "", err
	}
	total := float64(st.Blocks) * float64(st.Bsize)
	free := float64(st.Bfree) * float64(st.Bsize)
	used := total - free
	pct := 0
	if total > 0 {
		pct = int(used / total * 100)
	}
	return fmt.Sprintf("%s / %s (%d%%) - %s", formatGiBBytes(used), formatGiBBytes(total), pct, fsType(st.Type)), nil
}

func formatGiBBytes(b float64) string {
	const gib = 1024 * 1024 * 1024
	if b >= gib*1024 {
		return fmt.Sprintf("%.2f TiB", b/gib/1024)
	}
	return fmt.Sprintf("%.2f GiB", b/gib)
}

// fsType maps common Linux magic numbers (statfs.f_type) to names; falls
// back to the raw hex value for anything unrecognized.
func fsType(magic int64) string {
	switch magic {
	case 0xEF53:
		return "ext4"
	case 0x9123683E:
		return "btrfs"
	case 0x58465342:
		return "xfs"
	case 0x2011BAB0:
		return "exfat"
	case 0x01021994:
		return "tmpfs"
	default:
		return fmt.Sprintf("0x%X", magic)
	}
}

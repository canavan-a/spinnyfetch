package modules

// Module fetches a single line of system info, keyed like fastfetch's
// "Key: value" output rows. Fetch must fail soft: return an error to have
// the caller omit the line rather than crash the whole run.
type Module interface {
	Key() string
	Fetch() (string, error)
}

// registry maps config module names to constructors.
var registry = map[string]Module{
	"os":       osModule{},
	"host":     hostModule{},
	"kernel":   kernelModule{},
	"uptime":   uptimeModule{},
	"packages": packagesModule{},
	"shell":    shellModule{},
	"display":  displayModule{},
	"wm":       wmModule{},
	"terminal": terminalModule{},
	"cpu":      cpuModule{},
	"gpu":      gpuModule{},
	"memory":   memoryModule{},
	"swap":     swapModule{},
	"disk":     diskModule{},
	"local_ip": localIPModule{},
	"battery":  batteryModule{},
	"locale":   localeModule{},
}

// Resolve returns the modules for the given ordered list of names, skipping
// any name that doesn't match a known module.
func Resolve(names []string) []Module {
	out := make([]Module, 0, len(names))
	for _, n := range names {
		if m, ok := registry[n]; ok {
			out = append(out, m)
		}
	}
	return out
}

// Line is a resolved "Key: value" row ready for rendering.
type Line struct {
	Key   string
	Value string
}

// FetchAll runs each module and returns only the ones that succeeded, in
// the given order.
func FetchAll(mods []Module) []Line {
	out := make([]Line, 0, len(mods))
	for _, m := range mods {
		v, err := m.Fetch()
		if err != nil || v == "" {
			continue
		}
		out = append(out, Line{Key: m.Key(), Value: v})
	}
	return out
}

package config

import (
	"flag"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type SpinConfig struct {
	Enabled         bool    `yaml:"enabled"`
	SpeedDegPerTick float64 `yaml:"speed_deg_per_tick"`
	FPS             int     `yaml:"fps"`
	Axis            string  `yaml:"axis"`
}

type LogoConfig struct {
	Colors    []string   `yaml:"colors"`
	ColorMode string     `yaml:"color_mode"`
	Spin      SpinConfig `yaml:"spin"`
}

type OutputConfig struct {
	Gap       int    `yaml:"gap"`
	KeyColor  string `yaml:"key_color"`
	Separator string `yaml:"separator"`
}

type Config struct {
	Modules []string     `yaml:"modules"`
	Logo    LogoConfig   `yaml:"logo"`
	Output  OutputConfig `yaml:"output"`
}

func Default() Config {
	return Config{
		Modules: []string{
			"os", "host", "kernel", "uptime", "packages", "shell",
			"display", "wm", "terminal", "cpu", "gpu", "memory", "swap",
			"disk", "local_ip", "battery", "locale",
		},
		Logo: LogoConfig{
			Colors:    []string{"#7ebae4", "#4a90d9"},
			ColorMode: "truecolor",
			Spin: SpinConfig{
				Enabled:         true,
				SpeedDegPerTick: 4,
				FPS:             30,
				Axis:            "y",
			},
		},
		Output: OutputConfig{
			Gap:       4,
			KeyColor:  "#89b4fa",
			Separator: ":",
		},
	}
}

// DefaultPath returns ~/.config/spinnyfetch/config.yaml.
func DefaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "spinnyfetch", "config.yaml")
}

// Load reads and merges a YAML config file over the defaults. A missing
// file is not an error — the defaults are returned as-is.
func Load(path string) (Config, error) {
	cfg := Default()
	if path == "" {
		return cfg, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// Flags holds CLI overrides parsed on top of a loaded Config.
type Flags struct {
	NoSpin     bool
	Once       bool
	ConfigPath string
	FPS        int
	Speed      float64
}

// ParseFlags parses os.Args[1:] into Flags. FPS/Speed default to 0, meaning
// "not set" — callers should only apply them over the config when nonzero.
func ParseFlags(args []string) Flags {
	fs := flag.NewFlagSet("spinnyfetch", flag.ExitOnError)
	f := Flags{}
	fs.BoolVar(&f.NoSpin, "no-spin", false, "disable the logo spin animation")
	fs.BoolVar(&f.Once, "once", false, "print a single static frame and exit")
	fs.StringVar(&f.ConfigPath, "config", DefaultPath(), "path to config.yaml")
	fs.IntVar(&f.FPS, "fps", 0, "override spin frames per second")
	fs.Float64Var(&f.Speed, "speed", 0, "override spin speed in degrees/tick")
	fs.Parse(args)
	return f
}

// Apply merges CLI flag overrides onto a loaded Config.
func Apply(cfg Config, f Flags) Config {
	if f.NoSpin {
		cfg.Logo.Spin.Enabled = false
	}
	if f.FPS != 0 {
		cfg.Logo.Spin.FPS = f.FPS
	}
	if f.Speed != 0 {
		cfg.Logo.Spin.SpeedDegPerTick = f.Speed
	}
	return cfg
}

package main

import (
	"fmt"
	"os"

	"github.com/canavan-a/spinnyfetch/internal/config"
	"github.com/canavan-a/spinnyfetch/internal/logo"
	"github.com/canavan-a/spinnyfetch/internal/modules"
	"github.com/canavan-a/spinnyfetch/internal/render"
)

func main() {
	flags := config.ParseFlags(os.Args[1:])

	cfg, err := config.Load(flags.ConfigPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "spinnyfetch: failed to load config: %v\n", err)
		os.Exit(1)
	}
	cfg = config.Apply(cfg, flags)

	colorA, colorB := logo.ColorA, logo.ColorB
	if len(cfg.Logo.Colors) >= 2 {
		colorA, colorB = cfg.Logo.Colors[0], cfg.Logo.Colors[1]
	}
	grid := logo.ParseWithColors(logo.Source, colorA, colorB)
	mods := modules.Resolve(cfg.Modules)
	lines := modules.FetchAll(mods)

	opts := render.Options{
		Gap:       cfg.Output.Gap,
		KeyColor:  cfg.Output.KeyColor,
		Separator: cfg.Output.Separator,
		ColorMode: render.ParseColorMode(cfg.Logo.ColorMode),
		Spin: render.SpinOptions{
			Enabled:         cfg.Logo.Spin.Enabled && !flags.Once,
			SpeedDegPerTick: cfg.Logo.Spin.SpeedDegPerTick,
			FPS:             cfg.Logo.Spin.FPS,
		},
	}
	render.Run(grid, lines, opts)
}

# spinnyfetch

```
nix run .#default -- --once
```

## Usage

```
spinnyfetch [flags]

  --no-spin           disable the logo spin animation
  --once              print a single static frame and exit
  --config <path>     path to config.yaml (default ~/.config/spinnyfetch/config.yaml)
  --fps <n>           override spin frames per second
  --speed <deg>       override spin speed in degrees/tick
```

By default the logo spins continuously until interrupted (Ctrl-C).

## Config

`~/.config/spinnyfetch/config.yaml`:

```yaml
modules: [os, host, kernel, uptime, packages, shell, display, wm,
          terminal, cpu, gpu, memory, swap, disk, local_ip, battery, locale]

logo:
  colors: ["#7ebae4", "#4a90d9"]
  color_mode: truecolor   # truecolor | 256
  spin:
    enabled: true
    speed_deg_per_tick: 4
    fps: 30

output:
  gap: 4
  key_color: "#89b4fa"
  separator: ":"
```

Any subset of these keys may be set; unset keys fall back to defaults.

## Building

```
go build .
```

Or via Nix:

```
nix build .#default
```

## Development

```
go test ./...
go run ./cmd/spinpreview   # visual check of the spin transform alone
```

## License

MIT, see [LICENSE](LICENSE).

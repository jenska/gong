# Gong! — The Go Pong

A retro-styled [Pong](https://en.wikipedia.org/wiki/Pong) clone written in Go
with [Ebitengine](https://ebitengine.org/).

Features include:

- One-player, two-player, and AI-versus-AI modes
- Human-like computer players with reaction delay, imperfect prediction,
  gradual acceleration, braking, and occasional mistakes
- Fullscreen mode, sound effects, and volume controls
- Cached HUD rendering and allocation-free sprite trail updates

## Run from source

Install Go **1.26 or later**, then run:

```bash
git clone https://github.com/jenska/gong.git
cd gong
go run .
```

![Screenshot](game/assets/screenshot.png)

## Install

To build and install a standalone `gong` binary into your Go bin path:

```bash
go install
```

## Makefile targets

Common development commands:

```bash
make run
make build
make test
make fmt
make tidy
```

## Controls

### Main menu

- `V`: two players
- `A`: computer versus player
- `B`: computer versus computer
- `H`: show controls
- `F`: toggle fullscreen
- `W` / `S`: increase or decrease volume
- `Esc`: quit

### During play

- `W` / `S`: move the left paddle
- `↑` / `↓`: move the right paddle
- `Space`: pause or resume
- `R`: restart while paused

The first player to 10 points wins.

## Development

Run the complete validation suite:

```bash
gofmt -w main.go game/*.go
go test -race ./...
go vet ./...
go build ./...
go fix -diff ./...
```

Run the rendering microbenchmarks:

```bash
go test -run '^$' -bench 'Benchmark(HUD|Sprite)' -benchmem ./game
```

## Dependencies

Dependencies are managed with Go modules. Ebitengine may require additional
native development libraries on Linux.

### Debian / Ubuntu

```bash
sudo apt install libc6-dev libglu1-mesa-dev libgl1-mesa-dev libxcursor-dev libxi-dev libxinerama-dev libxrandr-dev libxxf86vm-dev libasound2-dev pkg-config
```

### Fedora

```bash
sudo dnf install mesa-libGLU-devel mesa-libGLES-devel libXrandr-devel libXcursor-devel libXinerama-devel libXi-devel libXxf86vm-devel alsa-lib-devel pkg-config
```

## Thanks

- [OpenGameArt.org](https://opengameart.org/)

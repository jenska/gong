# Gong! — The Go Pong

A retro-styled [Pong](https://en.wikipedia.org/wiki/Pong) clone written in Go
with [Ebitengine](https://ebitengine.org/).

[Play Gong in your browser](https://jenska.github.io/gong/) ·
[MIT license](LICENSE)

Features include:

- One-player, two-player, and AI-versus-AI modes
- Beginner, human-like, and perfect-prediction computer players
- An exported controller API for custom player and AI implementations
- Contact-based paddle deflection, movement-driven spin, directional serves,
  and capped speed progression
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
make web
```

## Custom controllers

Paddle input is defined by the exported `game.Controller` interface. A
controller receives an immutable view of the current game and returns its
desired movement:

```go
type Controller interface {
    Name() string
    Control(game.GameView) game.Control
}
```

Built-in implementations are created with:

```go
game.NewKeyboardController(upKey, downKey)
game.NewBeginnerAI()
game.NewHumanLikeAI()
game.NewPerfectAI()
```

This keeps gameplay, input, and AI strategy separate, making new controllers
straightforward to implement and test.

Custom controllers can be used directly:

```go
gong := game.NewGong()
gong.StartMatch(myLeftController, game.NewPerfectAI())
ebiten.RunGame(gong)
```

## Development

Run the complete validation suite:

```bash
gofmt -w main.go game/*.go
go test -race ./...
go vet ./...
go build ./...
go fix -diff ./...
```
Build and serve the WebAssembly version locally:

```bash
make serve-web
```

Then open <http://localhost:8080>.

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

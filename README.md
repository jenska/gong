# Gong! - The Go Pong

Retro styled classic [Pong](https://en.wikipedia.org/wiki/Pong) clone based upon [Ebiten](https://ebiten.org/) game framework.

With go version 16 or above installed, type

```bash
git clone http://github.com/jenska/gong
cd gong
go run main.go
```

to start

![Screenshot](game/assets/screenshot.png)

For a standalone gong game, type

```bash
go install
```


## Dependencies

Ebiten requires some additional libs installed

### Debian / Ubuntu

```bash
sudo apt install libc6-dev libglu1-mesa-dev libgl1-mesa-dev libxcursor-dev libxi-dev libxinerama-dev libxrandr-dev libxxf86vm-dev libasound2-dev pkg-config
```

### Fedora

```bash
sudo dnf install mesa-libGLU-devel mesa-libGLES-devel libXrandr-devel libXcursor-devel libXinerama-devel libXi-devel libXxf86vm-devel alsa-lib-devel pkg-config
```

## Thanks to

[OpenGameArt.Org](https://opengameart.org/)

## Soon to come

- ~~update package structure~~
- ~~stereo sounds~~
- ~~make window resizable~~
- ~~upgrade to ebiten 2.0~~
- ~~embedded assets~~
- AI support (beat the computer opponent, computer vs computer becomes more interessting)
- WebAssembly support

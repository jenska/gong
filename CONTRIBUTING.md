# Contributing to Gong

Contributions are welcome. Small, focused pull requests are the easiest to
review and merge.

## Getting started

1. Install Go 1.26 or later and the native libraries listed in the README.
2. Fork and clone the repository.
3. Create a branch from `master`.
4. Make and test your change.
5. Open a pull request explaining the behavior and design choices.

Run the validation suite before submitting:

```bash
make fmt
go test -race ./...
go vet ./...
go build ./...
go fix -diff ./...
make web
```

## Project guidelines

- Keep game simulation in `Update`; rendering in `Draw` should not mutate game
  state.
- Avoid allocations in per-frame update and rendering paths where practical.
- Add deterministic tests for gameplay and AI changes.
- Implement new player strategies through the exported `game.Controller`
  interface instead of adding controller-specific logic to `paddle`.
- Add benchmarks when changing performance-sensitive rendering code.
- Keep new gameplay mechanics optional when they significantly change classic
  Pong behavior.
- Do not add generated files from `dist/` to Git.

## Issues

Search existing issues before opening a new one. Bug reports should include the
operating system, Go version, steps to reproduce, and any relevant logs.

For larger features, open an issue first so the design can be discussed before
implementation begins.

## Pull requests

Pull requests should:

- Address one coherent change.
- Include tests or explain why testing is not practical.
- Update documentation for user-visible behavior.
- Pass all CI checks.

By participating, you agree to follow the project's
[Code of Conduct](CODE_OF_CONDUCT.md).

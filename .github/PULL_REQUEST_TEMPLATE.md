## Summary

Describe the change and why it is needed.

## Verification

- [ ] `make fmt`
- [ ] `go test -race ./...`
- [ ] `go vet ./...`
- [ ] `go build ./...`
- [ ] `make web`

## Checklist

- [ ] Tests cover new or changed behavior.
- [ ] User-visible changes are documented.
- [ ] Per-frame code avoids unnecessary allocations.
- [ ] The change is focused and contains no unrelated edits.

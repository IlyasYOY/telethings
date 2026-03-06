# Contributing to Telethings

Thank you for your interest in contributing!

## Prerequisites

- Go 1.25+
- macOS (required to run the bot; tests can run on any OS)
- Things 3 app (to test end-to-end)

## Development Setup

```bash
git clone https://github.com/IlyasYOY/telethings.git
cd telethings
```

## Building

```bash
make build       # builds to ./bin/telethings
make run         # build + run (requires env vars)
```

## Testing

```bash
make test        # run all tests
make vet         # run go vet
```

Run a single test:
```bash
go test ./internal/bot/... -run TestHandler_HandleAdd_ValidCommand
```

## Making Changes

1. Fork the repository
2. Create a feature branch: `git checkout -b my-fix`
3. Make your changes
4. Ensure tests pass: `make test && make vet`
5. Commit and push your branch
6. Open a pull request

## Code Style

- Follow standard Go conventions (`gofmt`)
- Use interfaces for testability — avoid concrete dependencies in handlers
- Use `minimock` (via `//go:generate`) for mocks; do not write hand-crafted fakes
- External test packages (`package foo_test`) are preferred for black-box testing

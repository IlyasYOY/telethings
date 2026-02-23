# Instructions

## Build, Test, and Lint

```bash
make build        # builds to ./bin/telethings
make test         # go test ./...
make vet          # go vet ./...
make run          # build + run
```

Run a single test:
```bash
go test ./internal/bot/... -run TestHandler_HandleAdd_ValidCommand
go test ./internal/thingsurl/... -run TestAdd_Title
```

Required environment variables to run the bot:
- `TELETHINGS_TELEGRAM_TOKEN` — Telegram bot token
- `TELETHINGS_THINGS_AUTH_TOKEN` — Things 3 URL scheme auth token

## Architecture

The bot is a single binary (`cmd/telethings/main.go`) that long-polls Telegram for updates and opens Things 3 URLs locally on macOS.

**Request flow:**
1. `bot.Bot.Run` receives a Telegram update via long-polling
2. `bot.Handler.Handle` dispatches on the command name (`/start`, `/add`, `/today`, `/inbox`)
3. For `/add`, `parseAddCommand` in `add_parser.go` converts the message text into a `thingsurl` URL string
4. `opener.MacOSOpener` invokes `open <url>` — this is macOS-only and triggers the Things 3 app directly

**Package responsibilities:**
- `internal/bot` — Telegram update handling, command parsing, `MessageSender` interface
- `internal/thingsreader` — reads Things 3 tasks via AppleScript (including paged list reads)
- `internal/thingsurl` — fluent builder for Things 3 URL scheme (`things:///add?...`)
- `internal/config` — reads env vars; returns typed errors for each missing variable
- `internal/opener` — `MacOSOpener` (production) and `openertest.RecordingOpener` (tests)

## Key Conventions

**Interfaces for testability:** `bot.Handler` depends on `MessageSender` and the unexported `opener` interface, not on concrete types. Tests use `fakeSender` (inline in `handler_test.go`) and `openertest.RecordingOpener` to avoid network/OS calls.

**`openertest` package:** Test helpers for the `opener` interface live in `internal/opener/openertest/`. When writing tests that need an opener, use `openertest.RecordingOpener` and inspect `.URLs` to verify the correct Things3 URL was constructed.

**External test packages:** Tests use `package bot_test` (not `package bot`), so only exported identifiers are accessible. Keep this pattern when adding new tests.

**`thingsurl` builder pattern:** All URL construction goes through the fluent builder:
```go
thingsurl.New(authToken).Add(title).WithWhen("today").WithTags("work").String()
```
URL encoding uses `%20` for spaces (not `+`) — see `encodeParams` in `types.go`.

**Adding a new command:** Add a case to the `switch` in `handler.go`, register it in `bot.New` via `tgbotapi.NewSetMyCommands`, and update `/start` response text to include it.

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
go test ./internal/... -run TestAdd_Title
```

Required environment variables to run the bot:
- `TELETHINGS_TELEGRAM_TOKEN` — Telegram bot token
- `TELETHINGS_ALLOWED_USER_IDS` — comma-separated list of allowed Telegram user IDs

## Architecture

The bot is a single binary (`cmd/telethings/main.go`) that long-polls Telegram for updates and interacts with Things 3 via AppleScript on macOS.

**Request flow:**
1. `bot.Bot.Run` receives a Telegram update via long-polling
2. `bot.Handler.Handle` dispatches on the command name (`/start`, `/add`, `/today`, `/inbox`)
3. For `/add`, `parseAddCommand` in `add_parser.go` parses the message text into a structured input
4. `thingser.AppleScriptReader` executes AppleScript to add the task to Things 3 directly

**Package responsibilities:**
- `internal/bot` — Telegram update handling, command parsing, `MessageSender` interface
- `internal/thingser` — reads and writes Things 3 tasks via AppleScript
- `internal/config` — reads env vars; returns typed errors for each missing variable
- `internal/db` — SQLite-backed task store (used for deferred/tracked tasks); defaults to `$XDG_DATA_HOME/telethings/telethings.db` when `TELETHINGS_DB_DSN` is unset

## Key Conventions

**Interfaces for testability:** `bot.Handler` depends on `MessageSender` and the unexported `thingser`/`taskStore` interfaces, not on concrete types.

**External test packages:** Tests use `package bot_test` (not `package bot`), so only exported identifiers are accessible. Keep this pattern when adding new tests.

**Adding a new command:** Add a case to the `switch` in `handler.go`, register it in `bot.New` via `tgbotapi.NewSetMyCommands`, and update `/start` response text to include it.

**Mocks:** Use minimock (via `//go:generate`) for all interface mocks in tests — do not write hand-crafted fakes. Generate mock files into `package bot_test` using the `-p bot_test` flag so they are accessible from external test packages.

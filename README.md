# Telethings

A Telegram bot that seamlessly integrates with the Things 3 task management app, allowing you to quickly add tasks directly from Telegram.

## Features

- **Quick Task Addition**: Add tasks to Things 3 via simple Telegram commands
- **Command Discovery**: Built-in `/help` and `/start` commands with command menu integration
- **Flexible Modifiers**: Enhance your tasks with:
  - `when:` - Set task timing (e.g., `when:today`, `when:next friday`)
  - `tags:` - Organize with tags (e.g., `tags:errands,personal`)
  - `notes:` - Add detailed notes (e.g., `notes:"pick up oat milk"`)
- **Native Integration**: Uses Things 3's URL scheme for reliable communication
- **Simple Setup**: Single environment variable configuration

## Quick Start

### Prerequisites

- Go 1.25 or later
- Telegram bot token (create one via [@BotFather](https://t.me/BotFather))
- Things 3 app on macOS (with auth token for URL scheme access)

### Installation

Clone the repository:

```bash
git clone https://github.com/IlyasYOY/telethings.git
cd telethings
```

### Building

Build the executable:

```bash
go build -o telethings ./cmd/telethings
```

### Running

Set your Telegram bot token and allowed user IDs, then run:

```bash
export TELETHINGS_TELEGRAM_TOKEN=your_bot_token_here
export TELETHINGS_THINGS_AUTH_TOKEN=your_things_auth_token_here
export TELETHINGS_ALLOWED_USER_IDS=123456789,987654321
./telethings
```

The bot will start polling for messages. Only users in the `TELETHINGS_ALLOWED_USER_IDS` list can interact with the bot.

### Available Commands

Once the bot is running, you can interact with it using these commands:

- **`/start`** - Welcome message with quick command overview
- **`/help`** - Display detailed information about all available commands
- **`/add <title>`** - Add a task to Things 3 with optional modifiers

### Usage Examples

Send these messages to your Telegram bot:

```
/add Buy milk
→ Creates task: "Buy milk"

/add Complete project notes:remember to review
→ Creates task with detailed notes

/add Gym session when:tomorrow tags:fitness,health
→ Creates task scheduled for tomorrow with tags
```

Use `/help` to see all available commands and their usage directly in Telegram.

## Technology Stack

- **Language**: Go 1.25+
- **API**: [Telegram Bot API](https://core.telegram.org/bots/api) via [go-telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api)
- **Integration**: Things 3 URL scheme protocol

## Configuration

The bot requires the following environment variables:

- `TELETHINGS_TELEGRAM_TOKEN` - Your Telegram bot token (required)
- `TELETHINGS_THINGS_AUTH_TOKEN` - Things 3 URL scheme auth token (required)
- `TELETHINGS_ALLOWED_USER_IDS` - Comma-separated list of Telegram user IDs allowed to use the bot (required)

Example:
```bash
export TELETHINGS_TELEGRAM_TOKEN=123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11
export TELETHINGS_THINGS_AUTH_TOKEN=your_things_auth_token
export TELETHINGS_ALLOWED_USER_IDS=123456789,987654321,555555555
```

Only users whose Telegram IDs are in the `TELETHINGS_ALLOWED_USER_IDS` list can interact with the bot.

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.

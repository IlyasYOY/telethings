# macOS Background Setup

This guide shows how to run `telethings` in the background on macOS using `launchd`.

## Quick interactive setup (recommended)

Run:

```bash
make setup
```

It will:
- build the binary;
- ask for required environment variables;
- create env file and LaunchAgent plist;
- optionally load and start the service immediately.

To remove/uninstall the background setup interactively:

```bash
make setup-remove
```

It unloads/disables the LaunchAgent and can also delete plist/env files.

## 1) Build binary

```bash
make build
```

This creates `./bin/telethings`.

## 2) Create env file

Create a file, for example `~/.config/telethings/env`:

```bash
mkdir -p ~/.config/telethings
cat > ~/.config/telethings/env <<'EOF'
TELETHINGS_TELEGRAM_TOKEN=your_bot_token_here
TELETHINGS_THINGS_AUTH_TOKEN=your_things_auth_token_here
TELETHINGS_ALLOWED_USER_IDS=123456789,987654321
EOF
```

## 3) Create LaunchAgent plist

Create `~/Library/LaunchAgents/com.ilyasyoy.telethings.plist`:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>
    <key>Label</key>
    <string>com.ilyasyoy.telethings</string>

    <key>ProgramArguments</key>
    <array>
      <string>/bin/zsh</string>
      <string>-lc</string>
      <string>set -a; source ~/.config/telethings/env; set +a; /ABSOLUTE/PATH/TO/telethings/bin/telethings</string>
    </array>

    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>

    <key>StandardOutPath</key>
    <string>/tmp/telethings.out.log</string>
    <key>StandardErrorPath</key>
    <string>/tmp/telethings.err.log</string>
  </dict>
</plist>
```

Replace `/ABSOLUTE/PATH/TO/telethings` with your real repository path.

## 4) Load and start

```bash
launchctl bootstrap gui/$(id -u) ~/Library/LaunchAgents/com.ilyasyoy.telethings.plist
launchctl enable gui/$(id -u)/com.ilyasyoy.telethings
launchctl kickstart -k gui/$(id -u)/com.ilyasyoy.telethings
```

## 5) Check status and logs

```bash
launchctl print gui/$(id -u)/com.ilyasyoy.telethings | head -n 40
tail -f /tmp/telethings.out.log /tmp/telethings.err.log
```

## 6) Stop / disable

```bash
launchctl bootout gui/$(id -u) ~/Library/LaunchAgents/com.ilyasyoy.telethings.plist
launchctl disable gui/$(id -u)/com.ilyasyoy.telethings
```

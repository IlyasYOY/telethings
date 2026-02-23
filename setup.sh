#!/usr/bin/env bash
set -euo pipefail

if [[ "$(uname -s)" != "Darwin" ]]; then
  echo "This setup helper supports macOS only."
  exit 1
fi

repo_root="$(cd "$(dirname "$0")" && pwd)"
default_env="$HOME/.config/telethings/env"
default_label="com.ilyasyoy.telethings"
default_plist="$HOME/Library/LaunchAgents/${default_label}.plist"
binary_path="$repo_root/bin/telethings"

read -r -p "Build binary now with 'make build'? [Y/n] " build_now
build_now="${build_now:-Y}"
if [[ "$build_now" =~ ^[Yy]$ ]]; then
  make -C "$repo_root" build
fi

if [[ ! -x "$binary_path" ]]; then
  echo "Binary not found at: $binary_path"
  echo "Run 'make build' first."
  exit 1
fi

echo
echo "Enter bot configuration values."
read -r -p "TELETHINGS_TELEGRAM_TOKEN: " telegram_token
read -r -p "TELETHINGS_THINGS_AUTH_TOKEN: " things_token
read -r -p "TELETHINGS_ALLOWED_USER_IDS (comma-separated): " allowed_user_ids

read -r -p "Env file path [$default_env]: " env_path
env_path="${env_path:-$default_env}"

read -r -p "LaunchAgent label [$default_label]: " launchd_label
launchd_label="${launchd_label:-$default_label}"

read -r -p "LaunchAgent plist path [$default_plist]: " plist_path
plist_path="${plist_path:-$default_plist}"

mkdir -p "$(dirname "$env_path")" "$(dirname "$plist_path")"

cat >"$env_path" <<EOF
TELETHINGS_TELEGRAM_TOKEN=$telegram_token
TELETHINGS_THINGS_AUTH_TOKEN=$things_token
TELETHINGS_ALLOWED_USER_IDS=$allowed_user_ids
EOF
chmod 600 "$env_path"

escaped_env_path="${env_path//\"/\\\"}"
escaped_binary_path="${binary_path//\"/\\\"}"
cat >"$plist_path" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>
    <key>Label</key>
    <string>$launchd_label</string>
    <key>ProgramArguments</key>
    <array>
      <string>/bin/zsh</string>
      <string>-lc</string>
      <string>set -a; source "$escaped_env_path"; set +a; "$escaped_binary_path"</string>
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
EOF

uid="$(id -u)"
service_id="gui/$uid/$launchd_label"

echo
echo "Created env file:   $env_path"
echo "Created plist file: $plist_path"

read -r -p "Load and start LaunchAgent now? [Y/n] " start_now
start_now="${start_now:-Y}"
if [[ "$start_now" =~ ^[Yy]$ ]]; then
  launchctl bootout "gui/$uid" "$plist_path" >/dev/null 2>&1 || true
  launchctl bootstrap "gui/$uid" "$plist_path"
  launchctl enable "$service_id"
  launchctl kickstart -k "$service_id"
  echo "LaunchAgent started: $service_id"
  echo "Check logs with: tail -f /tmp/telethings.out.log /tmp/telethings.err.log"
else
  echo "Start later with:"
  echo "  launchctl bootstrap gui/$uid $plist_path"
  echo "  launchctl enable $service_id"
  echo "  launchctl kickstart -k $service_id"
fi

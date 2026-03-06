#!/usr/bin/env bash
set -euo pipefail

if [[ "$(uname -s)" != "Darwin" ]]; then
  echo "This update helper supports macOS only."
  exit 1
fi

repo_root="$(cd "$(dirname "$0")" && pwd)"
default_label="io.github.ilyasyoy.telethings"

echo "Building new binary..."
make -C "$repo_root" build

# Auto-detect installed plist by scanning ~/Library/LaunchAgents for any telethings plist
detected_label=""
while IFS= read -r -d '' f; do
  detected_label="$(basename "$f" .plist)"
  break
done < <(find "$HOME/Library/LaunchAgents" -maxdepth 1 -name '*telethings*.plist' -print0 2>/dev/null)

prompt_default="${detected_label:-$default_label}"
read -r -p "LaunchAgent label [$prompt_default]: " launchd_label
launchd_label="${launchd_label:-$prompt_default}"

uid="$(id -u)"
service_id="gui/$uid/$launchd_label"

echo "Restarting launchd service: $service_id"
launchctl kickstart -k "$service_id"

echo "Done. New binary is running."
echo "Check logs with: tail -f /tmp/telethings.out.log /tmp/telethings.err.log"

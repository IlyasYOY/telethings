#!/usr/bin/env bash
set -euo pipefail

if [[ "$(uname -s)" != "Darwin" ]]; then
  echo "This setup removal helper supports macOS only."
  exit 1
fi

default_label="io.github.ilyasyoy.telethings"
default_plist="$HOME/Library/LaunchAgents/${default_label}.plist"
default_env="$HOME/.config/telethings/env"

read -r -p "LaunchAgent label [$default_label]: " launchd_label
launchd_label="${launchd_label:-$default_label}"

read -r -p "LaunchAgent plist path [$default_plist]: " plist_path
plist_path="${plist_path:-$default_plist}"

read -r -p "Env file path [$default_env]: " env_path
env_path="${env_path:-$default_env}"

uid="$(id -u)"
service_id="gui/$uid/$launchd_label"

echo
echo "Stopping launchd service (if running): $service_id"
launchctl bootout "gui/$uid" "$plist_path" >/dev/null 2>&1 || true
launchctl disable "$service_id" >/dev/null 2>&1 || true

read -r -p "Delete plist file '$plist_path'? [y/N] " delete_plist
if [[ "$delete_plist" =~ ^[Yy]$ ]]; then
  rm -f "$plist_path"
  echo "Deleted plist: $plist_path"
fi

read -r -p "Delete env file '$env_path'? [y/N] " delete_env
if [[ "$delete_env" =~ ^[Yy]$ ]]; then
  rm -f "$env_path"
  echo "Deleted env file: $env_path"
fi

echo "Done."

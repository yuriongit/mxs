set -euo pipefail

FIRST_NAME="${1:-}"

echo "◌ Gifts a greeting and a fun fact"

if [[ -z "$FIRST_NAME" ]]; then
  echo "First name required"
  exit 1
fi

sleep 1

echo "Hello, $FIRST_NAME :)"

sleep 1

echo "Did you know, you can see a million colors?"

sleep 1

echo "Mhhhmmmm, including THESE ones 🔥"

sleep 1

echo "Well, that's your first script with MXS"

sleep 1

echo "✓ Done"

sleep 0.5

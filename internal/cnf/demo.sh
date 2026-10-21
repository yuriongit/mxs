set -euo pipefail

FIRST_NAME="${1:-}"

echo "◌ gifts a greeting and a fun fact"

if [[ -z "$FIRST_NAME" ]]; then
  echo "first name required"
  exit 1
fi

sleep 1

echo "hello, $FIRST_NAME :)"

sleep 1

echo "did you know, you can see a million colors?"

sleep 1

echo "mhhhmmmm, including these ones 🔥"

sleep 1

echo "well, that's your first script with mxs"

sleep 1

echo "✓ done"

sleep 0.5

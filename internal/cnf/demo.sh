set -euo pipefail

COLOR="${1:-}"

echo "◌ paints your output with personality"

if [[ -z "$COLOR" ]]; then
  echo "color required"
  exit 1
fi

sleep 1

echo "i like ${COLOR,,} too :)"

sleep 1

echo "did you know, you can see a million colors?"

sleep 1

echo "mhhhmmmm, including these ones 🔥"

sleep 1

echo "well, that's your first script with mxs"

sleep 1

echo "✓ done"

sleep 0.5

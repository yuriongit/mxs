# Preserve interleaved stdout and stderr output

This patch fixes a race condition in the executor where concurrent `stdout` and
`stderr` readers could garble or reorder output when a script writes to both
streams. Both streams are now read by a single goroutine through a shared line
writer, so output lines are forwarded in the order the streams produced them.

To install or upgrade MXS: **[mxs-cli.vercel.app/install](https://mxs-cli.vercel.app/install)**

---

## 🛠️ Fixes & Improvements

- **Interleaved output preserved:** The executor now reads `stdout` and `stderr`
  from a single goroutine using a mutex-serialized line writer. This eliminates
  the race from concurrent channel writes and preserves the order in which output
  lines arrived, so mixed output no longer gets garbled.

## 📚 Documentation

- Updated the PRD and README to cover the executor output fix and CLI behavior.

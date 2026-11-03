# Verbose listings, last-run tags, & a cleaner execution UI

This release introduces the `la` command, which lists every script with its
description and when it last ran. It also redesigns the `ls` output, refreshes
the execution screen, and reorganizes the internal packages for the upcoming
script-management commands.

To install or upgrade MXS: **[mxs-cli.vercel.app/install](https://mxs-cli.vercel.app/install)**

---

## ✨ What's New

- **`la` command:** List all scripts verbosely, with each script's description
  aligned next to its name.
- **Last-executed tags:** Every `la` entry now shows when the script last ran,
  persisted in `~/.mxs/last-executed.json`. Scripts that have never run display
  `l-exec: never executed (n/a)` in a subtle style.
- **Redesigned `ls`:** Script names now render in a colored box under a scripts
  header, without the old tree symbols and without the `.sh` suffix.
- **Execution UI refresh:** The execution screen now has dedicated description
  and logs sections, a live spinner and stopwatch, and clearer success and
  failure output.

## 🛠️ General Improvements

- **Internal cleanup:** Renamed the `cnf` package to `cfg` and `manager` to
  `scr`, and extracted script creation, listing, and description logic into
  focused files.
- **UI consolidation:** Reorganized and renamed Lipgloss styles and shared UI
  helpers for consistent rendering across commands.
- **CI hardening:** Tightened lint rules and added a path filter so non-Go
  changes skip the Go test pipeline.

## 📚 Documentation

- Updated the README and PRD to cover the `la` command, the new listing output,
  and the v0.4.0 behavior. The installers now point at the v0.4.0 artifacts.

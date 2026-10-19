#!/usr/bin/env python3

import re
import subprocess
from collections import defaultdict
from pathlib import Path

OUTPUT = Path("dist/CHANGELOG.md")
REPOSITORY = "yuriongit/mxs"

EXCLUDE = [
    r"^test:",
    r"^merge",
    r"^chore",
    r"merge conflict",
    r"typo",
    r"bump version",
    r"Merge pull request",
    r"Merge remote-tracking branch",
    r"Merge branch",
    r"go mod tidy",
]

SCOPE = r"(\([^)]*\))?"

GROUPS = [
    ("New!", re.compile(rf"^.*feat{SCOPE}:.*$")),
    ("Fixed", re.compile(rf"^.*fix{SCOPE}:.*$")),
    ("Maintenance & Refactoring", re.compile(rf"^(refactor|chore|ci){SCOPE}:")),
    ("Docs", re.compile(rf"^.*docs{SCOPE}:.*$")),
    ("Other stuff", None),
]


def git(*args):
    return subprocess.check_output(["git", *args], text=True).strip()


def main():
    try:
        previous_tag = git("describe", "--tags", "--abbrev=0")
        revision = f"{previous_tag}..HEAD"
    except subprocess.CalledProcessError:
        revision = "HEAD"

    log = subprocess.check_output(
        [
            "git",
            "log",
            revision,
            "--format=%H%x1f%s%x1f%an",
        ],
        text=True,
    )

    commits = []

    for line in log.splitlines():
        sha, message, author = line.split("\x1f", 2)

        if any(re.search(pattern, message) for pattern in EXCLUDE):
            continue

        commits.append((sha, message, author))

    commits.sort(key=lambda commit: commit[1])

    grouped = defaultdict(list)

    for commit in commits:
        _, message, _ = commit

        for title, pattern in GROUPS:
            if pattern is None or pattern.search(message):
                grouped[title].append(commit)
                break

    lines = ["# Changelog", ""]

    for title, pattern in GROUPS:
        entries = grouped[title]

        if not entries:
            continue

        lines.extend([f"## {title}", ""])

        for sha, message, author in entries:
            short_sha = sha[:7]
            url = f"https://github.com/{REPOSITORY}/commit/{sha}"

            lines.append(f"- [`{short_sha}`]({url}): {message} ({author})")

        lines.append("")

    OUTPUT.parent.mkdir(parents=True, exist_ok=True)
    OUTPUT.write_text("\n".join(lines).rstrip() + "\n")

    print(f"✓ Changelog generated at {OUTPUT}")


if __name__ == "__main__":
    main()

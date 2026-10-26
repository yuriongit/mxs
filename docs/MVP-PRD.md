# MXS PRD

Oct 17, 2026

## 1. Overview

MXS is a command-line tool for running and managing Bash scripts. Scripts are kept
as `.sh` files in `~/.mxs/scripts` and run by name from any directory. MXS also
adds, updates, renames, removes and lists scripts, so you never have to manage
that folder by hand.

## 2. Problem it solves

Personal shell scripts pile up in project folders, the home directory and old notes
A few months later it is hard to find the one you wrote, remember what it does,
or run it without digging up its path

- Scripts live in scattered locations, so finding one means searching for it
- Nothing records what a script is for, so the file name has to carry all the meaning
- Running a script means knowing its path and whether it is executable
- Renaming or deleting by hand is easy to get wrong, especially once other files
  refer to it or its description

## 3. Goals

- Keep all personal scripts in one place: `~/.mxs/scripts`
- Run any script by name from any directory
- Cover the full script life cycle with short commands: add, update, rename, remove,
  list
- Store a short description with each script
- Use the user's own editor for writing code instead of building one
- Keep the command set small and the behavior predictable

### Non-goals

- Replacing task runners or build tools such as Make
- Syncing or sharing scripts across machines
- Scheduling scripts to run later
- Running script types other than `.sh` files
- Tracking version history of scripts

## 4. Target User

MXS is for developers who work in the terminal and keep their own collection of
shell scripts on one machine, from a handful to a few dozen

- Uses the terminal daily and has a preferred editor
- Writes small scripts for tasks they repeat
- Wants something lighter than a task runner or a full dotfiles setup
- Works on a single local machine and does not need sharing or team features

## 5. Core Features

MXS has two groups of features: running a script, and managing the scripts in `~/.mxs/scripts`

### 5.1 Script execution

Running a script by name is the main use of MXS. The command is `mxs <name>`

- [x] Resolves `<name>` to `~/.mxs/scripts/<name>.sh` and runs it
- [x] Runs the script using the user's current working directory
- [x] Passes any arguments to the script
- [x] Connects the script's input and output directly to the terminal
- [x] Exits with the same exit code as the script
- [x] Prints successes and errors explicitly and clearly&#32;
- [x] Prints a clear error if the script does not exist

### 5.2 Script management (CRUD)

These commands create, read, update and delete scripts. Names are always shown and
accepted without the `.sh` suffix, so the name printed by `ls` is the name every
other command expects

#### add

`mxs add`

- [ ] Prompts for the script's name, then for a short description
- [ ] Creates `~/.mxs/scripts/<name>.sh`
- [ ] Saves the name and description to the metadata store (store not chosen yet)
- [ ] Opens the new file in the user's editor so they can write the script
- [ ] Refuses a name that is already in use

#### upd

`mxs upd <name>`, for example `mxs upd demo`

- [ ] Prompts for the name, then for the short description
- [ ] Prompts code changes: If "yes", the corresponding script file is opened in
      the user's editor. If the editor is TUI-based, a new terminal window will
      open the file. Both routes will use some type of `--wait` flag functionality,
      then the
      interactive editing process will end.
- [ ] Updates the metadata store with any changes
- [ ] Prints a clear error if the script does not exist

#### ren

`mxs ren <old> <new>`, for example `mxs ren lorem ipsum`

- [ ] Renames the script file and updates its entry in the metadata store
- [ ] Fails if `<old>` does not exist, `<new>` is already taken, or if no changes
      are provided

#### rm

`mxs rm <name>`, for example `mxs rm demo`

- [ ] Deletes the script file and removes its entry from the metadata store
- [ ] Prints a clear error if the script does not exist

#### ls

`mxs ls`

- [x] Lists the `.sh` files in `~/.mxs/scripts`
- [x] Shows each name without the `.sh` suffix, to avoid confusion when running,
      renaming,
      or updating a script
- [x] Ignores any file in the folder that does not end in `.sh`
- [x] Prints a clear error if no scripts exist

## 6. Success Metrics

MXS has no usage tracking, so these are checks the author can run directly. The
numeric targets are starting points and can change

- [ ] All six commands (`add`, `upd`, `ren`, `rm`, `ls`) work end to end on a fresh
      install
- [ ] A new script can be added and run in under one minute
- [ ] Every name printed by `ls` is accepted as-is by `upd`, `ren` and `rm`
- [ ] `ren` and `rm` never leave the script file and its metadata out of sync, covered
      by tests
- [ ] MXS adds less than 100 ms before a script starts running

## 7. Open questions

- Q: Which metadata store holds names and descriptions: a single JSON or TOML file,
  SQLite, or one small file per script?
  - Maybe SQLite
- Q: What is the command that runs a script: `exec`, `run`, or just `mxs <name>`?
  - A: `mxs <name`
- Q: `upd` prompts for the name and `ren` also changes it. Should `upd` keep the
- name prompt, or leave renaming to `ren`?
  - A: Both should remain. Why? Well, `ren` provides a smoother UX because the user
    doesn't need to enter an interactive state just to rename a single script.
    Rather, when updating a script, the `upd` command provides the user with all
    that customizability
- Q: What happens to a `.sh` file added to the folder by hand? It shows up in `ls`
  but has no description?
  - A: Not entirely sure right now. I think maybe each script should be registered
    in the data-store, and any unregistered scripts will be immediately identified
    as unregistered. Those scripts will be displayed as red for `ls`'s output
- Q: Should `rm` ask for confirmation before deleting?
  - A:Yes
- Should `ls` show each script's description next to its name?
  - No, there should be a `la` command for that
- Which editor does MXS open: `$EDITOR`, `$VISUAL`, or a fallback if neither is
  set?
  - A: `$EDITOR` as primary, with `$VISUAL` as the fallback. No other fallbacks
    will be provided.
- Q: What characters are allowed in a script name (spaces, slashes, leading periods)?
  - A: Script names are allowed to be in the form lowercase characters, and can
    contain symbols other than periods, spaces, and slashes.
- Q: Is a script run through `bash` directly, or executed using its own shebang
  line?
  - A: Scripts are executed using bash directly, not their shebang line.
- Q: What is MXS written in, and how is it installed?
  - A: Go. It can be built from source and installed to the Go bin. Additionally,
    users can pick up a pre-built binary that can be downloaded and installed from
    the GitHub releases page, or use the provided installation script for automated
    setup.

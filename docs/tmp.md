# Tmp

Accidental commented out note left in README

Add custom configured colors: If there's any errors
in the config, they must first be resolved. Until
then, a script will not run.

Add this functionality to the cfg package.

These following colors will have a default color.

Also, in the cfg file, add a themes section:

```json
theme: {
  colors: {
    errors: ["#LIGHT", "#DARK"],
    execution_progress: ["#LIGHT", "#DARK"],
    successful_execution: ["#LIGHT", "#DARK"]
  },
  spinner: {
    extraspace: false,
    type: "dots | monkey | meter",
    color: "#HEX_CODE"
  },
  verbosity: {
    supplied_args: true,
    stopwatch: true
  },
  logs: {
    quotes: true
  }
}
```

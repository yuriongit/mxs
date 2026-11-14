Output "./demos/cfg/cfg.webm"

Require mxs

Set Shell "bash"
Set FontSize 32

Set Width 1800

Set Framerate 15
Set Height 900

# Setup
Hide
Type "mxs init" Enter
Sleep 250ms
Type "clear" Enter
Sleep 250ms

# Demonstration
Show
Type "mxs cfg" Sleep 500ms Enter
Sleep 2s

Type ":wq" Enter Sleep 500ms

Sleep 3s

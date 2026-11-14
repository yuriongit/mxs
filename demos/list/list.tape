Output "./demos/list/list.webm"

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
Type "mxs add foo" Enter
Type ":wq" Enter Sleep 250ms
Type "mxs demo" Enter
Sleep 250ms
Type "clear" Enter
Sleep 250ms

# Demonstration
Show
Type "mxs ls" Sleep 500ms Enter
Sleep 2s

Type "mxs la" Sleep 500ms Enter
Sleep 5s

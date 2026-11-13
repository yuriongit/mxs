Output "./demos/add/add.gif"

Require mxs
Require vim

Set Shell "bash"
Set FontSize 32

Set Width 1800

Set Framerate 15
Set Height 900

# Setup
Hide
Type "mxs init" Enter
Sleep 500ms
Type "clear" Enter
Sleep 250ms

# Demonstration
Show
Type "mxs add foo" Sleep 500ms Enter
Sleep 2s
Type ":wq" Sleep 500ms Enter
Sleep 5s

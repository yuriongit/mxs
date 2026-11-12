Output "./demos/add/add.gif"

Require mxs

Set Shell "bash"
Set FontSize 26

Set Width 1800
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
Sleep 5s

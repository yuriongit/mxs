Output "./demos/list/list.gif"

Require mxs

Set Shell "bash"
Set FontSize 26

Set Width 1800
Set Height 900

# Setup
Hide
Type "mxs init" Enter
Type "mxs add foo" Enter
Type "mxs add bar" Enter
Type "mxs bar" Enter
Sleep 500ms
Type "clear" Enter
Sleep 250ms

# Demonstration
Show
Type "mxs ls" Sleep 500ms Enter
Sleep 2s

Type "mxs la" Sleep 500ms Enter
Sleep 5s

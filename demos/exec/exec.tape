Output "./demos/exec/exec.webm"
Output "./demos/exec/exec.gif" # For README embedding compatability

Require mxs

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
Type "mxs demo" Sleep 500ms Enter
Sleep 1.5s
Ctrl+L

Type "mxs demo purple" Sleep 500ms Enter
Sleep 8s

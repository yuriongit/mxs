Output "./demos/open/open.gif"

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
Type "mxs foo" Enter
Sleep 250ms
Type "clear" Enter
Sleep 250ms

# Demonstration
Show

Type "mxs la" Sleep 500ms Enter
Sleep 1s

Type "mxs open foo" Sleep 500ms Enter
Sleep 500ms
Set TypingSpeed 75ms
Type 'G$ ci"' Sleep 250ms
Set TypingSpeed 75ms
Type "opened foo and made some edits"
Sleep 1s
Set TypingSpeed 50ms
Escape
Type ":wq" Enter Sleep 500ms
Sleep 0.5s

Type "clear" Enter Sleep 500ms
Sleep 0.75s
Type "mxs foo" Sleep 500ms Enter
Sleep 5s

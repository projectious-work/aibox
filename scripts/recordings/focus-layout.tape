# focus-layout.tape — capture the "focus" Zellij layout
# Run: vhs scripts/recordings/focus-layout.tape

Output docs/assets/screenshots/layout-focus.gif
Set Width 1400
Set Height 900
Set FontSize 14
Set Padding 10
Set Shell bash
Set TypingSpeed 50ms

Hide
Type "export TERM=xterm-256color COLORTERM=truecolor"
Enter
Sleep 500ms
Show

Type "zellij --layout focus"
Enter
Sleep 4s

Screenshot docs/assets/screenshots/layout-focus.png
Sleep 1s

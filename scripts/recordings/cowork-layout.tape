# cowork-layout.tape — capture the "cowork" Zellij layout
# Run: vhs scripts/recordings/cowork-layout.tape

Output docs/assets/screenshots/layout-cowork.gif
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

Type "zellij --layout cowork"
Enter
Sleep 4s

Screenshot docs/assets/screenshots/layout-cowork.png
Sleep 1s

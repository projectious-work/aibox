# dev-layout.tape — capture the default "dev" Zellij layout
# Run: vhs scripts/recordings/dev-layout.tape

Output docs/assets/screenshots/layout-dev.gif
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

Type "zellij --layout dev"
Enter
Sleep 4s

Screenshot docs/assets/screenshots/layout-dev.png
Sleep 1s

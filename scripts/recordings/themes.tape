# themes.tape — capture all 6 color themes in the dev layout
# This tape is run inside a dev-box container.
# The record-screenshots.sh script swaps the Zellij theme before each run.
#
# Usage (called by record-screenshots.sh, not directly):
#   DEVBOX_THEME=gruvbox-dark vhs scripts/recordings/themes.tape

Output docs/assets/screenshots/theme-current.gif
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

Screenshot docs/assets/screenshots/theme-current.png
Sleep 1s

# init-demo.tape — capture the dev-box init workflow
# Run: vhs scripts/recordings/init-demo.tape

Output docs/assets/screenshots/init-demo.gif
Set Width 1000
Set Height 600
Set FontSize 14
Set Padding 10
Set Shell bash
Set TypingSpeed 80ms

Type "mkdir my-project && cd my-project"
Enter
Sleep 500ms

Type "dev-box init --name my-project --image python --process managed"
Enter
Sleep 3s

Screenshot docs/assets/screenshots/init-output.png

Sleep 1s
Type "cat dev-box.toml"
Enter
Sleep 2s

Screenshot docs/assets/screenshots/init-toml.png
Sleep 1s

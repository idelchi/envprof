Output gifs/envprof.gif

Set Theme "CLRS"
Set FontSize 16
Set Width 1200
Set Height 650
Set TypingSpeed 50ms

Env SHELL "bash"

Type 'envprof ls'

Sleep 500ms
Enter
Sleep 1s

Type 'envprof ls dev -v'

Sleep 500ms
Enter
Sleep 1s

Type 'envprof export dev dev.env'

Sleep 500ms
Enter
Sleep 1s

Type 'cat dev.env'

Sleep 500ms
Enter
Sleep 1s

Type 'echo HOST=$HOST'

Sleep 500ms
Enter
Sleep 1s

Type 'envprof shell dev'

Sleep 500ms
Enter
Sleep 1s

Type 'echo HOST=$HOST'

Sleep 500ms
Enter
Sleep 1s

Sleep 2s

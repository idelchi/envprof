Output gifs/envprof.gif

Set Theme "CLRS"
Set FontSize 16
Set Width 1200
Set Height 950
Set TypingSpeed 50ms

Env SHELL "bash"
Env ENVPROF_FILE "envprof.toml"

Type '# List all profiles'
Enter
Type 'envprof profs -v'
Sleep 500ms
Enter

Sleep 1s
Enter

Type '# List all variables in a profile'
Enter
Type 'envprof ls -v'
Sleep 500ms
Enter

Sleep 1s
Enter

Type '# Write a profile to a dotenv file'
Enter
Type 'envprof write dev.env'
Sleep 500ms
Enter
Sleep 1s
Type 'cat dev.env'
Sleep 500ms
Enter

Sleep 1s
Enter

Type '# Execute a command with a profile'
Enter
Type 'envprof exec -i -- printenv'
Sleep 500ms
Enter

Sleep 1s
Enter

Type '# Enter a shell with an active profile'
Enter
Type 'echo HOST=$HOST'
Sleep 500ms
Enter
Sleep 1s
Type 'envprof shell'
Sleep 500ms
Enter
Sleep 1s
Type 'echo HOST=$HOST'
Sleep 500ms
Enter

Sleep 1s

# Finally
Sleep 2s

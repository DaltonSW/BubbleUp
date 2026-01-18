# Record a window demo (you’ll capture it with RetroClip + KeyCastr).
# VHS will still “run” the sequence deterministically.

Output example_main.gif

Set Shell "zsh"
Set FontSize 18
Set FontFamily "JetBrainsMono Nerd Font Mono"
Set Width 1024
Set Height 640
Set Theme "Dracula"

# Optional: add a little breathing room at start
Sleep 1500ms

# Run the example app
Type "go run example_main.go"
Enter

# Give the program time to start and render
Sleep 1500ms

############################
# Unicode (default)
############################
Type "i"
Sleep 2500ms

Type "w"
Sleep 2500ms

Type "e"
Sleep 2500ms

Type "d"
Sleep 2500ms

Type "I"
Sleep 2500ms
Escape
Sleep 2500ms

Type "W"
Sleep 2500ms
Escape

########################################
# Switch to NerdFont mode
########################################

Sleep 3000ms
Type "N"
Sleep 1500ms

Type "i"
Sleep 1500ms
Escape
Sleep 1500ms

Type "w"
Sleep 1500ms
Escape
Sleep 1500ms

Type "e"
Sleep 1500ms
Escape
Sleep 2500ms

Type "d"
Sleep 1500ms
Escape
Sleep 1500ms

Type "I"
Sleep 1500ms

Type "W"
Sleep 2500ms

########################################
# Switch to ASCII mode
########################################

Sleep 3000ms
Type "A"
Sleep 1500ms

Type "i"
Sleep 1500ms

Type "w"
Sleep 1500ms

Type "e"
Sleep 2500ms

Type "d"
Sleep 1500ms

Type "I"
Sleep 1500ms

Type "W"
Sleep 2500ms


########################################
# Quit
########################################

Type "q"
Sleep 3000ms

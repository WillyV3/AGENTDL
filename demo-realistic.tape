# Agent Search Demo - Realistic User Flow
# Natural pacing showing real user interaction

# Output configuration
Output demo-realistic.gif
Output demo-realistic.mp4

# Terminal setup for better visibility
Set FontSize 16
Set Width 1400
Set Height 900
Set Padding 30
Set Margin 30
Set MarginFill "#1F2937"
Set BorderRadius 12
Set WindowBar Colorful
Set Theme "Catppuccin Mocha"
Set Framerate 30
Set PlaybackSpeed 1.0

# Natural typing speed
Set TypingSpeed 100ms



# Start the application
Type "./search"
Enter
Sleep 3s  # Let the UI fully load

# Pause to show the beautiful search screen with animation
Sleep 2s

# Type search query naturally
Type@150ms "bash"
Sleep 1s  # User thinking/looking at screen

# Execute search
Enter


# User examining results
Sleep 2s

# View the top repository
Type "v"
Sleep 3s  # Loading repo structure

Enter # enter .claude directory 
Sleep 2s
Enter 
Enter
Sleep 2s
Enter
Sleep 3s
Down
Down
Down
Down
Down
Down
Down
Down
Down
Sleep 2s


Escape
Sleep 2
Escape
Sleep 1s
# select .md files
Space
Sleep 1s
Down
Sleep 1s
Space
Sleep 1s
Down
Sleep 1s
Space
Sleep 1s
Down
Sleep 1s
# go to download
Enter
Sleep 3s
# view files for download
Type "v"
Sleep 2s
Escape
Sleep 1s 
Enter
Sleep 2s
Enter

Type@150ms "animation"
Sleep 1s    
Enter 
Sleep 4s 
Space
Down
Sleep 1s 
Space
Sleep 1s
Down

Sleep 1s
Escape
Sleep 3s
Up
Sleep 2s
Enter
Sleep 5s

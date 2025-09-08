# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

GitHub Agent Search TUI - A Bubble Tea application for discovering and downloading Claude agent configurations (`.claude/agents/*.md` files) from GitHub repositories with intelligent keyword search, file preview, and multi-selection capabilities.

## File Structure & Architecture

### Core Principles
- **Modular Design**: Each file has a single responsibility (<300 lines per file)
- **Bubble Tea Patterns**: Leverage framework features (viewport, forms, spinners)
- **Single Source of Truth**: Global SelectionManager for all file selections
- **Clean Separation**: Views, Updates, and Commands in separate files

### File Organization
```
agent-search/
├── search.go          # Core model, states, and orchestrators (283 lines)
├── views.go           # All view rendering functions (358 lines)
├── updates.go         # State update handlers (383 lines)
├── commands.go        # Async tea.Cmd operations (104 lines)
├── github_viewer.go   # Repository browser component (248 lines)
├── selections.go      # Global selection management (120 lines)
└── github/
    └── github.go      # GitHub search module (100 lines)
```

## Development Commands

```bash
# Build and run (compile all .go files together)
go build -o search *.go
./search

# Development with hot reload
go run *.go

# Clean build
go mod tidy
go build -o search *.go
```

## Working with the Codebase

### Adding New Features
1. **State Machine**: Add new state in `search.go` constants
2. **Update Handler**: Add handler in `updates.go` following pattern:
   ```go
   func (m model) updateNewState(msg tea.Msg) (tea.Model, tea.Cmd)
   ```
3. **View Function**: Add view in `views.go` following pattern:
   ```go
   func (m model) viewNewState() string
   ```
4. **Wire Up**: Add cases in main Update() and View() orchestrators

### Bubble Tea Best Practices

#### Message Handling
- Global messages (WindowSize, Ctrl+C, spinner.TickMsg) handled in main Update()
- State-specific messages handled in dedicated update functions
- Never block spinner.TickMsg - return immediately with cmd

#### Component Patterns
- Use `viewport.Model` for scrollable content
- Use `huh.Form` for user selections (not custom key handling)
- Use `textinput.Model` for text entry
- Use `spinner.Model` for loading states

#### Styling with Lipgloss
```go
// Consistent theme usage
var theme = struct {
    primary   lipgloss.Color
    secondary lipgloss.Color
    // ...
}

// Centered, bordered boxes
boxStyle := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(theme.primary).
    Padding(2, 4).
    Width(60)

// Center on screen
lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
```

### Selection Management

All selections go through `SelectionManager`:
```go
// Toggle selection
sel := GlobalSelection{
    Repo:     result.Repo,
    Path:     result.Path,
    URL:      result.URL,
    FileName: ExtractFileName(result.Path),
    Source:   "search", // or "repo"
}
m.globalSelections.Toggle(sel)

// Check if selected
if m.globalSelections.IsSelected(repo, path) { ... }

// Get count
totalCount := m.globalSelections.Count()
```

### Async Operations

Use tea.Cmd for all async work:
```go
func searchGitHub(query string, mode string) tea.Cmd {
    return func() tea.Msg {
        // Perform async operation
        results := github.Search(query, opts)
        return searchResultsMsg{results: results}
    }
}
```

## Key Implementation Details

### GitHub Integration
- **Search**: Intelligent AND/OR keyword matching via `gh search code`
- **Repo Browser**: Navigate repositories using `gh api`
- **Downloads**: Convert blob URLs to raw.githubusercontent.com

### User Experience
- **Confirmation Dialogs**: Button-style selection with up/down navigation
- **Help Text**: Compact two-line format to prevent overflow
- **Visual Feedback**: Selected items show [✓], spinners for loading
- **Smart Defaults**: "oh shit, go back" selected by default in confirmations

### State Transitions
```
stateSearch -> stateSearching -> stateResults
                                      ↓
                              stateRepoViewer
                                      ↓
                              stateLocation -> stateDownloading -> stateComplete
```

## Common Patterns

### Adding a Confirmation Dialog
```go
// 1. Add state
const stateConfirm state = iota

// 2. Add field to model
confirmChoice int // 0 or 1

// 3. Update handler with up/down/enter
func (m model) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd)

// 4. View with button-style options
options := []string{"Yes", "No"}
for i, opt := range options {
    if i == m.confirmChoice {
        // Highlighted with ▶ and background color
    }
}
```

### File Size Limits
- Keep files under 300 lines
- Extract related functions into modules
- Use separate files for distinct features

## Testing

```bash
# Basic flow
./search
# Search: "tui" or "bash"
# Press 'v' to browse repo
# Space to select files
# Enter to download

# Test selection persistence
# Select files in repo browser
# Go back to results
# Selections remain until new search
```

## Dependencies

- Go 1.20+
- `gh` CLI (authenticated)
- Bubble Tea framework
- Lipgloss for styling

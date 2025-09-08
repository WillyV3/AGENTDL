package main

import (
	"fmt"
	"os"
	"path/filepath"

	"agent-search/github"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ============================
// Styles (grouped together)
// ============================

var (
	theme = struct {
		primary   lipgloss.Color
		secondary lipgloss.Color
		success   lipgloss.Color
		error     lipgloss.Color
		muted     lipgloss.Color
		bg        lipgloss.Color
	}{
		primary:   lipgloss.Color("#06B6D4"), // Teal/Cyan
		secondary: lipgloss.Color("#FCD34D"), // Yellow/Amber
		success:   lipgloss.Color("#10B981"), // Emerald green
		error:     lipgloss.Color("#EF4444"), // Red
		muted:     lipgloss.Color("#6B7280"), // Gray
		bg:        lipgloss.Color("#1F2937"), // Dark gray
	}

	titleStyle = lipgloss.NewStyle().
			Foreground(theme.primary).
			Bold(true).
			Padding(1, 0).
			Align(lipgloss.Center)

	selectedStyle = lipgloss.NewStyle().
			Foreground(theme.secondary).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA"))

	helpStyle = lipgloss.NewStyle().
			Foreground(theme.muted).
			Margin(1, 0)

	errorStyle = lipgloss.NewStyle().
			Foreground(theme.error).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(theme.success).
			Bold(true)
)

// ============================
// Domain Types
// ============================

type searchResult github.Result

type locationOption int

const (
	locationGlobal locationOption = iota
	locationCurrent
	locationCustom
)

var locationPaths = map[locationOption]func() string{
	locationGlobal: func() string {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".claude", "agents")
	},
	locationCurrent: func() string {
		return filepath.Join(".", ".claude", "agents")
	},
}

// ============================
// States
// ============================

type state int

const (
	stateSearch state = iota
	stateSearching
	stateResults
	statePreview
	stateLocation
	stateLocationFileList
	stateCustomPath
	stateDownloading
	stateComplete
	stateRepoViewer
	stateConfirmLoseSelections
)

// ============================
// Messages
// ============================

type searchResultsMsg struct {
	results []searchResult
	err     error
}

type fileContentMsg struct {
	content string
	err     error
}

type downloadCompleteMsg struct {
	count int
}

// ============================
// Model
// ============================

type model struct {
	state            state
	searchInput      textinput.Model
	results          []searchResult
	cursor           int
	err              error
	width            int
	height           int
	viewport         viewport.Model
	previewContent   string
	location         locationOption
	customPath       string
	customPathInput  textinput.Model
	repoViewer       *RepoViewer
	locationChoice   int               // 0=global, 1=current, 2=custom
	globalSelections *SelectionManager // Global selection manager
	returnToState    state             // State to return to after confirmation
	confirmChoice    int               // 0 = idgaf, 1 = oh shit go back
	resultsOffset    int               // Scroll offset for results list
	fileListCursor   int               // Separate cursor for file list view
}

// ============================
// Model Lifecycle Methods
// ============================

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter search keyword..."
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50

	customInput := textinput.New()
	customInput.Placeholder = "/path/to/download"
	customInput.CharLimit = 256
	customInput.Width = 50

	vp := viewport.New(80, 20)

	return model{
		state:            stateSearch,
		searchInput:      ti,
		results:          []searchResult{},
		cursor:           0,
		width:            80,
		height:           24,
		viewport:         vp,
		customPathInput:  customInput,
		globalSelections: NewSelectionManager(),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// ============================
// Main Update Orchestrator
// ============================

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle global messages
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 4
		// Adjust results scrolling on resize to keep cursor visible
		if m.state == stateResults {
			maxVisible := m.height - 8
			if maxVisible < 5 {
				maxVisible = 5
			}
			if maxVisible > 30 {
				maxVisible = 30
			}
			maxOffset := 0
			if len(m.results) > maxVisible {
				maxOffset = len(m.results) - maxVisible
			}
			if m.resultsOffset > maxOffset {
				m.resultsOffset = maxOffset
			}
			if m.cursor < m.resultsOffset {
				m.resultsOffset = m.cursor
			} else if m.cursor >= m.resultsOffset+maxVisible {
				m.resultsOffset = m.cursor - maxVisible + 1
			}
		}
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	// Route to state-specific handlers
	switch m.state {
	case stateSearch:
		return m.updateSearch(msg)
	case stateSearching:
		return m.updateSearching(msg)
	case stateResults:
		return m.updateResults(msg)
	case statePreview:
		return m.updatePreview(msg)
	case stateLocation:
		return m.updateLocation(msg)
	case stateLocationFileList:
		return m.updateLocationFileList(msg)
	case stateCustomPath:
		return m.updateCustomPath(msg)
	case stateDownloading:
		return m.updateDownloading(msg)
	case stateComplete:
		return m.updateComplete(msg)
	case stateRepoViewer:
		return m.updateRepoViewer(msg)
	case stateConfirmLoseSelections:
		return m.updateConfirmLoseSelections(msg)
	}

	return m, nil
}

// ============================
// Main View Orchestrator
// ============================

func (m model) View() string {
	switch m.state {
	case stateSearch:
		return m.viewSearch()
	case stateSearching:
		return m.viewSearching()
	case stateResults:
		return m.viewResults()
	case statePreview:
		return m.viewPreview()
	case stateLocation:
		return m.viewLocation()
	case stateLocationFileList:
		return m.viewLocationFileList()
	case stateCustomPath:
		return m.viewCustomPath()
	case stateDownloading:
		return m.viewDownloading()
	case stateComplete:
		return m.viewComplete()
	case stateRepoViewer:
		return m.viewRepoViewer()
	case stateConfirmLoseSelections:
		return m.viewConfirmLoseSelections()
	default:
		return "Unknown state"
	}
}

// ============================
// Main Entry Point
// ============================

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

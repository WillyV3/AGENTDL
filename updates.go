package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// ============================
// Update Handlers
// ============================

func (m model) updateSearch(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyTab:
			// Toggle search mode
			if m.searchMode == modeAgents {
				m.searchMode = modeCommands
			} else {
				m.searchMode = modeAgents
			}
			return m, nil
		case tea.KeyEnter:
			if m.searchInput.Value() != "" {
				// Clear global selections when starting a new search
				m.globalSelections.Clear()
				m.state = stateSearching
				return m, searchGitHub(m.searchInput.Value(), "all", m.searchMode)
			}
			return m, nil
		}
		// Special handling for "q" to quit
		if msg.String() == "q" {
			return m, tea.Quit
		}
	}

	// ALWAYS update text input for ALL messages (including KeyMsg that fall through)
	m.searchInput, cmd = m.searchInput.Update(msg)
	return m, cmd
}

func (m model) updateSearching(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case searchResultsMsg:
		if msg.err != nil {
			m.err = msg.err
			m.state = stateSearch
			return m, nil
		}
		m.results = msg.results
		m.cursor = 0
		m.resultsOffset = 0
		if len(m.results) > 0 {
			m.state = stateResults
		} else {
			m.err = fmt.Errorf("no results found")
			m.state = stateSearch
		}
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "esc" {
			m.state = stateSearch
			return m, nil
		}
	}

	return m, nil
}

func (m model) updateResults(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			// Check if we have global selections
			if m.globalSelections.Count() > 0 {
				m.returnToState = stateSearch
				m.state = stateConfirmLoseSelections
				m.confirmChoice = 1 // Default to "oh shit, go back"
				return m, nil
			}
			m.state = stateSearch
			m.cursor = 0
			m.resultsOffset = 0
			m.results = []searchResult{}
			return m, nil

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			// adjust scroll offset
			if m.cursor < m.resultsOffset {
				m.resultsOffset = m.cursor
			}

		case "down", "j":
			if m.cursor < len(m.results)-1 {
				m.cursor++
			}
			// adjust scroll offset based on visible window
			maxVisible := m.height - 8
			if maxVisible < 5 {
				maxVisible = 5
			}
			if maxVisible > 30 {
				maxVisible = 30
			}
			if m.cursor >= m.resultsOffset+maxVisible {
				m.resultsOffset = m.cursor - maxVisible + 1
			}

		case " ", "space":
			if m.cursor < len(m.results) {
				result := m.results[m.cursor]
				sel := GlobalSelection{
					Repo:     result.Repo,
					Path:     result.Path,
					URL:      result.URL,
					FileName: ExtractFileName(result.Path),
					Source:   "search",
				}
				m.globalSelections.Toggle(sel)
			}

		case "a":
			// Check if all are selected
			allSelected := true
			for _, r := range m.results {
				if !m.globalSelections.IsSelected(r.Repo, r.Path) {
					allSelected = false
					break
				}
			}
			// Toggle all
			for _, r := range m.results {
				sel := GlobalSelection{
					Repo:     r.Repo,
					Path:     r.Path,
					URL:      r.URL,
					FileName: ExtractFileName(r.Path),
					Source:   "search",
				}
				if allSelected {
					m.globalSelections.Remove(r.Repo, r.Path)
				} else {
					m.globalSelections.Add(sel)
				}
			}

		case "p":
			// Preview file
			if m.cursor < len(m.results) {
				result := m.results[m.cursor]
				m.state = statePreview
				return m, fetchFileContent(result.URL)
			}

		case "v":
			// View repository
			if m.cursor < len(m.results) {
				result := m.results[m.cursor]
				viewer := NewRepoViewer(result.Repo, result.Stars, m.globalSelections)
				m.repoViewer = &viewer
				m.state = stateRepoViewer
				return m, m.repoViewer.Init()
			}

		case "enter":
			// Check if we have any selections
			if m.globalSelections.Count() > 0 {
				m.state = stateLocation
				m.locationChoice = 0 // Reset to first option
				return m, nil
			}
		}
	}

	return m, nil
}

func (m model) updatePreview(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case fileContentMsg:
		if msg.err != nil {
			m.err = msg.err
			m.state = stateResults
			return m, nil
		}
		m.previewContent = msg.content
		m.viewport.SetContent(msg.content)
		m.viewport.GotoTop()
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			m.state = stateResults
			m.previewContent = ""
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m model) updateLocation(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.state = stateResults
			return m, nil
		case "v":
			// View file list
			m.state = stateLocationFileList
			m.fileListCursor = 0
			return m, nil
		case "up", "k":
			if m.locationChoice > 0 {
				m.locationChoice--
			}
		case "down", "j":
			if m.locationChoice < 2 {
				m.locationChoice++
			}
		case "enter":
			switch m.locationChoice {
			case 0: // Global
				m.location = locationGlobal
				m.state = stateDownloading
				return m, downloadSelectedFiles(m.results, locationPaths[locationGlobal][m.searchMode](), m.globalSelections)
			case 1: // Current
				m.location = locationCurrent
				m.state = stateDownloading
				return m, downloadSelectedFiles(m.results, locationPaths[locationCurrent][m.searchMode](), m.globalSelections)
			case 2: // Custom
				m.state = stateCustomPath
				m.customPathInput.Focus()
				return m, textinput.Blink
			}
		}
	}

	return m, nil
}

func (m model) updateLocationFileList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "v":
			m.state = stateLocation
			return m, nil
		}
	}
	return m, nil
}

func (m model) updateCustomPath(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.state = stateLocation
			m.locationChoice = 0
			return m, nil

		case "enter":
			if m.customPathInput.Value() != "" {
				path := m.customPathInput.Value()
				m.state = stateDownloading
				return m, downloadSelectedFiles(m.results, path, m.globalSelections)
			}
		}
	}

	var cmd tea.Cmd
	m.customPathInput, cmd = m.customPathInput.Update(msg)
	return m, cmd
}

func (m model) updateDownloading(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case downloadCompleteMsg:
		m.state = stateComplete
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "esc" {
			m.state = stateLocation
			m.locationChoice = 0
			return m, nil
		}
		// Don't handle other keys
	}

	return m, nil
}

func (m model) updateComplete(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "enter":
			// Reset and go back to search
			m.state = stateSearch
			m.results = []searchResult{}
			m.cursor = 0
			m.searchInput.SetValue("")
			m.searchInput.Focus()
			m.customPath = ""
			return m, nil
		}
	}

	return m, nil
}

func (m model) updateRepoViewer(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.repoViewer == nil {
		return m, nil
	}

	// Check for back to results message
	if _, ok := msg.(backToResultsMsg); ok {
		m.state = stateResults
		m.repoViewer = nil
		return m, nil
	}

	// Update the repo viewer
	viewer, cmd := m.repoViewer.Update(msg)
	m.repoViewer = &viewer
	return m, cmd
}

func (m model) updateConfirmLoseSelections(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.confirmChoice = 0
			return m, nil

		case "down", "j":
			m.confirmChoice = 1
			return m, nil

		case "enter", " ":
			if m.confirmChoice == 0 {
				// idgaf - Clear selections and proceed
				m.globalSelections.Clear()
				if m.returnToState == stateSearch {
					m.state = stateSearch
					m.cursor = 0
					m.results = []searchResult{}
				} else {
					m.state = m.returnToState
				}
			} else {
				// oh shit, go back - Return to results
				m.state = stateResults
			}
			return m, nil

		case "esc":
			// Escape always goes back
			m.state = stateResults
			return m, nil
		}
	}
	return m, nil
}

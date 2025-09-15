package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ============================
// View Methods
// ============================

func (m model) viewSearch() string {
	// Simple title
	title := lipgloss.NewStyle().
		Foreground(theme.primary).
		Bold(true).
		Render("AGENTDL")

	subtitle := lipgloss.NewStyle().
		Foreground(theme.muted).
		Italic(true).
		Render("Discover Claude Agents on GitHub")

	// Mode indicator
	var modeText string
	var modeStyle lipgloss.Style
	if m.searchMode == modeAgents {
		modeText = "[Agents]"
		modeStyle = lipgloss.NewStyle().Foreground(theme.secondary).Bold(true)
	} else {
		modeText = "[Commands]"
		modeStyle = lipgloss.NewStyle().Foreground(theme.success).Bold(true)
	}
	modeIndicator := modeStyle.Render(modeText)

	// Build content using simple formatting
	var content string
	if m.err != nil {
		content = fmt.Sprintf("\n\n%s\n%s\n\n%s\n\n%s\n\n%s\n\n%s\n\n%s",
			title,
			subtitle,
			modeIndicator,
			m.searchInput.View(),
			errorStyle.Render(fmt.Sprintf("⚠ %v", m.err)),
			helpStyle.Render("Tab: toggle mode • Enter: search • Esc: quit"),
			lipgloss.NewStyle().Foreground(theme.muted).Italic(true).Render("Made w/ ♥ by WillyV3"),
		)
	} else {
		content = fmt.Sprintf("\n\n%s\n%s\n\n%s\n\n%s\n\n%s\n\n%s",
			title,
			subtitle,
			modeIndicator,
			m.searchInput.View(),
			helpStyle.Render("Tab: toggle mode • Enter: search • Esc: quit"),
			lipgloss.NewStyle().Foreground(theme.muted).Italic(true).Render("Made w/ ♥ by WillyV3"),
		)
	}

	// Return content without lipgloss.Place() or box styling for now
	return content
}

func (m model) viewSearching() string {
	title := titleStyle.Render("Searching GitHub...")
	loadingText := "Finding .claude/agents files..."

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		loadingText,
		"",
		helpStyle.Render("Esc: cancel"),
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m model) viewResults() string {
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render(fmt.Sprintf("Found %d agent files", len(m.results))))
	b.WriteString("\n\n")

	// Stable scrolling: keep an offset and only scroll
	// when the cursor leaves the visible window
	maxVisible := m.height - 8 // leave room for title/help/padding
	if maxVisible < 5 {
		maxVisible = 5
	}
	if maxVisible > 30 {
		maxVisible = 30
	}

	// Clamp offset to a valid window based on current height
	if m.resultsOffset < 0 {
		m.resultsOffset = 0
	}
	// maxVisible computed above
	maxOffset := 0
	if len(m.results) > maxVisible {
		maxOffset = len(m.results) - maxVisible
	}
	if m.resultsOffset > maxOffset {
		m.resultsOffset = maxOffset
	}

	start := m.resultsOffset
	end := start + maxVisible
	if end > len(m.results) {
		end = len(m.results)
	}

	// Display results
	for i := start; i < end; i++ {
		r := m.results[i]

		// Build line: checkbox + filename + stars
		checkbox := "[ ]"
		if m.globalSelections.IsSelected(r.Repo, r.Path) {
			checkbox = "[✓]"
		}

		line := fmt.Sprintf("%s %s", checkbox, r.RelPath)
		if r.Stars > 0 {
			line += fmt.Sprintf(" ⭐ %d", r.Stars)
		}

		// Render with selection
		if i == m.cursor {
			b.WriteString(selectedStyle.Render("> " + line))
		} else {
			b.WriteString(normalStyle.Render("  " + line))
		}
		b.WriteString("\n")
	}

	// Selection count
	if count := m.globalSelections.Count(); count > 0 {
		b.WriteString(fmt.Sprintf("\n%d selected\n", count))
	}

	// Help
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("↑↓ move • space select • enter download • v repo • p preview • esc back"))

	return b.String()
}

func (m model) viewPreview() string {
	if m.previewContent == "" {
		return "Loading preview..."
	}

	var b strings.Builder
	title := titleStyle.Render("📄 Preview")
	b.WriteString(title + "\n\n")
	b.WriteString(m.viewport.View())
	b.WriteString("\n" + helpStyle.Render("↑↓/PgUp/PgDn: scroll • esc/q: back"))
	return b.String()
}

func (m model) viewLocation() string {
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("📦 Download Files"))
	b.WriteString("\n\n")

	// Static file info - just count
	count := m.globalSelections.Count()
	b.WriteString(fmt.Sprintf("Ready to pull down %d files\n\n", count))

	// Location options
	b.WriteString("Select Download Location:\n\n")

	// Option 1
	if m.locationChoice == 0 {
		b.WriteString(selectedStyle.Render("> Global: " + locationPaths[locationGlobal][m.searchMode]()))
	} else {
		b.WriteString(normalStyle.Render("  Global: " + locationPaths[locationGlobal][m.searchMode]()))
	}
	b.WriteString("\n")

	// Option 2
	if m.locationChoice == 1 {
		b.WriteString(selectedStyle.Render("> Current: " + locationPaths[locationCurrent][m.searchMode]()))
	} else {
		b.WriteString(normalStyle.Render("  Current: " + locationPaths[locationCurrent][m.searchMode]()))
	}
	b.WriteString("\n")

	// Option 3
	if m.locationChoice == 2 {
		b.WriteString(selectedStyle.Render("> Custom path..."))
	} else {
		b.WriteString(normalStyle.Render("  Custom path..."))
	}
	b.WriteString("\n")

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("↑/↓ select • v view files • enter download • esc cancel"))

	return b.String()
}

func (m model) viewLocationFileList() string {
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("📋 Files to Download"))
	b.WriteString("\n\n")

	selections := m.globalSelections.GetAll()

	// Show total count
	b.WriteString(fmt.Sprintf("Total: %d files\n\n", len(selections)))

	// Just list all files statically
	for i, sel := range selections {
		b.WriteString(fmt.Sprintf("%d. %s\n", i+1, sel.FileName))
		b.WriteString(fmt.Sprintf("   From: %s\n", sel.Repo))
		if i >= 20 {
			remaining := len(selections) - 20
			b.WriteString(fmt.Sprintf("\n... and %d more files\n", remaining))
			break
		}
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Press v or esc to go back"))

	return b.String()
}

func (m model) viewCustomPath() string {
	title := titleStyle.Render("📁 Enter Custom Path")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		m.customPathInput.View(),
		"",
		helpStyle.Render("Enter: confirm • Esc: back"),
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func (m model) viewDownloading() string {
	title := titleStyle.Render("⬇️ Downloading Files...")
	loadingText := "Saving agent files..."

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		loadingText,
		"",
		helpStyle.Render("Please wait..."),
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m model) viewComplete() string {
	var path string
	switch m.location {
	case locationGlobal:
		path = locationPaths[locationGlobal][m.searchMode]()
	case locationCurrent:
		path = locationPaths[locationCurrent][m.searchMode]()
	case locationCustom:
		path = m.customPath
	}

	// Calculate total files that should have been downloaded
	totalFiles := m.globalSelections.Count()

	title := successStyle.Render("✅ Download Complete!")
	details := fmt.Sprintf("%d files saved to:\n%s", totalFiles, path)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		normalStyle.Render(details),
		"",
		helpStyle.Render("Press Enter to search again • q to quit"),
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m model) viewRepoViewer() string {
	if m.repoViewer == nil {
		return "Loading repository viewer..."
	}
	return m.repoViewer.View()
}

func (m model) viewConfirmLoseSelections() string {
	count := m.globalSelections.Count()

	// Create warning box
	warningBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.error).
		Padding(2, 4).
		Width(60).
		Align(lipgloss.Center)

	warningTitle := lipgloss.NewStyle().
		Foreground(theme.error).
		Bold(true).
		Render("⚠️  Hold Up!")

	warningText := fmt.Sprintf("You've got %d files selected that will be lost\nif you go back to search.", count)

	// Create button-like options
	options := []string{
		"idgaf",
		"oh shit, go back",
	}

	var optionsList strings.Builder
	optionsList.WriteString("\n")
	for i, opt := range options {
		if i == m.confirmChoice {
			// Selected option
			btn := lipgloss.NewStyle().
				Background(theme.primary).
				Foreground(lipgloss.Color("#FFFFFF")).
				Padding(0, 2).
				Bold(true).
				Render(opt)
			optionsList.WriteString(fmt.Sprintf("  ▶ %s\n", btn))
		} else {
			// Unselected option
			btn := lipgloss.NewStyle().
				Foreground(theme.muted).
				Padding(0, 2).
				Render(opt)
			optionsList.WriteString(fmt.Sprintf("    %s\n", btn))
		}
	}

	help := helpStyle.Render("↑/↓: select • enter: confirm • esc: cancel")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		warningTitle,
		"",
		warningText,
		optionsList.String(),
		"",
		help,
	)

	boxed := warningBox.Render(content)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		boxed,
	)
}

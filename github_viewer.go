package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Add missing styles that are used in the View method
var (
	dimStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262"))
	
	subtitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4"))
)

// RepoViewer is a simple GitHub repository file browser
type RepoViewer struct {
	repo        string
	stars       int
	path        string
	items       []repoItem
	cursor      int
	viewport    viewport.Model
	viewingFile bool
	fileName   string
	err         error
	selections  *SelectionManager // Global selection manager
	width       int
	height      int
}

type repoItem struct {
	Name string
	Path string
	Type string // "file" or "dir"
}

// Messages for repo viewer
type repoContentsMsg struct {
	items []repoItem
	err   error
}

type repoFileMsg struct {
	content string
	err     error
}

type backToResultsMsg struct{}

// NewRepoViewer creates a new repository viewer
func NewRepoViewer(repo string, stars int, selections *SelectionManager) RepoViewer {
	vp := viewport.New(80, 20)
	return RepoViewer{
		repo:       repo,
		stars:      stars,
		path:       "",
		viewport:   vp,
		selections: selections,
		width:      80,
		height:     24,
	}
}

func (r RepoViewer) Init() tea.Cmd {
	return r.loadContents("")
}

func (r RepoViewer) Update(msg tea.Msg) (RepoViewer, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		r.width = msg.Width
		r.height = msg.Height
		r.viewport.Width = msg.Width
		r.viewport.Height = msg.Height - 5 // Leave room for header/footer
		
	case repoContentsMsg:
		if msg.err != nil {
			r.err = msg.err
			return r, nil
		}
		r.items = msg.items
		r.cursor = 0
		r.viewingFile = false
		
	case repoFileMsg:
		if msg.err != nil {
			r.err = msg.err
			return r, nil
		}
		r.viewport.SetContent(msg.content)
		r.viewport.GotoTop()
		r.viewingFile = true
		
	case tea.KeyMsg:
		if r.viewingFile {
			// File viewing mode
			switch msg.String() {
			case "esc", "backspace":
				r.viewingFile = false
				r.viewport.SetContent("")
				return r, nil
			case "q":
				// q quits the entire viewer, not just the file view
				return r, func() tea.Msg { return backToResultsMsg{} }
			}
			// Let viewport handle scrolling
			var cmd tea.Cmd
			r.viewport, cmd = r.viewport.Update(msg)
			return r, cmd
		} else {
			// Directory browsing mode
			switch msg.String() {
			case "q", "esc":
				return r, func() tea.Msg { return backToResultsMsg{} }
				
			case "up", "k":
				if r.cursor > 0 {
					r.cursor--
				}
				
			case "down", "j":
				if r.cursor < len(r.items)-1 {
					r.cursor++
				}
				
			case " ", "space":
				// Toggle selection on .md files only
				if r.cursor < len(r.items) {
					item := r.items[r.cursor]
					if item.Type == "file" && strings.HasSuffix(item.Name, ".md") {
						// Build the full GitHub URL for this file
						url := fmt.Sprintf("https://github.com/%s/blob/main/%s", r.repo, item.Path)
						
						sel := GlobalSelection{
							Repo:     r.repo,
							Path:     item.Path,
							URL:      url,
							FileName: item.Name,
							Source:   "repo",
						}
						r.selections.Toggle(sel)
					}
				}
				
			case "enter":
				if r.cursor < len(r.items) {
					item := r.items[r.cursor]
					if item.Type == "dir" {
						r.path = item.Path
						return r, r.loadContents(item.Path)
					} else {
						r.fileName = item.Name
						return r, r.loadFile(item.Path)
					}
				}
				
			case "backspace":
				// Go up a directory
				if r.path != "" {
					parts := strings.Split(r.path, "/")
					if len(parts) > 1 {
						r.path = strings.Join(parts[:len(parts)-1], "/")
					} else {
						r.path = ""
					}
					return r, r.loadContents(r.path)
				}
			}
		}
	}
	
	return r, nil
}

func (r RepoViewer) View() string {
	if r.err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v", r.err))
	}
	
	// Header
	repoName := r.repo
	// Truncate repo name if too long
	if len(repoName) > 30 {
		repoName = "..." + repoName[len(repoName)-27:]
	}
	header := titleStyle.Render(fmt.Sprintf("📂 %s ⭐ %d", repoName, r.stars))
	if r.path != "" {
		pathDisplay := r.path
		// Truncate path if too long
		if len(pathDisplay) > 40 {
			pathDisplay = "..." + pathDisplay[len(pathDisplay)-37:]
		}
		header += dimStyle.Render(fmt.Sprintf(" | %s", pathDisplay))
	}
	
	// Content
	var content string
	if r.viewingFile {
		// Show file content in viewport
		content = subtitleStyle.Render(fmt.Sprintf("📄 %s", r.fileName)) + "\n"
		content += r.viewport.View()
		content += "\n" + dimStyle.Render(fmt.Sprintf("%3.f%%", r.viewport.ScrollPercent()*100))
		content += "\n" + helpStyle.Render("↑/↓: scroll • esc: close file • q: quit app")
	} else {
		// Show directory listing with viewport scrolling
		// Dynamic max visible based on terminal height
		// Reserve space for: header (2), help (2), selections (2), padding (2)
		maxVisible := r.height - 8
		if maxVisible < 5 {
			maxVisible = 5 // Minimum visible items
		}
		if maxVisible > 30 {
			maxVisible = 30 // Cap at reasonable max
		}
		
		// Calculate visible range, keeping cursor in view
		visibleStart := r.cursor - maxVisible/2
		if visibleStart < 0 {
			visibleStart = 0
		}
		visibleEnd := visibleStart + maxVisible
		if visibleEnd > len(r.items) {
			visibleEnd = len(r.items)
			visibleStart = visibleEnd - maxVisible
			if visibleStart < 0 {
				visibleStart = 0
			}
		}
		
		var items []string
		for i := visibleStart; i < visibleEnd && i < len(r.items); i++ {
			item := r.items[i]
			icon := "📄"
			name := item.Name
			if item.Type == "dir" {
				icon = "📁"
				name += "/"
			}
			
			// Truncate long names to prevent overflow
			maxNameLen := 40
			if len(name) > maxNameLen {
				name = name[:maxNameLen-3] + "..."
			}
			
			// Add selection checkbox for .md files
			checkbox := "  "
			if item.Type == "file" && strings.HasSuffix(item.Name, ".md") {
				if r.selections.IsSelected(r.repo, item.Path) {
					checkbox = "[x]"
				} else {
					checkbox = "[ ]"
				}
			}
			
			line := fmt.Sprintf("%s %s %s", checkbox, icon, name)
			if i == r.cursor {
				items = append(items, selectedStyle.Render("> "+line))
			} else {
				items = append(items, normalStyle.Render("  "+line))
			}
		}
		
		if len(r.items) == 0 {
			items = append(items, dimStyle.Render("  (empty)"))
		}
		
		// Add scroll indicator if needed
		if len(r.items) > maxVisible {
			scrollInfo := fmt.Sprintf("\n[%d/%d]", r.cursor+1, len(r.items))
			items = append(items, dimStyle.Render(scrollInfo))
		}
		
		content = strings.Join(items, "\n")
		
		// Show selection count and help
		repoSelCount := len(r.selections.GetRepoSelections(r.repo))
		totalSelCount := r.selections.Count()
		
		helpText := "↑/↓: navigate • enter: open • space: select • backspace: up • q: quit"
		if totalSelCount > 0 {
			selInfo := fmt.Sprintf("%d selected", totalSelCount)
			if repoSelCount > 0 && repoSelCount != totalSelCount {
				selInfo = fmt.Sprintf("%d in repo, %d total", repoSelCount, totalSelCount)
			}
			content += "\n\n" + dimStyle.Render(selInfo)
		}
		content += "\n" + helpStyle.Render(helpText)
	}
	
	return lipgloss.JoinVertical(lipgloss.Left, header, "", content)
}

func (r RepoViewer) loadContents(path string) tea.Cmd {
	return func() tea.Msg {
		var url string
		if path == "" {
			url = fmt.Sprintf("repos/%s/contents", r.repo)
		} else {
			url = fmt.Sprintf("repos/%s/contents/%s", r.repo, path)
		}
		
		cmd := exec.Command("gh", "api", url, "--paginate")
		output, err := cmd.Output()
		if err != nil {
			return repoContentsMsg{err: fmt.Errorf("failed to load directory")}
		}
		
		var ghItems []struct {
			Name string `json:"name"`
			Path string `json:"path"`
			Type string `json:"type"`
		}
		
		if err := json.Unmarshal(output, &ghItems); err != nil {
			return repoContentsMsg{err: err}
		}
		
		items := make([]repoItem, 0, len(ghItems))
		for _, item := range ghItems {
			items = append(items, repoItem{
				Name: item.Name,
				Path: item.Path,
				Type: item.Type,
			})
		}
		
		return repoContentsMsg{items: items}
	}
}

func (r RepoViewer) loadFile(path string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("gh", "api",
			fmt.Sprintf("repos/%s/contents/%s", r.repo, path),
			"-H", "Accept: application/vnd.github.v3.raw")
		
		output, err := cmd.Output()
		if err != nil {
			return repoFileMsg{err: fmt.Errorf("failed to load file")}
		}
		
		return repoFileMsg{content: string(output)}
	}
}
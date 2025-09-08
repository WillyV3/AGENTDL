package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// GlobalSelection represents a file selected from any source (search results or repo browser)
type GlobalSelection struct {
	Repo     string // Repository name (e.g., "owner/repo")
	Path     string // Full path in repo
	URL      string // GitHub URL for downloading
	FileName string // Just the filename
	Source   string // "search" or "repo"
}

// SelectionManager manages the global list of selected files
type SelectionManager struct {
	selections map[string]*GlobalSelection // key is "repo:path" for uniqueness
}

// NewSelectionManager creates a new selection manager
func NewSelectionManager() *SelectionManager {
	return &SelectionManager{
		selections: make(map[string]*GlobalSelection),
	}
}

// Add adds a selection to the global list
func (sm *SelectionManager) Add(sel GlobalSelection) {
	key := sm.makeKey(sel.Repo, sel.Path)
	sm.selections[key] = &sel
}

// Remove removes a selection from the global list
func (sm *SelectionManager) Remove(repo, path string) {
	key := sm.makeKey(repo, path)
	delete(sm.selections, key)
}

// Toggle toggles a selection (add if not present, remove if present)
func (sm *SelectionManager) Toggle(sel GlobalSelection) bool {
	key := sm.makeKey(sel.Repo, sel.Path)
	if _, exists := sm.selections[key]; exists {
		delete(sm.selections, key)
		return false // now unselected
	}
	sm.selections[key] = &sel
	return true // now selected
}

// IsSelected checks if a file is selected
func (sm *SelectionManager) IsSelected(repo, path string) bool {
	key := sm.makeKey(repo, path)
	_, exists := sm.selections[key]
	return exists
}

// Clear removes all selections
func (sm *SelectionManager) Clear() {
	sm.selections = make(map[string]*GlobalSelection)
}

// GetAll returns all selections as a slice, sorted by repo then path
func (sm *SelectionManager) GetAll() []GlobalSelection {
	result := make([]GlobalSelection, 0, len(sm.selections))
	for _, sel := range sm.selections {
		result = append(result, *sel)
	}
	// Sort to ensure consistent order
	sort.Slice(result, func(i, j int) bool {
		if result[i].Repo != result[j].Repo {
			return result[i].Repo < result[j].Repo
		}
		return result[i].Path < result[j].Path
	})
	return result
}

// GetRepoSelections returns selections for a specific repo
func (sm *SelectionManager) GetRepoSelections(repo string) []GlobalSelection {
	var result []GlobalSelection
	for _, sel := range sm.selections {
		if sel.Repo == repo && sel.Source == "repo" {
			result = append(result, *sel)
		}
	}
	return result
}

// Count returns the total number of selections
func (sm *SelectionManager) Count() int {
	return len(sm.selections)
}

// CountBySource returns count of selections by source type
func (sm *SelectionManager) CountBySource(source string) int {
	count := 0
	for _, sel := range sm.selections {
		if sel.Source == source {
			count++
		}
	}
	return count
}

// makeKey creates a unique key for a selection
func (sm *SelectionManager) makeKey(repo, path string) string {
	return fmt.Sprintf("%s:%s", repo, path)
}


// GetDownloadURL converts a GitHub file URL to raw download URL
func GetDownloadURL(url string) string {
	// Convert from: https://github.com/owner/repo/blob/main/path/file.md
	// To: https://raw.githubusercontent.com/owner/repo/main/path/file.md
	url = strings.Replace(url, "github.com", "raw.githubusercontent.com", 1)
	url = strings.Replace(url, "/blob/", "/", 1)
	return url
}

// ExtractFileName gets just the filename from a path
func ExtractFileName(path string) string {
	return filepath.Base(path)
}
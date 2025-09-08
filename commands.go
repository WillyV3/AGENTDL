package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"agent-search/github"
	tea "github.com/charmbracelet/bubbletea"
)

// ============================
// Commands (Async Operations)
// ============================

func searchGitHub(query string, mode string) tea.Cmd {
	return func() tea.Msg {
		opts := github.SearchOptions{
			MatchMode: mode,
			Limit:     200,
		}
		
		githubResults := github.Search(query, opts)
		
		// Convert to our internal type
		results := make([]searchResult, len(githubResults))
		for i, r := range githubResults {
			results[i] = searchResult(r)
		}
		
		return searchResultsMsg{results: results}
	}
}

func fetchFileContent(url string) tea.Cmd {
	return func() tea.Msg {
		// Convert to raw URL
		rawURL := strings.Replace(url, "/blob/", "/raw/", 1)
		
		content, err := downloadFile(rawURL)
		if err != nil {
			return fileContentMsg{err: err}
		}
		
		// Limit preview to first 100 lines for performance
		lines := strings.Split(content, "\n")
		if len(lines) > 100 {
			lines = lines[:100]
			lines = append(lines, "", "... (preview truncated at 100 lines)")
		}
		content = strings.Join(lines, "\n")
		
		return fileContentMsg{content: content, err: err}
	}
}

func downloadSelectedFiles(results []searchResult, location string, selections *SelectionManager) tea.Cmd {
	return func() tea.Msg {
		count := 0
		
		// Download all selections from the global selection manager
		if selections != nil {
			for _, sel := range selections.GetAll() {
				url := sel.URL
				// Convert GitHub blob URL to raw URL
				url = strings.Replace(url, "github.com", "raw.githubusercontent.com", 1)
				url = strings.Replace(url, "/blob/", "/", 1)
				
				content, err := downloadFile(url)
				if err != nil {
					continue
				}
				
				destPath := filepath.Join(location, sel.FileName)
				os.MkdirAll(filepath.Dir(destPath), 0755)
				os.WriteFile(destPath, []byte(content), 0644)
				count++
			}
		}
		
		return downloadCompleteMsg{count: count}
	}
}

// Helper function to download a file from URL
func downloadFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return string(content), nil
}
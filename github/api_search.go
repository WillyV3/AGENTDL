package github

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// APISearchByFilename uses GitHub CLI content search then filters by filename
func APISearchByFilename(keywords string, opts SearchOptions) []Result {
	if opts.Limit == 0 {
		opts.Limit = 300
	}

	// Use content search to find files, then filter by filename locally
	searchQuery := buildContentSearchQuery(keywords, opts)

	// Get a larger set of results to filter from
	searchLimit := opts.Limit * 3 // Search 3x more to account for filtering
	if searchLimit > 1000 {
		searchLimit = 1000 // GitHub API limit
	}

	results := searchWithContentQuery(searchQuery, searchLimit, opts)

	// Filter results by filename matching
	filtered := filterByFilename(results, keywords, opts)

	// Limit final results
	if len(filtered) > opts.Limit {
		filtered = filtered[:opts.Limit]
	}

	// Fetch stars for final results
	if len(filtered) > 0 {
		repos := extractUniqueRepos(filtered)
		stars := fetchStarsParallel(repos)

		// Update star counts
		for i := range filtered {
			filtered[i].Stars = stars[filtered[i].Repo]
		}
	}

	return filtered
}

// buildContentSearchQuery creates a GitHub content search query
func buildContentSearchQuery(keywords string, opts SearchOptions) string {
	basePath := buildBasePath(opts)

	if keywords == "" {
		return basePath
	}

	// Simple approach: just search for the keywords in content with path restriction
	return fmt.Sprintf("%s %s", keywords, basePath)
}

// buildBasePath creates the base path restriction for the search
func buildBasePath(opts SearchOptions) string {
	if opts.SearchMode == ModeCommands {
		return "path:/.claude/commands/"
	}
	return "path:/.claude/agents/"
}

// searchWithContentQuery executes GitHub content search
func searchWithContentQuery(query string, limit int, opts SearchOptions) []Result {
	// Use GitHub CLI to search
	cmd := exec.Command("gh", "search", "code", query,
		"--limit", fmt.Sprintf("%d", limit),
		"--json", "repository,path,url")

	output, err := cmd.Output()
	if err != nil {
		log.Printf("GitHub search failed for query '%s': %v", query, err)
		return []Result{}
	}

	// Parse results
	var rawResults []struct {
		Repository struct {
			NameWithOwner string `json:"nameWithOwner"`
		} `json:"repository"`
		Path string `json:"path"`
		URL  string `json:"url"`
	}

	if err := json.Unmarshal(output, &rawResults); err != nil {
		log.Printf("JSON parse error for query '%s': %v", query, err)
		return []Result{}
	}

	// Convert to Results
	var results []Result
	for _, r := range rawResults {
		// Only include .md files
		if !strings.HasSuffix(r.Path, ".md") {
			continue
		}

		// Build relative path
		relPath := r.Repository.NameWithOwner + "/" + r.Path
		if idx := strings.Index(r.Path, ".claude/agents/"); idx >= 0 {
			relPath = r.Repository.NameWithOwner + "/" + r.Path[idx+15:]
		} else if idx := strings.Index(r.Path, ".claude/commands/"); idx >= 0 {
			relPath = r.Repository.NameWithOwner + "/" + r.Path[idx+16:]
		}

		results = append(results, Result{
			Repo:     r.Repository.NameWithOwner,
			Path:     r.Path,
			URL:      r.URL,
			Stars:    0, // Will be filled later
			RelPath:  relPath,
			Selected: false,
		})
	}

	return results
}

// filterByFilename filters results to only include files with keywords in filename
func filterByFilename(results []Result, keywords string, opts SearchOptions) []Result {
	if keywords == "" {
		return results
	}

	keywordList := strings.Fields(strings.ToLower(keywords))
	var filtered []Result

	for _, result := range results {
		// Extract filename from path
		parts := strings.Split(result.Path, "/")
		filename := strings.ToLower(parts[len(parts)-1])

		// Check if filename matches keywords
		if filenameMatches(filename, keywordList, opts.MatchMode) {
			filtered = append(filtered, result)
		}
	}

	return filtered
}

// filenameMatches checks if a filename contains the required keywords
func filenameMatches(filename string, keywords []string, matchMode string) bool {
	if len(keywords) == 0 {
		return true
	}

	if matchMode == "any" {
		// OR mode: filename must contain at least one keyword
		for _, keyword := range keywords {
			if strings.Contains(filename, keyword) {
				return true
			}
		}
		return false
	}

	// AND mode: filename must contain all keywords
	for _, keyword := range keywords {
		if !strings.Contains(filename, keyword) {
			return false
		}
	}
	return true
}

// deduplicateAndLimit removes duplicate results and applies limit
func deduplicateAndLimit(results []Result, limit int) []Result {
	seen := make(map[string]bool)
	var unique []Result

	for _, r := range results {
		key := r.Repo + ":" + r.Path
		if !seen[key] {
			seen[key] = true
			unique = append(unique, r)

			if len(unique) >= limit {
				break
			}
		}
	}

	// Fetch stars for unique results
	if len(unique) > 0 {
		repos := extractUniqueRepos(unique)
		stars := fetchStarsParallel(repos)

		// Update star counts
		for i := range unique {
			unique[i].Stars = stars[unique[i].Repo]
		}
	}

	return unique
}

// extractUniqueRepos gets unique repository names from results
func extractUniqueRepos(results []Result) []string {
	seen := make(map[string]bool)
	var repos []string

	for _, r := range results {
		if !seen[r.Repo] {
			seen[r.Repo] = true
			repos = append(repos, r.Repo)
		}
	}

	return repos
}
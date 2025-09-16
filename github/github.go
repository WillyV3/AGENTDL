package github

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// Result represents a search result from GitHub
type Result struct {
	Repo     string
	Path     string
	URL      string
	Stars    int
	RelPath  string
	Selected bool
}

// SearchMode represents the search target
type SearchMode int

const (
	ModeAgents SearchMode = iota
	ModeCommands
)

// SearchOptions configures the search behavior
type SearchOptions struct {
	MatchMode  string     // "all" (AND), "any" (OR)
	SearchMode SearchMode // agents or commands
	Limit      int
}


// Search performs intelligent keyword search on GitHub
func Search(query string, opts SearchOptions) []Result {
	// Use paginated search with rate limiting for filename-only search
	if query != "" {
		return PaginatedSearchByFilename(query, opts)
	}
	// Fall back to original approach for empty queries
	return searchFallback(query, opts)
}

// searchFallback is the original gh CLI implementation
func searchFallback(query string, opts SearchOptions) []Result {
	if opts.Limit == 0 {
		opts.Limit = 300
	}

	// Build the search query
	searchQuery := buildQuery(query, opts)

	// Execute GitHub search
	cmd := exec.Command("gh", "search", "code", searchQuery,
		"--limit", fmt.Sprintf("%d", opts.Limit),
		"--json", "repository,path,url")

	output, err := cmd.Output()
	if err != nil {
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
		return []Result{}
	}

	// Filter results to only .md files and apply filename matching if query provided
	filteredResults := []struct {
		Repository struct {
			NameWithOwner string `json:"nameWithOwner"`
		} `json:"repository"`
		Path string `json:"path"`
		URL  string `json:"url"`
	}{}

	for _, r := range rawResults {
		// Only include .md files
		if !strings.HasSuffix(r.Path, ".md") {
			continue
		}

		filteredResults = append(filteredResults, r)
	}

	// Fetch stars in parallel
	stars := fetchStarsParallel(extractFilteredRepos(filteredResults))

	// Build final results
	results := make([]Result, 0, len(filteredResults))
	for _, r := range filteredResults {
		relPath := r.Repository.NameWithOwner + "/" + r.Path
		if idx := strings.Index(r.Path, ".claude/agents/"); idx >= 0 {
			relPath = r.Repository.NameWithOwner + "/" + r.Path[idx+15:]
		} else if idx := strings.Index(r.Path, ".claude/commands/"); idx >= 0 {
			relPath = r.Repository.NameWithOwner + "/" + r.Path[idx+16:]
		}
		results = append(results, Result{
			Repo:    r.Repository.NameWithOwner,
			Path:    r.Path,
			URL:     r.URL,
			Stars:   stars[r.Repository.NameWithOwner],
			RelPath: relPath,
		})
	}

	// Sort by stars
	sort.Slice(results, func(i, j int) bool {
		return results[i].Stars > results[j].Stars
	})

	return results
}

// buildQuery constructs GitHub search query from user input
func buildQuery(input string, opts SearchOptions) string {
	keywords := strings.Fields(strings.TrimSpace(input))

	// Determine path based on search mode
	var pathQuery string
	if opts.SearchMode == ModeCommands {
		pathQuery = "path:/.claude/commands/"
	} else {
		pathQuery = "path:/.claude/agents/"
	}

	if len(keywords) == 0 {
		return pathQuery
	}

	// Build filename-specific search by combining keywords with wildcard patterns
	// This searches for files that contain the keywords in their actual path/filename
	var searchTerms []string
	for _, keyword := range keywords {
		// Search for files with keyword in the path (including filename)
		searchTerms = append(searchTerms, keyword)
	}

	if opts.MatchMode == "any" {
		// OR mode: (keyword1 OR keyword2) path
		return "(" + strings.Join(searchTerms, " OR ") + ") " + pathQuery
	}

	// Default AND mode: keyword1 keyword2 path (AND is implicit)
	return strings.Join(searchTerms, " ") + " " + pathQuery
}

// Helper functions
func extractRepos(results []struct {
	Repository struct{ NameWithOwner string `json:"nameWithOwner"` } `json:"repository"`
	Path string `json:"path"`
	URL  string `json:"url"`
}) []string {
	seen := make(map[string]bool)
	repos := []string{}
	for _, r := range results {
		if !seen[r.Repository.NameWithOwner] {
			seen[r.Repository.NameWithOwner] = true
			repos = append(repos, r.Repository.NameWithOwner)
		}
	}
	return repos
}

func extractFilteredRepos(results []struct {
	Repository struct{ NameWithOwner string `json:"nameWithOwner"` } `json:"repository"`
	Path string `json:"path"`
	URL  string `json:"url"`
}) []string {
	seen := make(map[string]bool)
	repos := []string{}
	for _, r := range results {
		if !seen[r.Repository.NameWithOwner] {
			seen[r.Repository.NameWithOwner] = true
			repos = append(repos, r.Repository.NameWithOwner)
		}
	}
	return repos
}

func fetchStarsParallel(repos []string) map[string]int {
	stars := make(map[string]int)
	var mu sync.Mutex
	var wg sync.WaitGroup
	
	// Limit concurrent requests
	sem := make(chan struct{}, 5)
	
	for _, repo := range repos {
		wg.Add(1)
		go func(r string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			
			cmd := exec.Command("gh", "api", fmt.Sprintf("repos/%s", r), "--jq", ".stargazers_count")
			if output, err := cmd.Output(); err == nil {
				if count, err := strconv.Atoi(strings.TrimSpace(string(output))); err == nil {
					mu.Lock()
					stars[r] = count
					mu.Unlock()
				}
			}
		}(repo)
	}
	
	wg.Wait()
	return stars
}
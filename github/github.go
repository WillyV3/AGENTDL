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
	if opts.Limit == 0 {
		opts.Limit = 200
	}
	
	// Build the search query
	searchQuery := buildQuery(query, opts)
	
	// Execute GitHub search
	cmd := exec.Command("gh", "search", "code", searchQuery,
		"--match", "path",
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
	
	// Fetch stars in parallel
	stars := fetchStarsParallel(extractRepos(rawResults))
	
	// Build final results
	results := make([]Result, 0, len(rawResults))
	for _, r := range rawResults {
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
	var path string
	if opts.SearchMode == ModeCommands {
		path = "path:.claude/commands/"
	} else {
		path = "path:.claude/agents/"
	}

	if len(keywords) == 0 {
		return path
	}

	// Handle exact phrase (quoted)
	if strings.HasPrefix(input, "\"") && strings.HasSuffix(input, "\"") {
		return input + " " + path
	}

	// Build keyword query based on mode
	if opts.MatchMode == "any" {
		// OR mode: keyword1 OR keyword2 OR keyword3
		return "(" + strings.Join(keywords, " OR ") + ") " + path
	}

	// Default AND mode: all keywords must match
	return strings.Join(keywords, " ") + " " + path
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
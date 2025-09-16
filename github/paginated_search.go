package github

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

// PaginatedSearchByFilename implements rate-limited search with filename filtering
func PaginatedSearchByFilename(keywords string, opts SearchOptions) []Result {
	if opts.Limit == 0 {
		opts.Limit = 100
	}

	// First try using the filename flag for direct filename search
	if keywords != "" {
		results := searchWithFilenameFlag(keywords, opts)
		if len(results) > 0 {
			return results
		}
	}

	// Fallback to content search with filename filtering
	return searchWithRateLimit(keywords, opts)
}

// searchWithFilenameFlag uses GitHub CLI's --filename flag for direct filename matching
func searchWithFilenameFlag(keywords string, opts SearchOptions) []Result {
	basePath := buildBasePath(opts)

	// Create filename pattern - search for files containing the keyword in filename
	filenamePattern := fmt.Sprintf("*%s*", keywords)

	// Build the search command with filename filtering
	cmd := exec.Command("gh", "search", "code",
		"--filename", filenamePattern,
		"--match", "path",
		basePath,
		"--limit", fmt.Sprintf("%d", opts.Limit),
		"--json", "repository,path,url")

	output, err := cmd.Output()
	if err != nil {
		log.Printf("GitHub filename search failed: %v", err)
		return []Result{}
	}

	return parseSearchResults(output)
}

// searchWithRateLimit implements rate-limited content search with retries
func searchWithRateLimit(keywords string, opts SearchOptions) []Result {
	var allResults []Result
	remainingLimit := opts.Limit
	batchSize := 30 // Conservative batch size to avoid rate limits

	for remainingLimit > 0 {
		currentBatch := batchSize
		if currentBatch > remainingLimit {
			currentBatch = remainingLimit
		}

		// Search for this batch
		query := buildContentSearchQuery(keywords, opts)
		results := searchBatchWithRetry(query, currentBatch)

		if len(results) == 0 {
			break // No more results
		}

		// Filter by filename
		filtered := filterByFilename(results, keywords, opts)
		allResults = append(allResults, filtered...)

		remainingLimit -= len(filtered)

		// Rate limiting: wait between requests
		if remainingLimit > 0 {
			log.Printf("Waiting 2 seconds to avoid rate limit...")
			time.Sleep(2 * time.Second)
		}

		// If we got fewer results than requested, we've reached the end
		if len(results) < currentBatch {
			break
		}
	}

	// Fetch stars for all results
	if len(allResults) > 0 {
		repos := extractUniqueRepos(allResults)
		stars := fetchStarsParallel(repos)

		for i := range allResults {
			allResults[i].Stars = stars[allResults[i].Repo]
		}
	}

	return allResults
}

// searchBatchWithRetry searches a single batch with retry logic for rate limits
func searchBatchWithRetry(query string, limit int) []Result {
	maxRetries := 3
	baseDelay := 10 * time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		cmd := exec.Command("gh", "search", "code", query,
			"--limit", fmt.Sprintf("%d", limit),
			"--json", "repository,path,url")

		output, err := cmd.Output()
		if err != nil {
			if strings.Contains(err.Error(), "rate limit") {
				waitTime := baseDelay * time.Duration(1<<attempt) // Exponential backoff
				log.Printf("Rate limit hit, waiting %v before retry %d/%d", waitTime, attempt+1, maxRetries)
				time.Sleep(waitTime)
				continue
			}
			log.Printf("Search failed: %v", err)
			return []Result{}
		}

		return parseSearchResults(output)
	}

	log.Printf("Max retries exceeded for query: %s", query)
	return []Result{}
}

// parseSearchResults parses GitHub CLI JSON output into Results
func parseSearchResults(output []byte) []Result {
	var rawResults []struct {
		Repository struct {
			NameWithOwner string `json:"nameWithOwner"`
		} `json:"repository"`
		Path string `json:"path"`
		URL  string `json:"url"`
	}

	if err := json.Unmarshal(output, &rawResults); err != nil {
		log.Printf("JSON parse error: %v", err)
		return []Result{}
	}

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
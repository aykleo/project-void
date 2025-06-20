package git

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type GitHubProvider struct {
	client *http.Client
	token  string
}

type GitHubCommit struct {
	SHA    string `json:"sha"`
	Commit struct {
		Author struct {
			Name string    `json:"name"`
			Date time.Time `json:"date"`
		} `json:"author"`
		Message string `json:"message"`
	} `json:"commit"`
}

type GitHubBranch struct {
	Name   string `json:"name"`
	Commit struct {
		SHA string `json:"sha"`
	} `json:"commit"`
}

func NewGitHubProvider() *GitHubProvider {
	return &GitHubProvider{
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (g *GitHubProvider) SetToken(token string) {
	g.token = token
}

func (g *GitHubProvider) parseGitHubURL(repoURL string) (owner, repo string, err error) {
	patterns := []string{
		`github\.com[:/]([^/]+)/([^/]+?)(?:\.git)?/?$`,
		`github\.com/([^/]+)/([^/]+?)(?:\.git)?/?$`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(repoURL)
		if len(matches) >= 3 {
			return matches[1], matches[2], nil
		}
	}

	return "", "", fmt.Errorf("invalid GitHub URL format: %s", repoURL)
}

func (g *GitHubProvider) makeRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if g.token != "" {
		req.Header.Set("Authorization", "token "+g.token)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}

	return resp, nil
}

func (g *GitHubProvider) GetCommitsSince(repoURL string, since time.Time) ([]Commit, error) {
	owner, repo, err := g.parseGitHubURL(repoURL)
	if err != nil {
		return nil, err
	}

	branches, err := g.getBranches(owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get branches: %w", err)
	}

	uniqueCommits := make(map[string]Commit)

	for _, branch := range branches {
		commits, err := g.getCommitsFromBranch(owner, repo, branch, since)
		if err != nil {
			continue
		}

		for _, commit := range commits {
			if _, exists := uniqueCommits[commit.Hash]; !exists {
				uniqueCommits[commit.Hash] = commit
			}
		}
	}

	result := make([]Commit, 0, len(uniqueCommits))
	for _, commit := range uniqueCommits {
		result = append(result, commit)
	}

	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].Timestamp.Before(result[j].Timestamp) {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result, nil
}

func (g *GitHubProvider) GetCommitsSinceByAuthors(repoURL string, since time.Time, authorNames []string) ([]Commit, error) {
	commits, err := g.GetCommitsSince(repoURL, since)
	if err != nil {
		return nil, err
	}

	var lowerAuthorNames []string
	for _, name := range authorNames {
		lowerAuthorNames = append(lowerAuthorNames, strings.ToLower(name))
	}

	var filtered []Commit
	for _, commit := range commits {
		authorLower := strings.ToLower(commit.Author)
		matchesAuthor := false

		for _, targetAuthor := range lowerAuthorNames {
			if strings.Contains(authorLower, targetAuthor) || strings.Contains(targetAuthor, authorLower) {
				matchesAuthor = true
				break
			}
		}

		if matchesAuthor {
			filtered = append(filtered, commit)
		}
	}

	return filtered, nil
}

func (g *GitHubProvider) GetCommitsSinceByBranches(repoURL string, since time.Time, branchNames []string) ([]Commit, error) {
	commits, err := g.GetCommitsSince(repoURL, since)
	if err != nil {
		return nil, err
	}

	var lowerBranchNames []string
	for _, name := range branchNames {
		lowerBranchNames = append(lowerBranchNames, strings.ToLower(name))
	}

	var filtered []Commit
	for _, commit := range commits {
		branchLower := strings.ToLower(commit.Branch)
		matchesBranch := false

		for _, targetBranch := range lowerBranchNames {
			if strings.Contains(branchLower, targetBranch) {
				matchesBranch = true
				break
			}
		}

		if matchesBranch {
			filtered = append(filtered, commit)
		}
	}

	return filtered, nil
}

func (g *GitHubProvider) GetCommitsSinceByAuthorsAndBranches(repoURL string, since time.Time, authorNames []string, branchNames []string) ([]Commit, error) {
	commits, err := g.GetCommitsSince(repoURL, since)
	if err != nil {
		return nil, err
	}

	var lowerAuthorNames []string
	for _, name := range authorNames {
		lowerAuthorNames = append(lowerAuthorNames, strings.ToLower(name))
	}

	var lowerBranchNames []string
	for _, name := range branchNames {
		lowerBranchNames = append(lowerBranchNames, strings.ToLower(name))
	}

	var filtered []Commit
	for _, commit := range commits {

		branchLower := strings.ToLower(commit.Branch)
		matchesBranch := false
		for _, targetBranch := range lowerBranchNames {
			if strings.Contains(branchLower, targetBranch) {
				matchesBranch = true
				break
			}
		}

		if !matchesBranch {
			continue
		}

		authorLower := strings.ToLower(commit.Author)
		matchesAuthor := false
		for _, targetAuthor := range lowerAuthorNames {
			if strings.Contains(authorLower, targetAuthor) || strings.Contains(targetAuthor, authorLower) {
				matchesAuthor = true
				break
			}
		}

		if matchesAuthor {
			filtered = append(filtered, commit)
		}
	}

	return filtered, nil
}

func (g *GitHubProvider) getBranches(owner, repo string) ([]GitHubBranch, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/branches", owner, repo)

	resp, err := g.makeRequest(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get branches for %s/%s: %w", owner, repo, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error for %s/%s branches: HTTP %d", owner, repo, resp.StatusCode)
	}

	var branches []GitHubBranch
	if err := json.NewDecoder(resp.Body).Decode(&branches); err != nil {
		return nil, fmt.Errorf("failed to decode branches response for %s/%s: %w", owner, repo, err)
	}

	return branches, nil
}

func (g *GitHubProvider) getCommitsFromBranch(owner, repo string, branch GitHubBranch, since time.Time) ([]Commit, error) {
	sinceStr := since.UTC().Format(time.RFC3339)

	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", owner, repo)
	params := url.Values{}
	params.Add("sha", branch.Name)
	params.Add("since", sinceStr)
	params.Add("per_page", "100")

	fullURL := apiURL + "?" + params.Encode()

	resp, err := g.makeRequest(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get commits for %s/%s branch %s: %w", owner, repo, branch.Name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error for %s/%s branch %s commits: HTTP %d", owner, repo, branch.Name, resp.StatusCode)
	}

	var githubCommits []GitHubCommit
	if err := json.NewDecoder(resp.Body).Decode(&githubCommits); err != nil {
		return nil, fmt.Errorf("failed to decode commits response for %s/%s branch %s: %w", owner, repo, branch.Name, err)
	}

	var commits []Commit
	for _, gc := range githubCommits {
		commits = append(commits, Commit{
			Hash:      gc.SHA,
			Branch:    branch.Name,
			Author:    gc.Commit.Author.Name,
			Message:   gc.Commit.Message,
			Timestamp: gc.Commit.Author.Date,
		})
	}

	return commits, nil
}

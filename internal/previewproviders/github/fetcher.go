package github

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// RepoInfo contains all the information fetched from a GitHub repository
type RepoInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	LicenseInfo struct {
		Name string `json:"name"`
	} `json:"licenseInfo"`
	UpdatedAt time.Time `json:"updatedAt"`
	Owner     struct {
		Login string `json:"login"`
	} `json:"owner"`
	StargazerCount  int  `json:"stargazerCount"`
	ForkCount       int  `json:"forkCount"`
	IsArchived      bool `json:"isArchived"`
	IsPrivate       bool `json:"isPrivate"`
	PrimaryLanguage struct {
		Name string `json:"name"`
	} `json:"primaryLanguage"`
	Issues struct {
		TotalCount int `json:"totalCount"`
	} `json:"issues"`
	PullRequests struct {
		TotalCount int `json:"totalCount"`
	} `json:"pullRequests"`
	DefaultBranchRef struct {
		Name string `json:"name"`
	} `json:"defaultBranchRef"`
	DiskUsage     int `json:"diskUsage"`
	LatestRelease *struct {
		Name        string    `json:"name"`
		TagName     string    `json:"tagName"`
		PublishedAt time.Time `json:"publishedAt"`
		URL         string    `json:"url"`
	} `json:"latestRelease"`
}

// FetchRepoInfo gathers information about a GitHub repository given its local path
func FetchRepoInfo(repoPath string) (*RepoInfo, string, error) {
	if _, err := exec.LookPath("git"); err != nil {
		return nil, "git not found in PATH", nil
	}

	if _, err := exec.Command("git", "-C", repoPath, "rev-parse").Output(); err != nil {
		return nil, "Not a git repository", nil
	}

	remoteCmd := exec.Command("git", "-C", repoPath, "remote", "get-url", "origin")
	remoteOutput, err := remoteCmd.Output()
	if err != nil {
		return nil, "No remote named origin found", nil
	}
	repoURL := strings.TrimSpace(string(remoteOutput))

	hostname, repoNWO, err := parseGitURL(repoURL)
	if err != nil {
		return nil, fmt.Sprintf("Not a GitHub repository: %s", repoURL), nil
	}

	if _, err := exec.LookPath("gh"); err != nil {
		return nil, "gh not found in PATH", nil
	}

	repoArg := fmt.Sprintf("%s/%s", hostname, repoNWO)
	repoViewCmd := exec.Command("gh", "repo", "view", repoArg, "--json", "name,description,licenseInfo,updatedAt,owner,stargazerCount,forkCount,isArchived,isPrivate,primaryLanguage,issues,pullRequests,defaultBranchRef,diskUsage,latestRelease")
	repoViewOutput, err := repoViewCmd.CombinedOutput()
	if err != nil {
		return nil, "", fmt.Errorf("error getting repo info from gh: %s", string(repoViewOutput))
	}

	var repoInfo RepoInfo
	if err := json.Unmarshal(repoViewOutput, &repoInfo); err != nil {
		return nil, "", fmt.Errorf("error parsing gh output: %w", err)
	}

	return &repoInfo, "", nil
}

func parseGitURL(url string) (hostname string, nwo string, err error) {
	url = strings.TrimSuffix(url, ".git")

	if after, ok := strings.CutPrefix(url, "https://"); ok {
		url = after
		parts := strings.SplitN(url, "/", 2)
		if len(parts) != 2 {
			return "", "", fmt.Errorf("invalid https git url: %s", url)
		}
		hostname = parts[0]
		nwo = parts[1]
		return hostname, nwo, nil
	}

	if after, ok := strings.CutPrefix(url, "git@"); ok {
		url = after
		parts := strings.SplitN(url, ":", 2)
		if len(parts) != 2 {
			return "", "", fmt.Errorf("invalid ssh git url: %s", url)
		}
		hostname = parts[0]
		nwo = parts[1]
		return hostname, nwo, nil
	}

	return "", "", fmt.Errorf("unsupported git remote URL format: %s", url)
}

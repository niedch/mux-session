package previewproviders

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/niedch/mux-session/internal/dataproviders"
)

// GithubPreviewProvider renders GitHub repository information synchronously
type GithubPreviewProvider struct {
	width int
}

// NewGithubPreviewProvider creates a new GitHub preview provider
func NewGithubPreviewProvider(width int) (*GithubPreviewProvider, error) {
	return &GithubPreviewProvider{
		width: width,
	}, nil
}

// Render synchronously fetches the GitHub repository preview for the given item
func (r *GithubPreviewProvider) Render(item any) (string, error) {
	dpItem, ok := item.(*dataproviders.Item)
	if !ok {
		return "", fmt.Errorf("expected *dataproviders.Item, got %T", item)
	}
	return r.fetch(dpItem)
}

func (r *GithubPreviewProvider) fetch(dpItem *dataproviders.Item) (string, error) {
	if _, err := exec.LookPath("git"); err != nil {
		return "git not found in PATH", nil
	}

	if _, err := exec.Command("git", "-C", dpItem.Path, "rev-parse").Output(); err != nil {
		return "Not a git repository", nil
	}

	remoteCmd := exec.Command("git", "-C", dpItem.Path, "remote", "get-url", "origin")
	remoteOutput, err := remoteCmd.Output()
	if err != nil {
		return "No remote named origin found", nil
	}
	repoURL := strings.TrimSpace(string(remoteOutput))

	hostname, repoNWO, err := parseGitURL(repoURL)
	if err != nil {
		return fmt.Sprintf("Not a GitHub repository: %s", repoURL), nil
	}

	if _, err := exec.LookPath("gh"); err != nil {
		return "gh not found in PATH", nil
	}

	repoArg := fmt.Sprintf("%s/%s", hostname, repoNWO)
	repoViewCmd := exec.Command("gh", "repo", "view", repoArg, "--json", "name,description,licenseInfo,pushedAt,owner")
	repoViewOutput, err := repoViewCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error getting repo info from gh: %s", string(repoViewOutput))
	}

	var repoInfo struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		LicenseInfo struct {
			Name string `json:"name"`
		} `json:"licenseInfo"`
		PushedAt time.Time `json:"pushedAt"`
		Owner    struct {
			Login string `json:"login"`
		} `json:"owner"`
	}

	if err := json.Unmarshal(repoViewOutput, &repoInfo); err != nil {
		return "", fmt.Errorf("error parsing gh output: %w", err)
	}

	return fmt.Sprintf(
		"Owner: %s\nName: %s\nDescription: %s\nLicense: %s\nPushed: %s",
		repoInfo.Owner.Login,
		repoInfo.Name,
		repoInfo.Description,
		repoInfo.LicenseInfo.Name,
		repoInfo.PushedAt.Format(time.RFC3339),
	), nil
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

// Name returns the identifier name of this provider
func (r *GithubPreviewProvider) Name() string {
	return "github"
}

// SetWidth updates the renderer with a new width for word wrapping
func (r *GithubPreviewProvider) SetWidth(width int) error {
	r.width = width
	return nil
}

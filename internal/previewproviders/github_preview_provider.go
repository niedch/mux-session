package previewproviders

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/niedch/mux-session/internal/dataproviders"
)

var (
	// Colors
	colorAccent  = lipgloss.AdaptiveColor{Light: "#044289", Dark: "#58A6FF"}
	colorSubtle  = lipgloss.AdaptiveColor{Light: "#57606A", Dark: "#8B949E"}
	colorWarning = lipgloss.AdaptiveColor{Light: "#9A6700", Dark: "#D29922"}
	colorDanger  = lipgloss.AdaptiveColor{Light: "#CF222E", Dark: "#F85149"}
	colorSuccess = lipgloss.AdaptiveColor{Light: "#1A7F37", Dark: "#3FB950"}
	colorDone    = lipgloss.AdaptiveColor{Light: "#8250DF", Dark: "#A371F7"}

	// Styles
	titleStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)

	descStyle = lipgloss.NewStyle().
			Foreground(colorSubtle).
			Italic(true).
			MarginTop(1).
			MarginBottom(1)

	badgePrivateStyle = lipgloss.NewStyle().
				Foreground(colorSubtle).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorSubtle).
				Padding(0, 1).
				MarginLeft(1)

	badgeArchivedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#000000")).
				Background(colorWarning).
				Padding(0, 1).
				MarginLeft(1)

	statItemStyle = lipgloss.NewStyle().
			MarginRight(3)

	statIconStyle = lipgloss.NewStyle().
			MarginRight(1)

	starIconStyle   = statIconStyle.Copy().Foreground(colorWarning)
	forkIconStyle   = statIconStyle.Copy().Foreground(colorSubtle)
	issueIconStyle  = statIconStyle.Copy().Foreground(colorSuccess)
	prIconStyle     = statIconStyle.Copy().Foreground(colorDone)
	langIconStyle   = statIconStyle.Copy().Foreground(colorAccent)
	branchIconStyle = statIconStyle.Copy().Foreground(colorSubtle)
	diskIconStyle   = statIconStyle.Copy().Foreground(colorSubtle)

	labelStyle = lipgloss.NewStyle().
			Foreground(colorSubtle).
			Width(14)
	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#24292F", Dark: "#C9D1D9"})

	sectionTitleStyle = lipgloss.NewStyle().
				Foreground(colorAccent).
				Bold(true).
				MarginTop(1).
				MarginBottom(1)
)

type githubRepoInfo struct {
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
	repoViewCmd := exec.Command("gh", "repo", "view", repoArg, "--json", "name,description,licenseInfo,updatedAt,owner,stargazerCount,forkCount,isArchived,isPrivate,primaryLanguage,issues,pullRequests,defaultBranchRef,diskUsage,latestRelease")
	repoViewOutput, err := repoViewCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error getting repo info from gh: %s", string(repoViewOutput))
	}

	var repoInfo githubRepoInfo

	if err := json.Unmarshal(repoViewOutput, &repoInfo); err != nil {
		return "", fmt.Errorf("error parsing gh output: %w", err)
	}

	return r.renderUI(&repoInfo), nil
}

func (r *GithubPreviewProvider) renderUI(info *githubRepoInfo) string {
	var b strings.Builder

	// Header
	title := titleStyle.Render(fmt.Sprintf(" %s /  %s", info.Owner.Login, info.Name))

	var badges []string
	if info.IsPrivate {
		badges = append(badges, badgePrivateStyle.Render("󰌾 Private"))
	} else {
		badges = append(badges, badgePrivateStyle.Render("󰛓 Public"))
	}
	if info.IsArchived {
		badges = append(badges, badgeArchivedStyle.Render(" Archived"))
	}

	headerRow := lipgloss.JoinHorizontal(lipgloss.Top, title, strings.Join(badges, ""))
	b.WriteString(headerRow + "\n\n")

	// Description
	desc := info.Description
	if desc == "" {
		desc = "No description provided."
	}
	desc = lipgloss.NewStyle().Width(r.width - 2).Render(desc)
	b.WriteString(descStyle.Render(desc) + "\n\n")

	// Stats Grid
	starStat := statItemStyle.Render(starIconStyle.Render("") + formatCount(info.StargazerCount))
	forkStat := statItemStyle.Render(forkIconStyle.Render("") + formatCount(info.ForkCount))
	issueStat := statItemStyle.Render(issueIconStyle.Render("") + formatCount(info.Issues.TotalCount))
	prStat := statItemStyle.Render(prIconStyle.Render("") + formatCount(info.PullRequests.TotalCount))

	statsRow1 := lipgloss.JoinHorizontal(lipgloss.Top, starStat, forkStat, issueStat, prStat)
	b.WriteString(statsRow1 + "\n")

	langName := info.PrimaryLanguage.Name
	if langName == "" {
		langName = "Unknown"
	}
	langStat := statItemStyle.Render(langIconStyle.Render("") + langName)

	branchName := info.DefaultBranchRef.Name
	if branchName == "" {
		branchName = "unknown"
	}
	branchStat := statItemStyle.Render(branchIconStyle.Render("") + branchName)

	diskStat := statItemStyle.Render(diskIconStyle.Render("") + formatSize(info.DiskUsage))

	statsRow2 := lipgloss.JoinHorizontal(lipgloss.Top, langStat, branchStat, diskStat)
	b.WriteString(statsRow2 + "\n\n")

	// Info List
	licenseName := info.LicenseInfo.Name
	if licenseName == "" {
		licenseName = "No License"
	}

	licenseRow := lipgloss.JoinHorizontal(lipgloss.Left,
		labelStyle.Render(" License:"),
		valueStyle.Render(licenseName),
	)

	updatedRow := lipgloss.JoinHorizontal(lipgloss.Left,
		labelStyle.Render("󰥔 Updated:"),
		valueStyle.Render(relativeTime(info.UpdatedAt)),
	)

	b.WriteString(licenseRow + "\n")
	b.WriteString(updatedRow + "\n")

	if info.LatestRelease != nil && info.LatestRelease.TagName != "" {
		b.WriteString("\n" + sectionTitleStyle.Render("Latest Release") + "\n")

		name := info.LatestRelease.Name
		if name == "" {
			name = info.LatestRelease.TagName
		}

		tagRow := lipgloss.JoinHorizontal(lipgloss.Left,
			labelStyle.Render(" Tag:"),
			valueStyle.Render(info.LatestRelease.TagName),
		)
		nameRow := lipgloss.JoinHorizontal(lipgloss.Left,
			labelStyle.Render("󰅂 Name:"),
			valueStyle.Render(name),
		)
		publishedRow := lipgloss.JoinHorizontal(lipgloss.Left,
			labelStyle.Render("󰥔 Published:"),
			valueStyle.Render(relativeTime(info.LatestRelease.PublishedAt)),
		)
		urlRow := lipgloss.JoinHorizontal(lipgloss.Left,
			labelStyle.Render(" URL:"),
			lipgloss.NewStyle().Foreground(colorAccent).Underline(true).Render(info.LatestRelease.URL),
		)

		b.WriteString(tagRow + "\n")
		b.WriteString(nameRow + "\n")
		b.WriteString(publishedRow + "\n")
		b.WriteString(urlRow + "\n")
	}

	return b.String()
}

func formatCount(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	return fmt.Sprintf("%.1fk", float64(n)/1000.0)
}

func formatSize(kb int) string {
	if kb < 1024 {
		return fmt.Sprintf("%d KB", kb)
	}
	mb := float64(kb) / 1024.0
	if mb < 1024 {
		return fmt.Sprintf("%.1f MB", mb)
	}
	gb := mb / 1024.0
	return fmt.Sprintf("%.1f GB", gb)
}

func relativeTime(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	case d < 30*24*time.Hour:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	case d < 365*24*time.Hour:
		return fmt.Sprintf("%dmo ago", int(d.Hours()/(24*30)))
	default:
		return fmt.Sprintf("%dy ago", int(d.Hours()/(24*365)))
	}
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

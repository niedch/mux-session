package github

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	colorAccent  = lipgloss.AdaptiveColor{Light: "#044289", Dark: "#58A6FF"}
	colorSubtle  = lipgloss.AdaptiveColor{Light: "#57606A", Dark: "#8B949E"}
	colorWarning = lipgloss.AdaptiveColor{Light: "#9A6700", Dark: "#D29922"}
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

// RenderUI takes the fetched RepoInfo and generates a formatted UI string.
func RenderUI(info *RepoInfo, width int) string {
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
	desc = lipgloss.NewStyle().Width(width - 2).Render(desc)
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

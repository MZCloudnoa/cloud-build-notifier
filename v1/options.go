package notifier

import (
	"os"
	"strconv"
	"strings"
)

func (n *Notifier) getSubEnv(name string) string {
	val := strings.TrimSpace(getProp(n.build, "substitutions."+name))
	if val != "" {
		return val
	}

	return strings.TrimSpace(os.Getenv(name))
}

func (n *Notifier) getTitle() string {
	title := n.getSubEnv("_BUILD_TITLE")
	if title != "" {
		return title
	}

	if n.options != nil && n.options.Title != "" {
		return n.options.Title
	}

	return "<NO TITLE>"
}

func (n *Notifier) getRepoName() string {
	repoName := n.getSubEnv("_REPO_NAME")
	if repoName != "" {
		return repoName
	}

	if n.options != nil && n.options.RepoName != "" {
		return n.options.RepoName
	}

	return ""
}

func (n *Notifier) getNotifyURL() string {
	url := n.getSubEnv("_NOTIFY_URL")
	if url != "" {
		return url
	}

	if n.options != nil {
		return n.options.NotifyURL
	}

	return ""
}

func (n *Notifier) getNotifyExcludes() []BuildStatus {
	excludes := n.getSubEnv("_NOTIFY_EXCLUDES")
	if excludes != "" {
		ret := []BuildStatus{}
		for _, exclude := range strings.Split(excludes, ",") {
			exclude = strings.ToUpper(strings.TrimSpace(exclude))
			if !isValidBuildStatus(exclude) {
				panic("invalid status: " + exclude)
			}
			ret = append(ret, BuildStatus(exclude))
		}
		return ret
	}

	if n.options != nil {
		return n.options.Excludes
	}

	return []BuildStatus{}
}

func (n *Notifier) getPlatformType() PlatformType {
	platformType := n.getSubEnv("_NOTIFY_PLATFORM")
	platformType = strings.ToUpper(strings.TrimSpace(platformType))
	if platformType != "" {
		if isValidPlatformType(platformType) {
			return PlatformType(platformType)
		}
		return PlatformTypeSlack
	}

	if n.options != nil && n.options.Platform != "" {
		return n.options.Platform
	}

	return PlatformTypeSlack
}

func (n *Notifier) getDisabled() bool {
	disabled := n.getSubEnv("_NOTIFY_DISABLED")
	if disabled != "" {
		val, err := strconv.ParseBool(disabled)
		if err != nil {
			return false
		}
		return val
	}

	if n.options != nil && n.options.Disabled {
		return true
	}

	return false
}

func (n *Notifier) getDryRun() bool {
	dryRun := n.getSubEnv("_DRY_RUN")
	if dryRun != "" {
		val, err := strconv.ParseBool(dryRun)
		if err != nil {
			return false
		}
		return val
	}

	if n.options != nil && n.options.DryRun {
		return true
	}

	return false
}

func (n *Notifier) getQuietMode() bool {
	dryRun := n.getSubEnv("_QUIET_MODE")
	if dryRun != "" {
		val, err := strconv.ParseBool(dryRun)
		if err != nil {
			return false
		}
		return val
	}

	if n.options != nil && n.options.DryRun {
		return true
	}

	return false
}

func (n *Notifier) getTemplate(status BuildStatus) string {
	tmpl := n.getSubEnv("_TEMPLATE_" + string(status))
	if tmpl != "" {
		return tmpl
	}

	if n.options != nil {
		if tmpl, found := n.options.Templates[status]; found && tmpl != "" {
			return tmpl
		}
	}

	tmpl = n.getSubEnv("_DEFAULT_TEMPLATE")
	if tmpl != "" {
		return tmpl
	}

	if n.options != nil && n.options.DefaultTemplate != "" {
		return n.options.DefaultTemplate
	}

	platformType := n.getPlatformType()
	if platformType == PlatformTypeUnknown {
		return ""
	}

	return platforms[platformType].DefaultTemplate
}

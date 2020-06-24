package notifier

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/joho/godotenv"
)

// Options Options
type Options struct {
	Title           string                 // _BUILD_TITLE, default "<NO TITLE>"
	RepoName        string                 // _REPO_NAME
	NotifyURL       string                 // _NOTIFY_URL
	Excludes        []BuildStatus          // _NOTIFY_EXCLUDES
	Platform        PlatformType           // _NOTIFY_PLATFORM, default "SLACK"
	Disabled        bool                   // _NOTIFY_DISABLED, default "false"
	DryRun          bool                   // _DRY_RUN, default "false"
	QuietMode       bool                   // _QUIET_MODE, default "false"
	DefaultTemplate string                 // _DEFAULT_TEMPLATE
	Templates       map[BuildStatus]string // _TEMPLATE_{status}
}

// PubSubMessage PubSubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// Notifier Notifier
type Notifier struct {
	options   *Options
	build     map[string]interface{}
	quietMode bool
}

// NewNotifier NewNotifier
func NewNotifier(options *Options) *Notifier {
	return &Notifier{options: options}
}

func (n *Notifier) log(msg string) {
	if !n.quietMode {
		println(msg)
	}
}

const triggerURLFormat = "https://console.cloud.google.com/cloud-build/triggers/edit/%v?project=%v"
const projectURLFormat = "https://console.cloud.google.com/home/dashboard?project=%v"

// HandlePubSub HandlePubSub
func (n *Notifier) HandlePubSub(data []byte) error {
	godotenv.Load()

	err := json.Unmarshal(data, &n.build)
	if err != nil {
		return err
	}

	n.quietMode = n.getQuietMode()

	status, statusExists := n.build["status"]
	if !statusExists {
		return fmt.Errorf("status not exists")
	} else if !isValidBuildStatus(status.(string)) {
		return fmt.Errorf("invalid status: " + status.(string))
	}

	platformType := n.getPlatformType()
	if platformType == PlatformTypeUnknown {
		return fmt.Errorf("unknown platform")
	}

	buildStatus := BuildStatus(status.(string))
	excludes := n.getNotifyExcludes()
	for _, exclude := range excludes {
		if exclude == buildStatus {
			n.log("skip notify: " + string(buildStatus))
			return nil
		}
	}

	projectID := getProp(n.build, "projectId")

	branchName := getProp(n.build, "source.repoSource.branchName")
	if branchName == "" {
		branchName = n.getSubEnv("BRANCH_NAME")
	}

	tagName := getProp(n.build, "source.repoSource.tagName")
	if tagName == "" {
		tagName = n.getSubEnv("TAG_NAME")
	}

	commitSha := getProp(n.build, "sourceProvenance.resolvedRepoSource.commitSha")
	if commitSha == "" {
		commitSha = n.getSubEnv("COMMIT_SHA")
	}

	repoName := getProp(n.build, "source.repoSource.repoName")
	if repoName == "" {
		repoName = n.getRepoName()
	}

	gitInfo := parseGitRepo(repoName)

	params := map[string]interface{}{
		"Build":      n.build,
		"Title":      n.getTitle(),
		"Status":     status.(string),
		"ProjectID":  projectID,
		"ProjectURL": fmt.Sprintf(projectURLFormat, projectID),
		"TriggerID":  getProp(n.build, "buildTriggerId"),
		"TriggerURL": fmt.Sprintf(triggerURLFormat, getProp(n.build, "buildTriggerId"), getProp(n.build, "projectId")),
		"Git": map[string]interface{}{
			"Provider":        gitInfo.ProviderName(),
			"ProviderURL":     gitInfo.ProviderURL(),
			"Organization":    gitInfo.Organization(),
			"OrganizationURL": gitInfo.OrganizationURL(),
			"Repository":      gitInfo.Repository(),
			"RepositoryURL":   gitInfo.RepositoryURL(),
			"Branch":          branchName,
			"BranchURL":       gitInfo.BranchURL(branchName),
			"Tag":             tagName,
			"TagURL":          gitInfo.TagURL(tagName),
			"Commit":          commitSha,
			"CommitURL":       gitInfo.CommitURL(commitSha),
			"CommitInfoURL":   gitInfo.CommitInfoURL(commitSha),
		},
	}

	tmpl := n.getTemplate(buildStatus)
	if tmpl == "" {
		return fmt.Errorf("template missing: " + string(buildStatus))
	}

	msg, err := executeTemplate(tmpl, params)
	if err != nil {
		return fmt.Errorf("execute template failed: " + err.Error())
	}

	if n.getDryRun() {
		n.log(strings.Replace(msg, "\n", "", -1))
		return nil
	}

	url := n.getNotifyURL()
	if url == "" {
		return fmt.Errorf("notifyURL missing")
	}

	return sendNotification(url, msg, platformType)
}

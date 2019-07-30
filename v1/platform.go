package notifier

import "net/http"

// PlatformType PlatformType
type PlatformType string

const (
	// PlatformTypeUnknown PlatformTypeUnknown
	PlatformTypeUnknown PlatformType = "UNKNOWN"

	// PlatformTypeSlack PlatformTypeSlack
	PlatformTypeSlack PlatformType = "SLACK"

	// PlatformTypeHangoutsChat PlatformTypeHangoutsChat
	PlatformTypeHangoutsChat PlatformType = "HANGOUTS_CHAT"
)

var platformTypes = map[PlatformType]bool{
	PlatformTypeSlack:        true,
	PlatformTypeHangoutsChat: true,
}

func isValidPlatformType(platformType string) bool {
	if value, found := platformTypes[PlatformType(platformType)]; found {
		return value
	}

	return false
}

// Platform Platform
type Platform struct {
	DefaultTemplate string
	HTTPMethod      string
	HTTPHeaders     map[string]string
}

const templateSlack = `{
	"attachments": [{
		"fallback": "{{.Title}} : {{.Status}}",
		{{if or (eq .Status "QUEUED") (eq .Status "WORKING") (eq .Status "SUCCESS") -}}
		"color": "good",
		{{else if or (eq .Status "FAILURE") (eq .Status "TIMEOUT") (eq .Status "INTERNAL_ERROR") -}}
		"color": "#ff0000",
		{{else -}}
		"color": "#cccccc",
		{{end -}}
		"title": "{{.Title}} : {{.Status}}",
		"title_link": "{{.Build.logUrl}}",
		"text": "
			{{- if ne .ProjectID "" -}}\nProject: {{.ProjectID}}{{- end -}}
			{{- if ne .RepoName "" -}}\nRepository: {{.RepoName}}{{- end -}}
			{{- if ne .TagName "" -}}\nTag: {{.TagName}}{{- end -}}
			{{- if ne .BranchName "" -}}\nBranch: {{.BranchName}}{{- end -}}
			{{- if ne .TriggerID "" -}}\nTrigger: <{{.TriggerURL}}|{{.TriggerID}}>{{- end -}}
		",
	}]
}`

const templateHangoutsChat = `{ "text": "{{.Title}} : {{.Status}}" }`

var platforms = map[PlatformType]Platform{
	PlatformTypeSlack: Platform{
		DefaultTemplate: templateSlack,
		HTTPMethod:      http.MethodPost,
		HTTPHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	},

	PlatformTypeHangoutsChat: Platform{
		DefaultTemplate: templateHangoutsChat,
		HTTPMethod:      http.MethodPost,
		HTTPHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	},
}

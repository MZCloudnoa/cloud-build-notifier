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
		{{if (eq .Status "SUCCESS") -}}
		"color": "good",
		{{else if or (eq .Status "FAILURE") (eq .Status "TIMEOUT") (eq .Status "INTERNAL_ERROR") -}}
		"color": "#ff0000",
		{{else -}}
		"color": "#cccccc",
		{{end -}}
		"title": "{{.Title}} : {{.Status}}",
		"text": "
			{{- if ne .ProjectID "" -}}\nProject: <{{.ProjectURL}}|{{.ProjectID}}>{{- end -}}
			{{- if ne .Git.Repository "" -}}\nRepository: <{{.Git.ProviderURL}}|{{.Git.Provider}}>/<{{.Git.OrgnizationURL}}|{{.Git.Orgnization}}>/<{{.Git.RepositoryURL}}|{{.Git.Repository}}>{{- end -}}
			{{- if ne .Git.Branch "" -}}\nBranch: <{{.Git.BranchURL}}|{{.Git.Branch}}>{{- end -}}
			{{- if ne .Git.Tag "" -}}\nTag: <{{.Git.TagURL}}|{{.Git.Tag}}>{{- end -}}
			\n\n[ <{{.Build.logUrl}}|Log>
			{{- if ne .TriggerURL ""}} | <{{.TriggerURL}}|Trigger>{{- end -}}
			{{- if ne .Git.CommitURL ""}} | <{{.Git.CommitURL}}|Source>{{- end -}}
			{{- if ne .Git.CommitInfoURL ""}} | <{{.Git.CommitInfoURL}}|Commit>{{- end -}} ]
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

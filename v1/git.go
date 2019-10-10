package notifier

import (
	"fmt"
	"strings"
)

type gitProvider struct {
	name                 string
	baseURL              string
	orgnizationURLFormat string
	repositoryURLFormat  string
	branchURLFormat      string
	tagURLFormat         string
	commitURLFormat      string
	commitInfoURLFormat  string
}

func (g *gitProvider) Name() string {
	return g.name
}

func (g *gitProvider) URL() string {
	return g.baseURL
}

func (g *gitProvider) OrgnizationURL(orgnization string) string {
	return fmt.Sprintf(g.orgnizationURLFormat, g.baseURL, orgnization)
}

func (g *gitProvider) RepositoryURL(orgnization, repository string) string {
	return fmt.Sprintf(g.repositoryURLFormat, g.baseURL, orgnization, repository)
}

func (g *gitProvider) BranchURL(orgnization, repository, branch string) string {
	return fmt.Sprintf(g.branchURLFormat, g.baseURL, orgnization, repository, branch)
}

func (g *gitProvider) TagURL(orgnization, repository, tag string) string {
	return fmt.Sprintf(g.tagURLFormat, g.baseURL, orgnization, repository, tag)
}

func (g *gitProvider) CommitURL(orgnization, repository, commit string) string {
	return fmt.Sprintf(g.commitURLFormat, g.baseURL, orgnization, repository, commit)
}

func (g *gitProvider) CommitInfoURL(orgnization, repository, commit string) string {
	return fmt.Sprintf(g.commitInfoURLFormat, g.baseURL, orgnization, repository, commit)
}

var gitProviders = map[string]*gitProvider{
	"github": &gitProvider{
		name:                 "GitHub",
		baseURL:              "https://github.com",
		orgnizationURLFormat: "%[1]v/%[2]v",
		repositoryURLFormat:  "%[1]v/%[2]v/%[3]v",
		branchURLFormat:      "%[1]v/%[2]v/%[3]v/tree/%[4]v",
		tagURLFormat:         "%[1]v/%[2]v/%[3]v/tree/%[4]v",
		commitURLFormat:      "%[1]v/%[2]v/%[3]v/tree/%[4]v",
		commitInfoURLFormat:  "%[1]v/%[2]v/%[3]v/commit/%[4]v",
	},
	"bitbucket": &gitProvider{
		name:                 "Bitbucket",
		baseURL:              "https://bitbucket.org",
		orgnizationURLFormat: "%[1]v/%[2]v",
		repositoryURLFormat:  "%[1]v/%[2]v/%[3]v",
		branchURLFormat:      "%[1]v/%[2]v/%[3]v/src/%[4]v",
		tagURLFormat:         "%[1]v/%[2]v/%[3]v/src/%[4]v",
		commitURLFormat:      "%[1]v/%[2]v/%[3]v/src/%[4]v",
		commitInfoURLFormat:  "%[1]v/%[2]v/%[3]v/commits/%[4]v",
	},
}

var defaultGitProvider = &gitProvider{
	name:                 "Unknown",
	baseURL:              "",
	orgnizationURLFormat: "",
}

type gitInfo struct {
	provider     *gitProvider
	organization string
	repository   string
}

func (g *gitInfo) ProviderName() string {
	return g.provider.Name()
}

func (g *gitInfo) ProviderURL() string {
	return g.provider.URL()
}

func (g *gitInfo) Orgnization() string {
	return g.organization
}

func (g *gitInfo) OrgnizationURL() string {
	return g.provider.OrgnizationURL(g.organization)
}

func (g *gitInfo) Repository() string {
	return g.repository
}

func (g *gitInfo) RepositoryURL() string {
	return g.provider.RepositoryURL(g.organization, g.repository)
}

func (g *gitInfo) BranchURL(branch string) string {
	return g.provider.BranchURL(g.organization, g.repository, branch)
}

func (g *gitInfo) TagURL(tag string) string {
	return g.provider.TagURL(g.organization, g.repository, tag)
}

func (g *gitInfo) CommitURL(commit string) string {
	return g.provider.CommitURL(g.organization, g.repository, commit)
}

func (g *gitInfo) CommitInfoURL(commit string) string {
	return g.provider.CommitInfoURL(g.organization, g.repository, commit)
}

func parseGitRepo(repo string) *gitInfo {
	tokens := strings.SplitN(repo, "_", 3)
	if len(tokens) == 3 {
		if provider, found := gitProviders[tokens[0]]; found {
			return &gitInfo{
				provider:     provider,
				organization: tokens[1],
				repository:   tokens[2],
			}
		}
	}

	return &gitInfo{
		provider:   defaultGitProvider,
		repository: repo,
	}
}

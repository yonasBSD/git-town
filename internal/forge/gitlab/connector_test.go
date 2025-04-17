package gitlab_test

import (
	"testing"

	"github.com/git-town/git-town/v19/internal/cli/print"
	"github.com/git-town/git-town/v19/internal/config/configdomain"
	"github.com/git-town/git-town/v19/internal/forge/forgedomain"
	"github.com/git-town/git-town/v19/internal/forge/gitlab"
	"github.com/git-town/git-town/v19/internal/git/gitdomain"
	"github.com/git-town/git-town/v19/internal/git/giturl"
	. "github.com/git-town/git-town/v19/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestGitlabConnector(t *testing.T) {
	t.Parallel()

	t.Run("DefaultProposalMessage", func(t *testing.T) {
		t.Parallel()
		config := gitlab.Data{
			Data: forgedomain.Data{
				Hostname:     "",
				Organization: "",
				Repository:   "",
			},
			APIToken: None[configdomain.GitLabToken](),
		}
		give := forgedomain.Proposal{
			Number:       1,
			MergeWithAPI: true,
			Target:       "",
			Title:        "my title",
		}
		have := config.DefaultProposalMessage(give)
		want := "my title (!1)"
		must.EqOp(t, want, have)
	})

	t.Run("NewProposalURL", func(t *testing.T) {
		t.Parallel()
		tests := map[string]struct {
			branch gitdomain.LocalBranchName
			parent gitdomain.LocalBranchName
			want   string
		}{
			"top-level branch": {
				branch: "feature",
				parent: "main",
				want:   "https://gitlab.com/organization/repo/-/merge_requests/new?merge_request%5Bsource_branch%5D=feature&merge_request%5Btarget_branch%5D=main",
			},
			"stacked change": {
				branch: "feature-3",
				parent: "feature-2",
				want:   "https://gitlab.com/organization/repo/-/merge_requests/new?merge_request%5Bsource_branch%5D=feature-3&merge_request%5Btarget_branch%5D=feature-2",
			},
			"special characters in branch name": {
				branch: "feature-#",
				parent: "main",
				want:   "https://gitlab.com/organization/repo/-/merge_requests/new?merge_request%5Bsource_branch%5D=feature-%23&merge_request%5Btarget_branch%5D=main",
			},
		}
		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				connector := gitlab.Connector{
					Data: gitlab.Data{
						APIToken: None[configdomain.GitLabToken](),
						Data: forgedomain.Data{
							Hostname:     "gitlab.com",
							Organization: "organization",
							Repository:   "repo",
						},
					},
				}
				have, err := connector.NewProposalURL(tt.branch, tt.parent, "main", "", "")
				must.NoError(t, err)
				must.EqOp(t, tt.want, have)
			})
		}
	})
}

func TestNewGitlabConnector(t *testing.T) {
	t.Parallel()

	t.Run("GitLab SaaS", func(t *testing.T) {
		t.Parallel()
		remoteURL, has := giturl.Parse("git@gitlab.com:git-town/docs.git").Get()
		must.True(t, has)
		have, err := gitlab.NewConnector(gitlab.NewConnectorArgs{
			APIToken:  None[configdomain.GitLabToken](),
			Log:       print.Logger{},
			RemoteURL: remoteURL,
		})
		must.NoError(t, err)
		wantConfig := gitlab.Data{
			Data: forgedomain.Data{
				Hostname:     "gitlab.com",
				Organization: "git-town",
				Repository:   "docs",
			},
			APIToken: None[configdomain.GitLabToken](),
		}
		must.Eq(t, wantConfig, have.Data)
	})

	t.Run("custom URL", func(t *testing.T) {
		t.Parallel()
		remoteURL, has := giturl.Parse("git@custom-url.com:git-town/docs.git").Get()
		must.True(t, has)
		have, err := gitlab.NewConnector(gitlab.NewConnectorArgs{
			APIToken:  None[configdomain.GitLabToken](),
			Log:       print.Logger{},
			RemoteURL: remoteURL,
		})
		must.NoError(t, err)
		wantConfig := gitlab.Data{
			Data: forgedomain.Data{
				Hostname:     "custom-url.com",
				Organization: "git-town",
				Repository:   "docs",
			},
			APIToken: None[configdomain.GitLabToken](),
		}
		must.Eq(t, wantConfig, have.Data)
	})

	t.Run("hosted GitLab instance with custom SSH port", func(t *testing.T) {
		t.Parallel()
		remoteURL, has := giturl.Parse("git@gitlab.domain:1234/group/project").Get()
		must.True(t, has)
		have, err := gitlab.NewConnector(gitlab.NewConnectorArgs{
			APIToken:  None[configdomain.GitLabToken](),
			Log:       print.Logger{},
			RemoteURL: remoteURL,
		})
		must.NoError(t, err)
		wantConfig := gitlab.Data{
			Data: forgedomain.Data{
				Hostname:     "gitlab.domain",
				Organization: "group",
				Repository:   "project",
			},
			APIToken: None[configdomain.GitLabToken](),
		}
		must.Eq(t, wantConfig, have.Data)
	})
}

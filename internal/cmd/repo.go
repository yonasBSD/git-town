package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v18/internal/browser"
	"github.com/git-town/git-town/v18/internal/cli/flags"
	"github.com/git-town/git-town/v18/internal/cli/print"
	"github.com/git-town/git-town/v18/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v18/internal/config/configdomain"
	"github.com/git-town/git-town/v18/internal/execute"
	"github.com/git-town/git-town/v18/internal/forge"
	"github.com/git-town/git-town/v18/internal/forge/forgedomain"
	"github.com/git-town/git-town/v18/internal/git/gitdomain"
	. "github.com/git-town/git-town/v18/pkg/prelude"
	"github.com/spf13/cobra"
)

const repoDesc = "Open the repository homepage in the browser"

const repoHelp = `
Supported for repositories hosted on GitHub, GitLab, Gitea, and Bitbucket. Derives the Git provider from the "origin" remote. You can override this detection with "git config %s <DRIVER>" where DRIVER is "github", "gitlab", "gitea", or "bitbucket".

When using SSH identities, run "git config %s <HOSTNAME>" where HOSTNAME matches what is in your ssh config file.`

func repoCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:   "repo [remote]",
		Args:  cobra.MaximumNArgs(1),
		Short: repoDesc,
		Long:  cmdhelpers.Long(repoDesc, fmt.Sprintf(repoHelp, configdomain.KeyHostingPlatform, configdomain.KeyHostingOriginHostname)),
		RunE: func(cmd *cobra.Command, args []string) error {
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			return executeRepo(args, verbose)
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeRepo(args []string, verbose configdomain.Verbose) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: true,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	data, err := determineRepoData(args, repo)
	if err != nil {
		return err
	}
	browser.Open(data.connector.RepositoryURL(), repo.Frontend, repo.Backend)
	print.Footer(verbose, repo.CommandsCounter.Immutable(), repo.FinalMessages.Result())
	return nil
}

func determineRepoData(args []string, repo execute.OpenRepoResult) (data repoData, err error) {
	var remoteOpt Option[gitdomain.Remote]
	if len(args) > 0 {
		remoteOpt = gitdomain.NewRemote(args[0])
	} else {
		remoteOpt = Some(repo.UnvalidatedConfig.NormalConfig.DevRemote)
	}
	remote, hasRemote := remoteOpt.Get()
	if !hasRemote {
		return repoData{connector: nil}, nil
	}
	connectorOpt, err := forge.NewConnector(repo.UnvalidatedConfig, remote, print.Logger{})
	if err != nil {
		return data, err
	}
	connector, hasConnector := connectorOpt.Get()
	if !hasConnector {
		return data, forgedomain.UnsupportedServiceError()
	}
	return repoData{
		connector: connector,
	}, err
}

type repoData struct {
	connector forgedomain.Connector
}

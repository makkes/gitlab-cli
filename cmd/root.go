package cmd

import (
	"os"

	"github.com/makkes/gitlab-cli/v3/cmd/update"

	"github.com/makkes/gitlab-cli/v3/cmd/version"

	"github.com/makkes/gitlab-cli/v3/cmd/completion"

	"github.com/makkes/gitlab-cli/v3/api"
	"github.com/makkes/gitlab-cli/v3/cmd/issue"
	"github.com/makkes/gitlab-cli/v3/cmd/issues"
	"github.com/makkes/gitlab-cli/v3/cmd/login"
	"github.com/makkes/gitlab-cli/v3/cmd/pipeline"
	"github.com/makkes/gitlab-cli/v3/cmd/pipelines"
	"github.com/makkes/gitlab-cli/v3/cmd/project"
	"github.com/makkes/gitlab-cli/v3/cmd/projects"
	"github.com/makkes/gitlab-cli/v3/cmd/status"
	"github.com/makkes/gitlab-cli/v3/cmd/variable"
	"github.com/makkes/gitlab-cli/v3/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "gitlab",
	SilenceUsage: true,
}

func Execute(cfg config.Config) {
	var apiClient api.Client
	apiClient = api.NewAPIClient(cfg)
	rootCmd.AddCommand(projects.NewCommand(apiClient, cfg))
	rootCmd.AddCommand(project.NewCommand(apiClient))
	rootCmd.AddCommand(pipelines.NewCommand(apiClient))
	rootCmd.AddCommand(pipeline.NewCommand(apiClient))
	rootCmd.AddCommand(login.NewCommand(apiClient, cfg))
	rootCmd.AddCommand(issues.NewCommand(apiClient))
	rootCmd.AddCommand(issue.NewCommand(apiClient))
	rootCmd.AddCommand(variable.NewCommand(apiClient))
	rootCmd.AddCommand(status.NewCommand(apiClient, cfg))
	rootCmd.AddCommand(completion.NewCommand(rootCmd))
	rootCmd.AddCommand(version.NewCommand())
	rootCmd.AddCommand(update.NewCommand())
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

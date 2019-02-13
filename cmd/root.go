package cmd

import (
	"os"

	"github.com/makkes/gitlab-cli/cmd/completion"

	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/cmd/issue"
	"github.com/makkes/gitlab-cli/cmd/issues"
	"github.com/makkes/gitlab-cli/cmd/login"
	"github.com/makkes/gitlab-cli/cmd/pipeline"
	"github.com/makkes/gitlab-cli/cmd/pipelines"
	"github.com/makkes/gitlab-cli/cmd/project"
	"github.com/makkes/gitlab-cli/cmd/projects"
	"github.com/makkes/gitlab-cli/cmd/status"
	"github.com/makkes/gitlab-cli/cmd/variable"
	"github.com/makkes/gitlab-cli/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "gitlab-cli",
	SilenceUsage: true,
}

func Execute(cfg config.Config) {
	apiClient := api.NewAPIClient(cfg)
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
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

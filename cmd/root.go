package cmd

import (
	"os"

	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/cmd/completion"
	"github.com/makkes/gitlab-cli/cmd/get"
	"github.com/makkes/gitlab-cli/cmd/inspect"
	"github.com/makkes/gitlab-cli/cmd/login"
	"github.com/makkes/gitlab-cli/cmd/project"
	"github.com/makkes/gitlab-cli/cmd/status"
	"github.com/makkes/gitlab-cli/cmd/update"
	"github.com/makkes/gitlab-cli/cmd/variable"
	"github.com/makkes/gitlab-cli/cmd/version"
	"github.com/makkes/gitlab-cli/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "gitlab",
	SilenceUsage: true,
}

func Execute(cfg config.Config) {
	var apiClient api.Client
	apiClient = api.NewAPIClient(cfg)

	rootCmd.AddCommand(inspect.NewCommand(apiClient))
	rootCmd.AddCommand(get.NewCommand(apiClient, cfg))

	rootCmd.AddCommand(project.NewCommand(apiClient))
	rootCmd.AddCommand(login.NewCommand(apiClient, cfg))
	rootCmd.AddCommand(variable.NewCommand(apiClient))
	rootCmd.AddCommand(status.NewCommand(apiClient, cfg))
	rootCmd.AddCommand(completion.NewCommand(rootCmd))
	rootCmd.AddCommand(version.NewCommand())
	rootCmd.AddCommand(update.NewCommand())
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

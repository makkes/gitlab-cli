package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/makkes/gitlab-cli/api"
	apicmd "github.com/makkes/gitlab-cli/cmd/api"
	"github.com/makkes/gitlab-cli/cmd/completion"
	"github.com/makkes/gitlab-cli/cmd/create"
	"github.com/makkes/gitlab-cli/cmd/delete"
	"github.com/makkes/gitlab-cli/cmd/get"
	"github.com/makkes/gitlab-cli/cmd/inspect"
	"github.com/makkes/gitlab-cli/cmd/login"
	"github.com/makkes/gitlab-cli/cmd/status"
	"github.com/makkes/gitlab-cli/cmd/update"
	"github.com/makkes/gitlab-cli/cmd/version"
	"github.com/makkes/gitlab-cli/config"
)

var rootCmd = &cobra.Command{
	Use:          "gitlab",
	SilenceUsage: true,
}

func Execute(cfg config.Config) {
	apiClient := api.NewAPIClient(cfg)

	rootCmd.AddCommand(apicmd.NewCommand(apiClient, cfg))
	rootCmd.AddCommand(inspect.NewCommand(apiClient))
	rootCmd.AddCommand(get.NewCommand(apiClient, cfg))
	rootCmd.AddCommand(create.NewCommand(apiClient, cfg))
	rootCmd.AddCommand(delete.NewCommand(apiClient, cfg))

	rootCmd.AddCommand(login.NewCommand(apiClient, cfg))
	rootCmd.AddCommand(status.NewCommand(apiClient, cfg))
	rootCmd.AddCommand(completion.NewCommand(rootCmd))
	rootCmd.AddCommand(version.NewCommand())
	rootCmd.AddCommand(update.NewCommand())
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

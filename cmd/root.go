package cmd

import (
	"fmt"
	"os"

	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/cmd/login"
	"github.com/makkes/gitlab-cli/cmd/pipelines"
	"github.com/makkes/gitlab-cli/cmd/project"
	"github.com/makkes/gitlab-cli/cmd/projects"
	"github.com/makkes/gitlab-cli/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "gitlab-cli",
}

func Execute(cfg *config.Config) {
	apiClient := api.NewAPIClient(cfg)
	rootCmd.AddCommand(projects.NewCommand(apiClient))
	rootCmd.AddCommand(project.NewCommand(apiClient))
	rootCmd.AddCommand(pipelines.NewCommand(apiClient))
	rootCmd.AddCommand(login.NewCommand(apiClient, cfg))
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

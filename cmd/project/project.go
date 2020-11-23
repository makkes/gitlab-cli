package project

import (
	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/cmd/project/create"
	"github.com/spf13/cobra"
)

func NewCommand(client api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project COMMAND",
		Short: "Manage projects",
	}

	cmd.AddCommand(create.NewCommand(client))

	return cmd
}

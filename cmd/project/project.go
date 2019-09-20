package project

import (
	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/cmd/project/create"
	"github.com/makkes/gitlab-cli/cmd/project/inspect"
	"github.com/spf13/cobra"
)

func NewCommand(client api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project COMMAND",
		Short: "Manage projects",
	}

	cmd.AddCommand(create.NewCommand(client))
	cmd.AddCommand(inspect.NewCommand(client))

	return cmd
}

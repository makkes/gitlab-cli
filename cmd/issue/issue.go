package issue

import (
	"github.com/makkes/gitlab-cli/v3/api"
	"github.com/makkes/gitlab-cli/v3/cmd/issue/inspect"
	"github.com/spf13/cobra"
)

func NewCommand(client api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue COMMAND",
		Short: "Manage issues",
	}

	cmd.AddCommand(inspect.NewCommand(client))

	return cmd
}

package issue

import (
	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/cmd/issue/inspect"
	"github.com/spf13/cobra"
)

func NewCommand(client api.APIClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue COMMAND",
		Short: "Manage issues",
	}

	cmd.AddCommand(inspect.NewCommand(client))

	return cmd
}

package variable

import (
	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/cmd/variable/remove"
	"github.com/spf13/cobra"
)

func NewCommand(client api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "var",
		Short: "Manage project variables",
	}

	cmd.AddCommand(remove.NewCommand(client))

	return cmd
}

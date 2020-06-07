package variable

import (
	"github.com/makkes/gitlab-cli/v3/api"
	"github.com/makkes/gitlab-cli/v3/cmd/variable/create"
	"github.com/makkes/gitlab-cli/v3/cmd/variable/list"
	"github.com/makkes/gitlab-cli/v3/cmd/variable/remove"
	"github.com/spf13/cobra"
)

func NewCommand(client api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "var",
		Short: "Manage project variables",
	}

	cmd.AddCommand(create.NewCommand(client))
	cmd.AddCommand(list.NewCommand(client))
	cmd.AddCommand(remove.NewCommand(client))

	return cmd
}

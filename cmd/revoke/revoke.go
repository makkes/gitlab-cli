package revoke

import (
	"github.com/spf13/cobra"

	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/cmd/revoke/accesstoken"
)

func NewCommand(client api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke",
		Short: "Revoke an object, e.g. a project access token",
	}

	cmd.AddCommand(accesstoken.NewCommand(client))

	return cmd
}

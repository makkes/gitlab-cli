package revoke

import (
	"github.com/spf13/cobra"

	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/cmd/revoke/access_token"
)

func NewCommand(client api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke",
		Short: "Revoke an object, e.g. a project access token",
	}

	cmd.AddCommand(access_token.NewCommand(client))

	return cmd
}

package accesstoken

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/makkes/gitlab-cli/api"
)

func NewCommand(client api.Client) *cobra.Command {
	var projectFlag string

	cmd := &cobra.Command{
		Use:   "access-token ID|NAME",
		Args:  cobra.ExactArgs(1),
		Short: "Revoke a project access token",
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectFlag == "" {
				fmt.Printf("missing project name/ID\n\n")
				return cmd.Usage()
			}

			p, err := client.FindProject(projectFlag)
			if err != nil {
				return fmt.Errorf("failed to find project %q: %w", projectFlag, err)
			}

			atID := -1
			atl, err := client.GetAccessTokens(strconv.Itoa(p.ID))
			if err != nil {
				return fmt.Errorf("failed to list access tokens: %w", err)
			}

			reqAtID, _ := strconv.Atoi(args[0])
			for _, at := range atl {
				if at.Name == args[0] {
					atID = at.ID
					break
				}
				if at.ID == reqAtID {
					atID = at.ID
					break
				}
			}

			if atID == -1 {
				return fmt.Errorf("access token %q not found", args[0])
			}

			sc, err := client.Delete(fmt.Sprintf("/projects/%d/access_tokens/%d", p.ID, atID))
			if err != nil {
				return fmt.Errorf("failed revoking token: %w", err)
			}

			if sc != http.StatusNoContent {
				return fmt.Errorf("received unexpected status code from GitLab API: %d", sc)
			}

			fmt.Printf("revoked access token %d in %q\n", atID, p.Name)

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectFlag, "project", "p", "", "project to revoke the token from")

	return cmd
}

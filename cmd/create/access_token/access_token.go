package access_token

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/makkes/gitlab-cli/api"
)

var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyz")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func NewCommand(client api.Client) *cobra.Command {
	var projectFlag string
	var scopesFlag []string
	var nameFlag string

	cmd := &cobra.Command{
		Use:   "access-token",
		Short: "Create a project access token",
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectFlag == "" {
				fmt.Printf("missing project name/ID\n\n")
				return cmd.Usage()
			}

			if nameFlag == "" {
				nameFlag = randSeq(16)
			}

			p, err := client.FindProject(projectFlag)
			if err != nil {
				return fmt.Errorf("failed finding project: %w", err)
			}

			pat, err := client.CreateAccessToken(p.ID, nameFlag, time.Now().Add(24*time.Hour), scopesFlag)
			if err != nil {
				return fmt.Errorf("error creating access token: %w", err)
			}

			fmt.Fprintf(os.Stderr, "project access token %q created in %q\n", pat.Name, p.Name)
			fmt.Fprintf(os.Stdout, "%s\n", pat.Token)

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectFlag, "project", "p", "", "project for which to list the tokens")
	cmd.Flags().StringVarP(&nameFlag, "name", "n", "", "name of the new token")
	cmd.Flags().StringSliceVarP(&scopesFlag, "scopes", "s", []string{"api"}, "scopes to apply to the new access token")

	return cmd
}

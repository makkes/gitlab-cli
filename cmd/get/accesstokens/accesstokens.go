package accesstokens

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/cmd/get/output"
	"github.com/makkes/gitlab-cli/table"
)

func NewCommand(client api.Client, format *string) *cobra.Command {
	var projectFlag string

	cmd := &cobra.Command{
		Use:   "access-tokens",
		Short: "List access tokens in a project",
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectFlag == "" {
				fmt.Printf("missing project name/ID\n\n")
				return cmd.Usage()
			}

			p, err := client.FindProject(projectFlag)
			if err != nil {
				return fmt.Errorf("failed finding project: %w", err)
			}
			atl, err := client.GetAccessTokens(strconv.Itoa(p.ID))
			if err != nil {
				return fmt.Errorf("error retrieving access tokens: %w", err)
			}

			return output.NewPrinter(os.Stdout).Print(*format, func() error {
				table.PrintProjectAccessTokens(os.Stdout, atl)
				return nil
			}, func() error {
				for _, at := range atl {
					fmt.Fprintf(os.Stdout, "%s\n", at.Name)
				}
				return nil
			}, atl)
		},
	}

	cmd.Flags().StringVarP(&projectFlag, "project", "p", "", "project for which to list the tokens")

	return cmd
}

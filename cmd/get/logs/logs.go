package logs

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/makkes/gitlab-cli/api"
	"github.com/spf13/cobra"
)

func NewCommand(client api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Show logs of a job",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ids := strings.Split(args[0], ":")
			if len(ids) < 2 || ids[0] == "" || ids[1] == "" {
				return fmt.Errorf("ID must be of the form PROJECT_ID:JOB_ID")
			}
			logs, _, err := client.Get(fmt.Sprintf("/projects/%s/jobs/%s/trace", url.PathEscape(ids[0]), url.PathEscape(ids[1])))
			if err != nil {
				return fmt.Errorf("error retrieving logs: %w", err)
			}

			fmt.Printf("%s\n", logs)

			return nil
		},
	}

	return cmd
}

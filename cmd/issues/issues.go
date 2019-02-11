package issues

import (
	"encoding/json"
	"io"
	"os"
	"strconv"

	"github.com/makkes/gitlab-cli/table"

	"github.com/makkes/gitlab-cli/api"
	"github.com/spf13/cobra"
)

func issuesCommand(args []string, client api.Client, all bool, out io.Writer) error {
	project, err := client.FindProject(args[0])
	if err != nil {
		return err
	}
	path := "/projects/" + strconv.Itoa(project.ID) + "/issues"
	if !all {
		path += "?state=opened"
	}
	resp, _, err := client.Get(path)
	if err != nil {
		return err
	}

	issues := make([]api.Issue, 0)
	err = json.Unmarshal(resp, &issues)
	if err != nil {
		return err
	}

	table.PrintIssues(out, issues)
	return nil
}

func NewCommand(client api.APIClient) *cobra.Command {
	var all *bool
	cmd := &cobra.Command{
		Use:   "issues PROJECT",
		Short: "List issues in a project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return issuesCommand(args, client, *all, os.Stdout)
		},
	}

	all = cmd.Flags().BoolP("all", "a", false, "Show all issues (default shows just open)")
	return cmd
}

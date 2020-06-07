package issues

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/makkes/gitlab-cli/v3/table"

	"github.com/makkes/gitlab-cli/v3/api"
	"github.com/spf13/cobra"
)

func issuesCommand(args []string, client api.Client, all bool, page int, out io.Writer) error {
	project, err := client.FindProject(args[0])
	if err != nil {
		return err
	}
	path := "/projects/" + strconv.Itoa(project.ID) + fmt.Sprintf("/issues?page=%d", page)
	if !all {
		path += "&state=opened"
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

func NewCommand(client api.Client) *cobra.Command {
	var all *bool
	var page *int
	cmd := &cobra.Command{
		Use:   "issues PROJECT",
		Short: "List issues in a project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return issuesCommand(args, client, *all, *page, os.Stdout)
		},
	}

	all = cmd.Flags().BoolP("all", "a", false, "Show all issues (default shows just open)")
	page = cmd.Flags().Int("page", 1, "Page of results to display")
	return cmd
}

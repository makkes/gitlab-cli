package issues

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/makkes/gitlab-cli/cmd/get/output"

	"github.com/makkes/gitlab-cli/table"

	"github.com/makkes/gitlab-cli/api"
	"github.com/spf13/cobra"
)

func issuesCommand(scope string, format string, client api.Client, all bool, page int, out io.Writer) error {
	project, err := client.FindProject(scope)
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
	issuesIf := make([]interface{}, len(issues))
	for idx, i := range issues {
		issuesIf[idx] = i
	}

	return output.Print(resp, format, out, func() error {
		table.PrintIssues(out, issues)
		return nil
	}, func() error {
		for _, issue := range issues {
			fmt.Fprintf(out, "%s\n", issue.Title)
		}
		return nil
	}, issuesIf)
}

func NewCommand(client api.Client, project *string, format *string) *cobra.Command {
	var all *bool
	var page *int
	cmd := &cobra.Command{
		Use:   "issues",
		Short: "List issues in a project",
		RunE: func(cmd *cobra.Command, args []string) error {
			if project == nil || *project == "" {
				return fmt.Errorf("please provide a project scope")
			}
			return issuesCommand(*project, *format, client, *all, *page, os.Stdout)
		},
	}

	all = cmd.Flags().BoolP("all", "a", false, "Show all issues (default shows just open)")
	page = cmd.Flags().Int("page", 1, "Page of results to display")
	return cmd
}

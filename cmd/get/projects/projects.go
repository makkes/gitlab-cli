package projects

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/cmd/get/output"
	"github.com/makkes/gitlab-cli/config"
	"github.com/makkes/gitlab-cli/table"
	"github.com/spf13/cobra"
)

func projectsCommand(client api.Client, cfg config.Config, format string, page int, membership bool, out io.Writer) error {
	var path string
	if membership {
		path = fmt.Sprintf("/projects?membership=true&page=%d", page)
	} else {
		path = fmt.Sprintf("/users/${user}/projects?page=%d", page)
	}
	resp, status, err := client.Get(path)

	if err != nil {
		if status == 404 {
			return fmt.Errorf("cannot list projects: User %s not found. Please check your configuration", cfg.Get("user"))
		}
		return fmt.Errorf("Cannot list projects: %s", err)
	}
	projects := make([]api.Project, 0)
	err = json.Unmarshal(resp, &projects)
	if err != nil {
		return err
	}

	// put all project name => ID mappings into the cache
	for _, p := range projects {
		cfg.Cache().Put("projects", p.Name, strconv.Itoa(p.ID))
	}
	cfg.Write()

	return output.Print(resp, format, out, func() error {
		table.PrintProjects(out, projects)
		return nil
	}, func() error {
		for _, p := range projects {
			fmt.Fprintf(out, "%s\n", p.Name)
		}
		return nil
	}, projects)
}

func NewCommand(client api.Client, cfg config.Config, format *string) *cobra.Command {
	var page *int
	var membership *bool

	cmd := &cobra.Command{
		Use:   "projects",
		Short: "List all your projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			return projectsCommand(client, cfg, *format, *page, *membership, os.Stdout)
		},
	}

	page = cmd.Flags().Int("page", 1, "Page of results to display")
	membership = cmd.Flags().Bool("member", false, "Displays projects you are a member of")

	return cmd
}

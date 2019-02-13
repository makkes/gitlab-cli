package completion

import (
	"os"

	"github.com/spf13/cobra"
)

func NewCommand(rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion scripts",
		Long: `To load completions in the current shell run

. <(gitlab-cli completion)
		
To configure your bash shell to load completions for each session add the
following line to your ~/.bashrc or ~/.profile:

. <(gitlab-cli completion)
`,
		Run: func(cmd *cobra.Command, args []string) {
			rootCmd.GenBashCompletion(os.Stdout)
		},
	}
}

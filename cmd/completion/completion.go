package completion

import (
	"github.com/makkes/gitlab-cli/cmd/completion/bash"
	"github.com/makkes/gitlab-cli/cmd/completion/zsh"
	"github.com/spf13/cobra"
)

func NewCommand(rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts
		
To load completions in the current shell run

source <(gitlab completion SHELL)
		
To configure your bash shell to load completions for each session add the
following line to your ~/.bashrc or ~/.profile:

source <(gitlab completion bash)

If you use the zsh shell, run this command to permanently load completions:

gitlab completion zsh |sudo tee "${fpath[1]}/_gitlab"
`,
	}

	cmd.AddCommand(bash.NewCommand(rootCmd))
	cmd.AddCommand(zsh.NewCommand(rootCmd))

	return cmd
}

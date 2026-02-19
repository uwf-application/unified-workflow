package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [shell]",
		Short: "Generate shell completion",
		Long: `Generate shell completion scripts for uwf-cli.

To load completions:

Bash:
  $ source <(uwf-cli completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ uwf-cli completion bash > /etc/bash_completion.d/uwf-cli
  # macOS:
  $ uwf-cli completion bash > /usr/local/etc/bash_completion.d/uwf-cli

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it. You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ uwf-cli completion zsh > "${fpath[1]}/_uwf-cli"

  # You will need to start a new shell for this setup to take effect.

fish:
  $ uwf-cli completion fish | source

  # To load completions for each session, execute once:
  $ uwf-cli completion fish > ~/.config/fish/completions/uwf-cli.fish

PowerShell:
  PS> uwf-cli completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> uwf-cli completion powershell > uwf-cli.ps1
  # and source this file from your PowerShell profile.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		RunE:                  runCompletionCmd,
	}

	return cmd
}

func runCompletionCmd(cmd *cobra.Command, args []string) error {
	shell := args[0]

	switch shell {
	case "bash":
		return cmd.Root().GenBashCompletion(os.Stdout)
	case "zsh":
		return cmd.Root().GenZshCompletion(os.Stdout)
	case "fish":
		return cmd.Root().GenFishCompletion(os.Stdout, true)
	case "powershell":
		return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
	default:
		return fmt.Errorf("unsupported shell type %q", shell)
	}
}

package cmd

import (
	"fmt"
	"os"

	"dep/internal/command"

	"github.com/spf13/cobra"
)

var projectDir string

var rootCmd = &cobra.Command{
	Use:   "dep",
	Short: "Git multi-repository project manager",
	Long:  "DEP manages a collection of independent Git repositories as a single project.",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&projectDir, "project-dir", ".", "DEP project root directory")

	rootCmd.AddCommand(
		command.NewInitCmd(&projectDir),
		command.NewAddCmd(&projectDir),
		command.NewRemoveCmd(&projectDir),
		command.NewLockCmd(&projectDir),
		command.NewSyncCmd(&projectDir),
		command.NewStatusCmd(&projectDir),
		command.NewExecCmd(&projectDir),
	)
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
}

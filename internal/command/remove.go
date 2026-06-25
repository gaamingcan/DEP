package command

import (
	"fmt"
	"os"

	"dep/internal/config"
	"dep/internal/lock"
	"dep/internal/repo"

	"github.com/spf13/cobra"
)

// NewRemoveCmd creates the dep remove command.
func NewRemoveCmd(projectDir *string) *cobra.Command {
	var deleteLocal bool

	cmd := &cobra.Command{
		Use:   "remove <repo>",
		Short: "Remove a repository from the project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			identifier := args[0]

			cfg, lk, paths, err := loadProject(projectDir)
			if err != nil {
				return err
			}

			url, err := repo.ResolveURL(identifier, cfg.URLs())
			if err != nil {
				return err
			}

			if err := cfg.Remove(url); err != nil {
				return err
			}
			if err := config.Write(paths.ConfigPath, cfg); err != nil {
				return err
			}

			if err := lk.Remove(url); err != nil {
				return err
			}
			if err := lock.Write(paths.LockPath, lk); err != nil {
				return err
			}

			if deleteLocal {
				repoPath := repo.RepoPath(paths.Root, url)
				if err := os.RemoveAll(repoPath); err != nil {
					return fmt.Errorf("%s: delete local repository failed: %w", repo.DirName(url), err)
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&deleteLocal, "delete", false, "Delete the local repository directory")
	return cmd
}

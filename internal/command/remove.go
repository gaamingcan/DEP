package command

import (
	"fmt"
	"os"

	"dep/internal/repo"
	"dep/internal/store"

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

			proj, paths, err := loadProject(projectDir)
			if err != nil {
				return err
			}

			url, err := repo.ResolveURL(identifier, proj.URLs())
			if err != nil {
				return err
			}

			if err := proj.Remove(url); err != nil {
				return err
			}
			if err := store.Write(paths.LockPath, proj); err != nil {
				return err
			}

			if deleteLocal {
				repoPath := repo.RepoPath(paths.Root, url)
				if err := os.RemoveAll(repoPath); err != nil {
					return fmt.Errorf("%s: delete local repository failed: %w", url.DirName(), err)
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&deleteLocal, "delete", false, "Delete the local repository directory")
	return cmd
}

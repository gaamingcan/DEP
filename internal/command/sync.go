package command

import (
	"fmt"
	"os"

	"dep/internal/git"
	"dep/internal/repo"

	"github.com/spf13/cobra"
)

// NewSyncCmd creates the dep sync command.
func NewSyncCmd(projectDir *string) *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "Synchronize repositories to dep.lock commits",
		RunE: func(cmd *cobra.Command, args []string) error {
			proj, paths, err := loadProject(projectDir)
			if err != nil {
				return err
			}

			for _, r := range proj.Repositories {
				if r.Commit == nil {
					continue
				}
				url := r.URL
				lockedCommit := r.Commit
				dirName := url.DirName()
				repoPath := repo.RepoPath(paths.Root, url)

				if _, err := os.Stat(repoPath); os.IsNotExist(err) {
					fmt.Printf("%s: cloning...\n", dirName)
					if err := git.Clone(url.String(), repoPath); err != nil {
						return fmt.Errorf("%s: clone failed: %w", dirName, err)
					}
				} else if err != nil {
					return err
				}

				dirty, err := git.IsDirty(repoPath)
				if err != nil {
					return fmt.Errorf("%s: check status failed: %w", dirName, err)
				}
				if dirty {
					return fmt.Errorf("%s: working tree is dirty", dirName)
				}

				head, headErr := git.HeadCommit(repoPath)
				if headErr == nil && head == lockedCommit.String() {
					fmt.Printf("%s: already synced\n", dirName)
					continue
				}

				fmt.Printf("%s: fetching...\n", dirName)
				if err := git.Fetch(repoPath); err != nil {
					return fmt.Errorf("%s: fetch failed: %w", dirName, err)
				}

				fmt.Printf("%s: checking out %s...\n", dirName, lockedCommit.Short())
				if err := git.Checkout(repoPath, lockedCommit.String()); err != nil {
					return fmt.Errorf("%s: checkout failed: %w", dirName, err)
				}
			}

			return nil
		},
	}
}

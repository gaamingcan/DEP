package command

import (
	"fmt"

	"dep/internal/domain"
	"dep/internal/git"
	"dep/internal/repo"
	"dep/internal/store"

	"github.com/spf13/cobra"
)

// NewLockCmd creates the dep lock command.
func NewLockCmd(projectDir *string) *cobra.Command {
	return &cobra.Command{
		Use:   "lock",
		Short: "Rebuild dep.lock from repository HEAD commits",
		RunE: func(cmd *cobra.Command, args []string) error {
			proj, paths, err := loadProject(projectDir)
			if err != nil {
				return err
			}

			heads := make(map[domain.URL]domain.Commit, len(proj.Repositories))
			for _, r := range proj.Repositories {
				dirName := r.URL.DirName()
				repoPath := repo.RepoPath(paths.Root, r.URL)

				dirty, err := git.IsDirty(repoPath)
				if err != nil {
					return fmt.Errorf("%s: check status failed: %w", dirName, err)
				}
				if dirty {
					return fmt.Errorf("%s: working tree is dirty", dirName)
				}

				hash, err := git.HeadCommit(repoPath)
				if err != nil {
					return fmt.Errorf("%s: read HEAD failed: %w", dirName, err)
				}
				commit, err := domain.NewCommit(hash)
				if err != nil {
					return fmt.Errorf("%s: %w", dirName, err)
				}
				heads[r.URL] = commit
				fmt.Printf("[%s] commit: %s\n", dirName, commit.Short())
			}

			if err := proj.Lock(heads); err != nil {
				return err
			}
			return store.Write(paths.LockPath, proj)
		},
	}
}

package command

import (
	"fmt"

	"dep/internal/git"
	"dep/internal/lock"
	"dep/internal/repo"

	"github.com/spf13/cobra"
)

// NewLockCmd creates the dep lock command.
func NewLockCmd(projectDir *string) *cobra.Command {
	return &cobra.Command{
		Use:   "lock",
		Short: "Rebuild dep.lock from repository HEAD commits",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, _, paths, err := loadProject(projectDir)
			if err != nil {
				return err
			}

			headCommits := make(map[string]string, len(cfg.Repositories))
			for _, r := range cfg.Repositories {
				dirName := repo.DirName(r.URL)
				repoPath := repo.RepoPath(paths.Root, r.URL)

				dirty, err := git.IsDirty(repoPath)
				if err != nil {
					return fmt.Errorf("%s: check status failed: %w", dirName, err)
				}
				if dirty {
					return fmt.Errorf("%s: working tree is dirty", dirName)
				}

				commit, err := git.HeadCommit(repoPath)
				if err != nil {
					return fmt.Errorf("%s: read HEAD failed: %w", dirName, err)
				}
				headCommits[r.URL] = commit
			}

			lk, err := lock.Rebuild(cfg, headCommits)
			if err != nil {
				return err
			}
			return lock.Write(paths.LockPath, lk)
		},
	}
}

package command

import (
	"fmt"
	"os"

	"dep/internal/config"
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
			cfg, lk, paths, err := loadProject(projectDir)
			if err != nil {
				return err
			}

			for _, entry := range lk.Repositories {
				url := entry.URL
				lockedCommit := entry.Commit
				dirName := repo.DirName(url)
				repoPath := repo.RepoPath(paths.Root, url)

				if !cfg.ContainsURL(url) {
					if err := cfg.Add(url); err != nil {
						return fmt.Errorf("%s: add to config failed: %w", dirName, err)
					}
					if err := config.Write(paths.ConfigPath, cfg); err != nil {
						return fmt.Errorf("%s: write config failed: %w", dirName, err)
					}
					fmt.Printf("%s: added to dep.toml\n", dirName)
				}

				if _, err := os.Stat(repoPath); os.IsNotExist(err) {
					fmt.Printf("%s: cloning...\n", dirName)
					if err := git.Clone(url, repoPath); err != nil {
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
				if headErr == nil && head == lockedCommit {
					fmt.Printf("%s: already synced\n", dirName)
					continue
				}

				fmt.Printf("%s: fetching...\n", dirName)
				if err := git.Fetch(repoPath); err != nil {
					return fmt.Errorf("%s: fetch failed: %w", dirName, err)
				}

				fmt.Printf("%s: checking out %s...\n", dirName, lockedCommit[:7])
				if err := git.Checkout(repoPath, lockedCommit); err != nil {
					return fmt.Errorf("%s: checkout failed: %w", dirName, err)
				}
			}

			return nil
		},
	}
}

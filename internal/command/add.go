package command

import (
	"fmt"
	"os"

	"dep/internal/config"
	"dep/internal/git"
	"dep/internal/repo"

	"github.com/spf13/cobra"
)

// NewAddCmd creates the dep add command.
func NewAddCmd(projectDir *string) *cobra.Command {
	return &cobra.Command{
		Use:   "add <url>",
		Short: "Add a repository to the project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			url := args[0]
			if err := repo.ValidateURL(url); err != nil {
				return err
			}

			cfg, paths, err := loadConfig(projectDir)
			if err != nil {
				return err
			}

			dirName := repo.DirName(url)
			repoPath := repo.RepoPath(paths.Root, url)

			// URL already exists in config + directory exists = reject
			if cfg.ContainsURL(url) {
				if _, err := os.Stat(repoPath); err == nil {
					return fmt.Errorf("repository URL already exists: %s", url)
				}
				// URL exists but directory missing — repair case: clone and return
				if err := git.Clone(url, repoPath); err != nil {
					return fmt.Errorf("%s: clone failed: %w", dirName, err)
				}
				return nil
			}

			// New URL — validate no directory name conflict
			for _, r := range cfg.Repositories {
				if repo.DirName(r.URL) == dirName {
					return fmt.Errorf("derived directory %q already belongs to another repository", dirName)
				}
			}

			// Check existing path for various conflict conditions
			if info, err := os.Stat(repoPath); err == nil {
				if !info.IsDir() {
					return fmt.Errorf("derived path exists and is not a directory: %s", dirName)
				}
				if !git.IsRepository(repoPath) {
					return fmt.Errorf("derived directory exists and is not a Git repository: %s", dirName)
				}
				existingURL, remoteErr := git.RemoteOriginURL(repoPath)
				if remoteErr == nil && existingURL != url {
					return fmt.Errorf("derived directory %q exists and belongs to a different repository: %s", dirName, existingURL)
				}
				return fmt.Errorf("derived directory already exists: %s", dirName)
			} else if !os.IsNotExist(err) {
				return err
			}

			// Clone and add to config
			if err := git.Clone(url, repoPath); err != nil {
				return fmt.Errorf("%s: clone failed: %w", dirName, err)
			}
			if err := cfg.Add(url); err != nil {
				return err
			}
			return config.Write(paths.ConfigPath, cfg)
		},
	}
}

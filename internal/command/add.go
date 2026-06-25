package command

import (
	"fmt"
	"os"

	"dep/internal/domain"
	"dep/internal/git"
	"dep/internal/repo"
	"dep/internal/store"

	"github.com/spf13/cobra"
)

// NewAddCmd creates the dep add command.
func NewAddCmd(projectDir *string) *cobra.Command {
	return &cobra.Command{
		Use:   "add <url>",
		Short: "Add a repository to the project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			url, err := domain.NewURL(args[0])
			if err != nil {
				return err
			}

			proj, paths, err := loadProject(projectDir)
			if err != nil {
				return err
			}

			dirName := url.DirName()
			repoPath := repo.RepoPath(paths.Root, url)

			// URL already exists in project + directory exists = reject
			if proj.ContainsURL(url) {
				if _, err := os.Stat(repoPath); err == nil {
					return fmt.Errorf("repository URL already exists: %s", url)
				}
				// URL exists but directory missing — repair case: clone and return
				if err := git.Clone(url.String(), repoPath); err != nil {
					return fmt.Errorf("%s: clone failed: %w", dirName, err)
				}
				return nil
			}

			// New URL — validate no directory name conflict
			for _, r := range proj.Repositories {
				if r.URL.DirName() == dirName {
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
				if remoteErr == nil && existingURL != url.String() {
					return fmt.Errorf("derived directory %q exists and belongs to a different repository: %s", dirName, existingURL)
				}
				return fmt.Errorf("derived directory already exists: %s", dirName)
			} else if !os.IsNotExist(err) {
				return err
			}

			// Clone and add to project
			if err := git.Clone(url.String(), repoPath); err != nil {
				return fmt.Errorf("%s: clone failed: %w", dirName, err)
			}
			if err := proj.Add(url); err != nil {
				return err
			}
			return store.Write(paths.LockPath, proj)
		},
	}
}

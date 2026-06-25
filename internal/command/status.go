package command

import (
	"fmt"
	"os"

	"dep/internal/git"
	"dep/internal/repo"

	"github.com/spf13/cobra"
)

// NewStatusCmd creates the dep status command.
func NewStatusCmd(projectDir *string) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show abnormal repository states",
		RunE: func(cmd *cobra.Command, args []string) error {
			proj, paths, err := loadProject(projectDir)
			if err != nil {
				return err
			}

			abnormal := false
			for _, r := range proj.Repositories {
				dirName := r.URL.DirName()
				repoPath := repo.RepoPath(paths.Root, r.URL)

				if r.Commit == nil {
					fmt.Printf("%s unlocked\n", dirName)
					abnormal = true
					continue
				}
				lockedCommit := r.Commit

				info, err := os.Stat(repoPath)
				if os.IsNotExist(err) {
					fmt.Printf("%s missing\n", dirName)
					abnormal = true
					continue
				}
				if err != nil {
					return err
				}
				if !info.IsDir() {
					fmt.Printf("%s invalid\n", dirName)
					abnormal = true
					continue
				}
				if !git.IsRepository(repoPath) {
					fmt.Printf("%s invalid\n", dirName)
					abnormal = true
					continue
				}

				dirty, err := git.IsDirty(repoPath)
				if err != nil {
					return fmt.Errorf("%s: check status failed: %w", dirName, err)
				}
				if dirty {
					fmt.Printf("%s dirty\n", dirName)
					abnormal = true
					continue
				}

				head, err := git.HeadCommit(repoPath)
				if err != nil {
					return fmt.Errorf("%s: read HEAD failed: %w", dirName, err)
				}
				if head != lockedCommit.String() {
					fmt.Printf("%s drifted\n", dirName)
					abnormal = true
				}
			}

			if !abnormal {
				fmt.Println("all synced")
			}
			return nil
		},
	}
}

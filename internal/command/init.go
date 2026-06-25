package command

import (
	"fmt"
	"os"

	"dep/internal/domain"
	"dep/internal/store"

	"github.com/spf13/cobra"
)

// NewInitCmd creates the dep init command.
func NewInitCmd(projectDir *string) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a new DEP project",
		RunE: func(cmd *cobra.Command, args []string) error {
			paths, err := resolveProjectPaths(projectDir)
			if err != nil {
				return err
			}

			if _, err := os.Stat(paths.LockPath); err == nil {
				return fmt.Errorf("%s already exists", store.FileName)
			} else if !os.IsNotExist(err) {
				return err
			}

			proj := domain.NewProject()
			return store.Write(paths.LockPath, proj)
		},
	}
}

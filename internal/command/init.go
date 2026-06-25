package command

import (
	"fmt"
	"os"

	"dep/internal/config"
	"dep/internal/lock"

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

			if _, err := os.Stat(paths.ConfigPath); err == nil {
				return fmt.Errorf("%s already exists", config.FileName)
			} else if !os.IsNotExist(err) {
				return err
			}

			if _, err := os.Stat(paths.LockPath); err == nil {
				return fmt.Errorf("%s already exists", lock.FileName)
			} else if !os.IsNotExist(err) {
				return err
			}

			cfg := config.NewEmpty()
			if err := config.Write(paths.ConfigPath, cfg); err != nil {
				return err
			}

			lk := lock.NewEmpty()
			if err := lock.Write(paths.LockPath, lk); err != nil {
				return err
			}

			return nil
		},
	}
}

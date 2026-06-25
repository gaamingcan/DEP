package command

import (
	"fmt"
	"os"
	"path/filepath"

	"dep/internal/config"
	"dep/internal/lock"
)

// ProjectPaths holds resolved DEP project file paths.
type ProjectPaths struct {
	Root       string
	ConfigPath string
	LockPath   string
}

func resolveProjectPaths(projectDir *string) (ProjectPaths, error) {
	root, err := filepath.Abs(*projectDir)
	if err != nil {
		return ProjectPaths{}, fmt.Errorf("resolve project directory: %w", err)
	}
	return ProjectPaths{
		Root:       root,
		ConfigPath: filepath.Join(root, config.FileName),
		LockPath:   filepath.Join(root, lock.FileName),
	}, nil
}

func loadProject(projectDir *string) (*config.Config, *lock.Lock, ProjectPaths, error) {
	paths, err := resolveProjectPaths(projectDir)
	if err != nil {
		return nil, nil, ProjectPaths{}, err
	}

	cfg, err := config.Read(paths.ConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, paths, fmt.Errorf("%s not found in %s", config.FileName, paths.Root)
		}
		return nil, nil, paths, err
	}

	lk, err := lock.Read(paths.LockPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, paths, fmt.Errorf("%s not found in %s", lock.FileName, paths.Root)
		}
		return nil, nil, paths, err
	}

	return cfg, lk, paths, nil
}

// loadConfig reads only dep.toml, without requiring dep.lock.
func loadConfig(projectDir *string) (*config.Config, ProjectPaths, error) {
	paths, err := resolveProjectPaths(projectDir)
	if err != nil {
		return nil, ProjectPaths{}, err
	}

	cfg, err := config.Read(paths.ConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, paths, fmt.Errorf("%s not found in %s", config.FileName, paths.Root)
		}
		return nil, paths, err
	}

	return cfg, paths, nil
}

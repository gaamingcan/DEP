package command

import (
	"fmt"
	"os"
	"path/filepath"

	"dep/internal/domain"
	"dep/internal/store"
)

// ProjectPaths holds resolved DEP project file paths.
type ProjectPaths struct {
	Root     string
	LockPath string
}

func resolveProjectPaths(projectDir *string) (ProjectPaths, error) {
	root, err := filepath.Abs(*projectDir)
	if err != nil {
		return ProjectPaths{}, fmt.Errorf("resolve project directory: %w", err)
	}
	return ProjectPaths{
		Root:     root,
		LockPath: filepath.Join(root, store.FileName),
	}, nil
}

func loadProject(projectDir *string) (*domain.Project, ProjectPaths, error) {
	paths, err := resolveProjectPaths(projectDir)
	if err != nil {
		return nil, ProjectPaths{}, err
	}

	proj, err := store.Read(paths.LockPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, paths, fmt.Errorf("%s not found in %s", store.FileName, paths.Root)
		}
		return nil, paths, err
	}
	return proj, paths, nil
}

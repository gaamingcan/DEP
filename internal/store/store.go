package store

import (
	"fmt"
	"os"

	"dep/internal/domain"

	"github.com/pelletier/go-toml/v2"
)

// FileName is the single DEP project file.
const FileName = "dep.lock"

// fileEntry is the TOML serialisation form of a repository.
type fileEntry struct {
	URL    string  `toml:"url"`
	Commit *string `toml:"commit,omitempty"`
}

// lockFile is the TOML structure written to disk.
type lockFile struct {
	Version      int         `toml:"version"`
	Repositories []fileEntry `toml:"repo"`
}

// Read parses dep.lock from path into a domain Project.
func Read(path string) (*domain.Project, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var lf lockFile
	if err := toml.Unmarshal(data, &lf); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return fromFile(&lf)
}

// Write writes a domain Project to dep.lock at path.
// It copies entries internally and never mutates the input Project.
func Write(path string, proj *domain.Project) error {
	lf := toFile(proj)
	data, err := toml.Marshal(lf)
	if err != nil {
		return fmt.Errorf("marshal %s: %w", path, err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func fromFile(lf *lockFile) (*domain.Project, error) {
	repos := make([]domain.Repository, 0, len(lf.Repositories))
	for _, fe := range lf.Repositories {
		url, err := domain.NewURL(fe.URL)
		if err != nil {
			return nil, err
		}
		repo := domain.Repository{URL: url}
		if fe.Commit != nil {
			commit, err := domain.NewCommit(*fe.Commit)
			if err != nil {
				return nil, err
			}
			repo.Commit = &commit
		}
		repos = append(repos, repo)
	}
	return &domain.Project{
		Version:      lf.Version,
		Repositories: repos,
	}, nil
}

func toFile(proj *domain.Project) *lockFile {
	lf := &lockFile{
		Version:      proj.Version,
		Repositories: make([]fileEntry, len(proj.Repositories)),
	}
	for i, r := range proj.Repositories {
		fe := fileEntry{URL: r.URL.String()}
		if r.Commit != nil {
			s := r.Commit.String()
			fe.Commit = &s
		}
		lf.Repositories[i] = fe
	}
	return lf
}

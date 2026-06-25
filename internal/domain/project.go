package domain

import (
	"fmt"
	"sort"
)

const ProjectVersion = 1

// Repository represents a Git repository within the project.
type Repository struct {
	URL    URL
	Commit *Commit // nil means not locked
}

// Project is the single aggregate root for a DEP project.
type Project struct {
	Version      int
	Repositories []Repository
}

// NewProject returns an empty project.
func NewProject() *Project {
	return &Project{
		Version:      ProjectVersion,
		Repositories: []Repository{},
	}
}

// ContainsURL reports whether the URL is already registered.
func (p *Project) ContainsURL(url URL) bool {
	for _, r := range p.Repositories {
		if r.URL == url {
			return true
		}
	}
	return false
}

// Add registers a new repository.
func (p *Project) Add(url URL) error {
	if p.ContainsURL(url) {
		return fmt.Errorf("repository URL already exists: %s", url)
	}
	p.Repositories = append(p.Repositories, Repository{URL: url})
	p.sort()
	return nil
}

// Remove deletes a repository by URL.
func (p *Project) Remove(url URL) error {
	for i, r := range p.Repositories {
		if r.URL == url {
			p.Repositories = append(p.Repositories[:i], p.Repositories[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("repository not found: %s", url)
}

// Lock sets the locked commit for every repository.
func (p *Project) Lock(heads map[URL]Commit) error {
	for i, r := range p.Repositories {
		commit, ok := heads[r.URL]
		if !ok {
			return fmt.Errorf("missing HEAD commit for repository: %s", r.URL)
		}
		p.Repositories[i].Commit = &commit
	}
	return nil
}

// URLs returns all registered URLs.
func (p *Project) URLs() []URL {
	urls := make([]URL, len(p.Repositories))
	for i, r := range p.Repositories {
		urls[i] = r.URL
	}
	return urls
}

// CommitForURL returns the locked commit, if any.
func (p *Project) CommitForURL(url URL) (*Commit, bool) {
	for _, r := range p.Repositories {
		if r.URL == url {
			return r.Commit, r.Commit != nil
		}
	}
	return nil, false
}

func (p *Project) sort() {
	sort.Slice(p.Repositories, func(i, j int) bool {
		return p.Repositories[i].URL.Less(p.Repositories[j].URL)
	})
}

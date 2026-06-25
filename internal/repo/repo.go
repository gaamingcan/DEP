package repo

import (
	"fmt"
	"path/filepath"

	"dep/internal/domain"
)

// RepoPath joins project root with the derived directory name for a URL.
func RepoPath(projectRoot string, url domain.URL) string {
	return filepath.Join(projectRoot, url.DirName())
}

// ResolveURL finds a repository URL by local directory name or exact URL match.
func ResolveURL(identifier string, urls []domain.URL) (domain.URL, error) {
	for _, url := range urls {
		if url.DirName() == identifier || url.String() == identifier {
			return url, nil
		}
	}
	return domain.URL{}, fmt.Errorf("repository not found: %s", identifier)
}

// ResolveURLs resolves multiple identifiers to URLs.
func ResolveURLs(identifiers []string, allURLs []domain.URL) ([]domain.URL, error) {
	if len(identifiers) == 0 {
		out := make([]domain.URL, len(allURLs))
		copy(out, allURLs)
		return out, nil
	}

	seen := make(map[string]struct{}, len(identifiers))
	urls := make([]domain.URL, 0, len(identifiers))
	for _, id := range identifiers {
		url, err := ResolveURL(id, allURLs)
		if err != nil {
			return nil, err
		}
		if _, ok := seen[url.String()]; ok {
			continue
		}
		seen[url.String()] = struct{}{}
		urls = append(urls, url)
	}
	return urls, nil
}

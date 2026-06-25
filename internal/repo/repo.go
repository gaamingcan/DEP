package repo

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ValidateURL checks that a repository URL ends with .git.
func ValidateURL(url string) error {
	if url == "" {
		return fmt.Errorf("repository URL is empty")
	}
	if !strings.HasSuffix(url, ".git") {
		return fmt.Errorf("repository URL must end with .git: %s", url)
	}
	dirName := DirName(url)
	if dirName == "" {
		return fmt.Errorf("cannot derive directory name from URL: %s", url)
	}
	return nil
}

// DirName derives the local directory name from a repository URL.
func DirName(url string) string {
	last := url
	if idx := strings.LastIndexAny(url, "/:"); idx >= 0 && idx < len(url)-1 {
		last = url[idx+1:]
	}
	return strings.TrimSuffix(last, ".git")
}

// ResolveURL finds a repository URL by local directory name or exact URL match.
func ResolveURL(identifier string, urls []string) (string, error) {
	for _, url := range urls {
		if DirName(url) == identifier {
			return url, nil
		}
	}
	for _, url := range urls {
		if url == identifier {
			return url, nil
		}
	}
	return "", fmt.Errorf("repository not found: %s", identifier)
}

// ResolveURLs resolves multiple identifiers to URLs.
func ResolveURLs(identifiers, allURLs []string) ([]string, error) {
	if len(identifiers) == 0 {
		out := make([]string, len(allURLs))
		copy(out, allURLs)
		return out, nil
	}

	seen := make(map[string]struct{}, len(identifiers))
	urls := make([]string, 0, len(identifiers))
	for _, id := range identifiers {
		url, err := ResolveURL(id, allURLs)
		if err != nil {
			return nil, err
		}
		if _, ok := seen[url]; ok {
			continue
		}
		seen[url] = struct{}{}
		urls = append(urls, url)
	}
	return urls, nil
}

// RepoPath joins project root with the derived directory name for a URL.
func RepoPath(projectDir, url string) string {
	return filepath.Join(projectDir, DirName(url))
}

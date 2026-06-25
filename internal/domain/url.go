package domain

import (
	"fmt"
	"strings"
)

// URL is a validated Git repository URL.
// It guarantees .git suffix and non-empty directory derivation.
type URL struct {
	value string
}

// NewURL validates and creates a URL value object.
func NewURL(raw string) (URL, error) {
	if raw == "" {
		return URL{}, fmt.Errorf("repository URL is empty")
	}
	if !strings.HasSuffix(raw, ".git") {
		return URL{}, fmt.Errorf("repository URL must end with .git: %s", raw)
	}
	if dirName(raw) == "" {
		return URL{}, fmt.Errorf("cannot derive directory name from URL: %s", raw)
	}
	return URL{value: raw}, nil
}

// String returns the original URL string.
func (u URL) String() string {
	return u.value
}

// DirName returns the local directory name derived from the URL.
func (u URL) DirName() string {
	return dirName(u.value)
}

// Less compares URLs lexicographically for sorting.
func (u URL) Less(other URL) bool {
	return u.value < other.value
}

// dirName derives the local directory name from a raw URL string.
func dirName(raw string) string {
	last := raw
	if idx := strings.LastIndexAny(raw, "/:"); idx >= 0 && idx < len(raw)-1 {
		last = raw[idx+1:]
	}
	return strings.TrimSuffix(last, ".git")
}

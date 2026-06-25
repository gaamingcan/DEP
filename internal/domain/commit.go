package domain

import (
	"fmt"
	"regexp"
)

var commitHashPattern = regexp.MustCompile(`^[0-9a-f]{40}$`)

// Commit is a validated 40-character Git commit hash.
type Commit struct {
	value string
}

// NewCommit validates and creates a Commit value object.
func NewCommit(hash string) (Commit, error) {
	if !commitHashPattern.MatchString(hash) {
		return Commit{}, fmt.Errorf("commit must be a full 40-character hash: %s", hash)
	}
	return Commit{value: hash}, nil
}

// String returns the full 40-character hash.
func (c Commit) String() string {
	return c.value
}

// Short returns the first 7 characters for display.
func (c Commit) Short() string {
	return c.value[:7]
}

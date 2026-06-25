package git

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ErrNotInstalled indicates the git executable is not available.
var ErrNotInstalled = errors.New("git is not installed or not found in PATH")

func run(dir string, args ...string) (stdout, stderr string, err error) {
	cmd := exec.Command("git", args...)
	if dir != "" {
		cmd.Dir = dir
	}
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return "", "", ErrNotInstalled
		}
		if errBuf.Len() > 0 {
			return outBuf.String(), errBuf.String(), fmt.Errorf("%s", strings.TrimSpace(errBuf.String()))
		}
		return outBuf.String(), errBuf.String(), err
	}
	return outBuf.String(), errBuf.String(), nil
}

// runLive runs a git command with stderr piped to the terminal in real time,
// preserving progress output (e.g. clone/fetch object counting).
func runLive(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = io.Discard
	if err := cmd.Run(); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return ErrNotInstalled
		}
		return err
	}
	return nil
}

// Clone clones url into dir.
func Clone(url, dir string) error {
	return runLive("", "clone", url, dir)
}

// Fetch fetches objects for the repository in dir.
func Fetch(dir string) error {
	return runLive(dir, "fetch")
}

// Checkout checks out commit in dir quietly.
func Checkout(dir, commit string) error {
	return runLive(dir, "checkout", "--quiet", commit)
}

// HeadCommit returns the current HEAD commit hash in dir.
func HeadCommit(dir string) (string, error) {
	out, _, err := run(dir, "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// IsDirty reports whether the working tree has uncommitted changes.
func IsDirty(dir string) (bool, error) {
	out, _, err := run(dir, "status", "--porcelain")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(out) != "", nil
}

// Status returns the raw git status output.
func Status(dir string) (string, error) {
	out, stderr, err := run(dir, "status")
	if err != nil {
		if stderr != "" {
			return "", fmt.Errorf("%s", strings.TrimSpace(stderr))
		}
		return "", err
	}
	return strings.TrimRight(out, "\n"), nil
}

// Exec runs a command in dir and returns combined stdout/stderr output.
// The context can be used for cancellation and timeout.
func Exec(ctx context.Context, dir string, command string, args []string) (string, error) {
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Dir = dir
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return "", fmt.Errorf("%s: %w", command, ErrNotInstalled)
		}
		combined := strings.TrimSpace(outBuf.String() + errBuf.String())
		if combined != "" {
			return combined, fmt.Errorf("%s", combined)
		}
		return "", err
	}
	return strings.TrimRight(outBuf.String(), "\n"), nil
}

// IsRepository reports whether dir contains a valid Git repository.
func IsRepository(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, ".git"))
	return err == nil
}

// RemoteOriginURL returns the remote origin URL of the repository in dir.
func RemoteOriginURL(dir string) (string, error) {
	out, _, err := run(dir, "remote", "get-url", "origin")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

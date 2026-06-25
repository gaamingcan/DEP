package command

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"dep/internal/domain"
	"dep/internal/git"
	"dep/internal/repo"

	"github.com/spf13/cobra"
)

// NewExecCmd creates the dep exec command.
func NewExecCmd(projectDir *string) *cobra.Command {
	var parallel bool
	var timeout time.Duration

	cmd := &cobra.Command{
		Use:                "exec [repo] -- command [args...]",
		Short:              "Execute a command in project repositories",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			repoIdentifiers, commandName, commandArgs, err := parseExecArgs(args, &parallel, &timeout, projectDir)
			if err != nil {
				return err
			}
			if commandName == "" {
				return fmt.Errorf("missing command after --")
			}

			proj, paths, err := loadProject(projectDir)
			if err != nil {
				return err
			}

			urls, err := repo.ResolveURLs(repoIdentifiers, proj.URLs())
			if err != nil {
				return err
			}

			if parallel {
				return execParallel(urls, paths.Root, commandName, commandArgs, timeout)
			}
			return execSerial(urls, paths.Root, commandName, commandArgs)
		},
	}

	cmd.Flags().BoolVar(&parallel, "parallel", false, "Execute commands in parallel")
	return cmd
}

func parseExecArgs(args []string, parallel *bool, timeout *time.Duration, projectDir *string) (repoIdentifiers []string, commandName string, commandArgs []string, err error) {
	sep := -1
	for i, arg := range args {
		if arg == "--" {
			sep = i
			break
		}
	}
	if sep < 0 {
		return nil, "", nil, fmt.Errorf("usage: dep exec [repo] [-p] [--project-dir <dir>] [--timeout <duration>] -- command [args...]")
	}
	if sep == len(args)-1 {
		return nil, "", nil, fmt.Errorf("missing command after --")
	}

	before := args[:sep]
	after := args[sep+1:]
	commandName = after[0]
	commandArgs = after[1:]

	i := 0
	for i < len(before) {
		arg := before[i]
		switch arg {
		case "--parallel", "-p":
			*parallel = true
			i++
		case "--project-dir":
			if i+1 >= len(before) {
				return nil, "", nil, fmt.Errorf("--project-dir requires a value")
			}
			*projectDir = before[i+1]
			i += 2
		case "--timeout":
			if i+1 >= len(before) {
				return nil, "", nil, fmt.Errorf("--timeout requires a value (e.g. 30s, 5m)")
			}
			d, err := time.ParseDuration(before[i+1])
			if err != nil {
				return nil, "", nil, fmt.Errorf("invalid --timeout value: %w", err)
			}
			*timeout = d
			i += 2
		default:
			repoIdentifiers = append(repoIdentifiers, arg)
			i++
		}
	}

	return repoIdentifiers, commandName, commandArgs, nil
}

func execSerial(urls []domain.URL, projectRoot, commandName string, commandArgs []string) error {
	for _, url := range urls {
		dirName := url.DirName()
		repoPath := repo.RepoPath(projectRoot, url)
		output, err := git.Exec(context.Background(), repoPath, commandName, commandArgs)
		if output != "" {
			w := prefixLineWriter{
				prefix: fmt.Sprintf("[%s] ", dirName),
				out:    os.Stdout,
			}
			_, _ = w.Write([]byte(output))
		}
		if err != nil {
			return fmt.Errorf("%s: %w", dirName, err)
		}
	}
	return nil
}

func execParallel(urls []domain.URL, projectRoot, commandName string, commandArgs []string, timeout time.Duration) error {
	ctx := context.Background()
	var cancel context.CancelFunc

	if timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	var wg sync.WaitGroup
	errCh := make(chan error, len(urls))

	for _, url := range urls {
		url := url
		wg.Add(1)
		go func() {
			defer wg.Done()

			dirName := url.DirName()
			repoPath := repo.RepoPath(projectRoot, url)

			output, err := git.Exec(ctx, repoPath, commandName, commandArgs)
			if output != "" {
				w := prefixLineWriter{
					prefix: fmt.Sprintf("[%s] ", dirName),
					out:    os.Stdout,
				}
				_, _ = w.Write([]byte(output))
			}
			if err != nil {
				errCh <- fmt.Errorf("%s: %w", dirName, err)
				cancel()
			}
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return err
		}
	}
	return nil
}

type prefixLineWriter struct {
	prefix string
	out    *os.File
}

func (w *prefixLineWriter) Write(p []byte) (int, error) {
	scanner := bufio.NewScanner(strings.NewReader(string(p)))
	written := 0
	for scanner.Scan() {
		line := scanner.Text()
		if _, err := fmt.Fprintf(w.out, "%s%s\n", w.prefix, line); err != nil {
			return written, err
		}
		written += len(line) + 1
	}
	return len(p), scanner.Err()
}

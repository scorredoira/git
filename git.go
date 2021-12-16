// Package git allows to read git repositories
package git

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

func take(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

func (c *Commit) ShaShort() string {
	return take(c.Sha, 7)
}

// Summary returns the first line of the commit message
func (c *Commit) Summary() string {
	i := strings.Index(c.Message, "\n")
	if i == -1 {
		return c.Message
	}
	return c.Message[:i]
}

// ExtendedMessage returns the rest of the commit message beyond the summary.
func (c *Commit) ExtendedMessage() string {
	i := strings.Index(c.Message, "\n")
	if i != -1 {
		return c.Message[i+1:]
	}
	return ""
}

func (c *Commit) FomatDate() string {
	t, err := time.Parse("2006-01-02 15:04:05 -0700", c.Date)
	if err != nil {
		return err.Error()
	}

	return t.Format("02-01-06 15:04")
}

// ReadCommits reads the list of commits from a revision
func ReadCommits(path, revision string, skip, limit int) ([]*Commit, error) {
	cmd := "log --decorate --date=iso " + revision

	if skip > 0 {
		cmd += fmt.Sprintf(" --skip=%d", skip)
	}

	if limit > 0 {
		cmd += fmt.Sprintf(" --max-count=%d", limit)
	}

	r, err := execGit(path, cmd)
	if err != nil {
		return nil, err
	}

	return readCommits(r)
}

// ReadCommit reads a commit including diffs
func ReadCommit(path, sha string) (*Commit, error) {
	cmd := "log --decorate --max-count=1 --date=iso -p " + sha

	r, err := execGit(path, cmd)
	if err != nil {
		return nil, err
	}

	return readCommit(r)
}

func readCommits(r io.Reader) ([]*Commit, error) {
	p := newParser(r)
	p.run(parseCommits)
	return p.Commits, p.Error
}

func readCommit(r io.Reader) (*Commit, error) {
	p := newParser(r)
	p.run(parseCommit)
	return p.current(), p.Error
}

func execGit(path, command string) (io.Reader, error) {
	args := strings.Split(command, " ")
	cmd := exec.Command("git", args...)
	cmd.Dir = path
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return stdout, nil
}

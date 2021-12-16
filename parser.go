// Package git allows to read git repositories
package git

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

func parseCommits(p *parser) scanFn {
	l, ok := p.peek()
	if !ok || l == "" {
		return nil
	}

	if !strings.HasPrefix(l, "commit ") {
		p.Error = fmt.Errorf("Error parsing hash: %s", p.line)
		return nil
	}

	return parseCommit
}

func parseCommit(p *parser) scanFn {
	line, ok := p.next()
	if !ok {
		return nil
	}

	if !strings.HasPrefix(line, "commit ") {
		p.Error = fmt.Errorf("Error parsing hash: %s", p.line)
		return nil
	}

	c := &Commit{}
	p.Commits = append(p.Commits, c)

	line = line[len("commit "):]

	i := strings.Index(line, "(")
	if i > -1 {
		c.Sha = line[:i]
		// tags come after the sha, example: commit sha (HEAD, tag: 8.0, dev)
		tags := line[i+1 : len(line)-1]
		for _, t := range strings.Split(tags, ", ") {
			if strings.HasPrefix(t, "tag: ") {
				c.Tags = append(c.Tags, t[len("tag: "):])
			}
		}
	} else {
		c.Sha = line
	}

	return parseAuthor
}

func parseAuthor(p *parser) scanFn {
	line, ok := p.next()
	if !ok {
		return nil
	}

	if !strings.HasPrefix(line, "Author: ") {
		p.Error = fmt.Errorf("Error parsing Author: %s", p.line)
		return nil
	}

	c := p.current()

	line = line[len("Author: "):]

	i := strings.LastIndex(line, "<")
	if i > -1 {
		c.Author = line[:i]
		c.Email = line[i+1 : len(line)-1]
	} else {
		c.Author = line
	}

	return parseDate
}

func parseDate(p *parser) scanFn {
	line, ok := p.next()
	if !ok {
		return nil
	}

	if !strings.HasPrefix(line, "Date: ") {
		p.Error = fmt.Errorf("Error parsing Date: %s", p.line)
		return nil
	}

	c := p.current()
	c.Date = strings.TrimLeft(line[len("Date: "):], " ")

	return parseMessage
}

func parseMessage(p *parser) scanFn {
	line, ok := p.next()
	if !ok {
		return nil
	}

	// the start of the message is an empty line
	if line != "" {
		p.Error = fmt.Errorf("Error parsing message")
		return nil
	}

	var buff bytes.Buffer
	for {
		line, ok = p.next()
		if !ok || line == "" {
			break
		}
		// message lines are indented with 4 spaces.
		buff.WriteString(line[4:])
		buff.WriteRune('\n')
	}

	c := p.current()
	c.Message = buff.String()

	return parseDiff
}

func parseDiff(p *parser) scanFn {
	line, ok := p.peek()
	if !ok {
		return nil
	}

	if !strings.HasPrefix(line, "diff ") {
		return parseCommits
	}

	// now that we know its a diff, advance the reader
	p.next()

	return parseDiffFilePath
}

func parseDiffFilePath(p *parser) scanFn {
	var line string
	var ok bool
	for {
		line, ok = p.next()
		if !ok {
			return nil
		}
		if strings.HasPrefix(line, "--- ") && line != "--- /dev/null" {
			break
		}
		if strings.HasPrefix(line, "+++ ") && line != "+++ /dev/null" {
			break
		}
	}

	f := &DiffFile{Path: line[4:]}
	c := p.current()
	c.Files = append(c.Files, f)

	line, ok = p.peek()
	if !ok {
		return nil
	}
	if strings.HasPrefix(line, "+++ ") {
		p.next()
	}

	return parseDiffBlock
}

func parseDiffBlock(p *parser) scanFn {
	line, ok := p.peek()
	if !ok {
		return nil
	}

	if strings.HasPrefix(line, "diff ") {
		return parseDiff
	}

	if !strings.HasPrefix(line, "@@ ") {
		return parseCommits
	}

	// now that we know its a diff block, advance the reader
	p.next()

	d := &DiffBlock{Unified: line}
	c := p.current()
	f := c.Files[len(c.Files)-1]
	f.Diffs = append(f.Diffs, d)
	return parseDiffLine
}

func parseDiffLine(p *parser) scanFn {
	ln, ok := p.peek()
	if !ok {
		return nil
	}

	// If the file does not end with newline git emmits this error.
	// It is just a warning so we need to ignore it.
	if ln == `\ No newline at end of file` {
		p.next()
		return parseDiffLine
	}

	if ln == "" {
		return parseCommits
	}

	if strings.HasPrefix(ln, "@@ ") {
		return parseDiffBlock
	}

	if strings.HasPrefix(ln, "diff ") {
		return parseDiff
	}

	var lineType int

	switch ln[0] {
	case ' ':
		lineType = Unchanged
	case '+':
		lineType = Added
	case '-':
		lineType = Deleted
	default:
		// not a diff line
		return parseCommits
	}

	// now that we know its a diff line, advance the reader
	p.next()

	l := &DiffLine{Type: lineType, Text: ln}
	c := p.current()
	f := c.Files[len(c.Files)-1]
	d := f.Diffs[len(f.Diffs)-1]
	d.Lines = append(d.Lines, l)

	return parseDiffLine
}

func newParser(r io.Reader) *parser {
	return &parser{sc: bufio.NewScanner(r)}
}

type Commit struct {
	Sha     string
	Author  string
	Email   string
	Date    string
	Message string
	Tags    []string
	Files   []*DiffFile
}

type DiffFile struct {
	Path  string
	Diffs []*DiffBlock
}

type DiffBlock struct {
	Unified string
	Lines   []*DiffLine
}

type DiffLine struct {
	Text string
	Type int
}

const (
	Unchanged = iota
	Added
	Deleted
)

// Parser is a line reader that can buffer the last one
type parser struct {
	Commits  []*Commit
	Error    error
	sc       *bufio.Scanner
	line     string
	buffered bool
}

type scanFn func(*parser) scanFn

// Run starts the parsing
func (p *parser) run(fn scanFn) {
	for fn != nil {
		fn = fn(p)
	}
}

// current returns the last parsed commit.
func (p *parser) current() *Commit {
	if len(p.Commits) == 0 {
		return nil
	}
	return p.Commits[len(p.Commits)-1]
}

// Peek reads the next line without advancing the reader.
func (p *parser) peek() (string, bool) {
	if p.buffered {
		return p.line, true
	}

	ok := p.sc.Scan()
	if !ok {
		return "", false
	}

	p.line = p.sc.Text()
	p.buffered = true
	return p.line, true
}

// Next advances and reads the next line.
func (p *parser) next() (string, bool) {
	if p.buffered {
		p.buffered = false
		return p.line, true
	}

	ok := p.sc.Scan()
	if !ok {
		return "", false
	}

	p.line = p.sc.Text()
	return p.line, true
}

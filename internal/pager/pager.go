package pager

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Pager manages an external pager process (e.g. less, more).
// When no suitable pager is found it falls back to direct stdout writes.
type Pager struct {
	cmd    *exec.Cmd
	writer io.WriteCloser // pipe to the pager's stdin, or a nopCloser wrapping stdout
}

// New creates and starts a pager process.
//
// Resolution order:
//  1. MDVIEWER_PAGER environment variable
//  2. PAGER environment variable
//  3. "less -R" (ANSI-safe)
//  4. "more"
//  5. Direct stdout (no pager)
func New() *Pager {
	name, args := resolvePager()
	if name == "" {
		// No pager available – write directly to stdout.
		return &Pager{writer: nopCloser{os.Stdout}}
	}

	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	pipe, err := cmd.StdinPipe()
	if err != nil {
		// Cannot create pipe – fall back to stdout.
		return &Pager{writer: nopCloser{os.Stdout}}
	}

	if err := cmd.Start(); err != nil {
		// Pager binary failed to start – fall back to stdout.
		return &Pager{writer: nopCloser{os.Stdout}}
	}

	return &Pager{cmd: cmd, writer: pipe}
}

// Writer returns the io.Writer that callers should write rendered content to.
func (p *Pager) Writer() io.Writer {
	return p.writer
}

// Close signals the pager that input is complete and waits for it to exit.
// It is safe to call Close even when no external process is running.
func (p *Pager) Close() error {
	if err := p.writer.Close(); err != nil && p.cmd != nil {
		return fmt.Errorf("pager: close pipe: %w", err)
	}
	if p.cmd != nil {
		return p.cmd.Wait()
	}
	return nil
}

// resolvePager returns the pager binary name and arguments.
// An empty name means no usable pager was found.
func resolvePager() (string, []string) {
	candidates := []string{
		os.Getenv("MDVIEWER_PAGER"),
		os.Getenv("PAGER"),
		"less -R",
		"more",
	}

	for _, raw := range candidates {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}

		parts := strings.Fields(raw)
		name := parts[0]
		var args []string
		if len(parts) > 1 {
			args = parts[1:]
		}

		if _, err := exec.LookPath(name); err == nil {
			return name, args
		}
	}

	return "", nil
}

// nopCloser wraps an io.Writer so it satisfies io.WriteCloser.
// Close is a no-op (we never want to close real stdout).
type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }

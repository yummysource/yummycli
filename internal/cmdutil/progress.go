package cmdutil

import (
	"fmt"
	"io"
)

// Progress displays status messages during long-running operations.
// Implementations write to stderr so that stdout remains clean for JSON output.
// In non-terminal environments (pipes, CI), implementations should be no-ops
// so they do not pollute captured output.
type Progress interface {
	// Set replaces the current progress line with the given message.
	// Calling Set multiple times overwrites the previous message in place
	// when running in a terminal.
	Set(message string)

	// Clear erases the current progress line, leaving the cursor at the
	// start of a blank line. Call this before writing any final output.
	Clear()
}

// TerminalProgress writes progress messages to w using carriage-return overwriting.
// When isTerminal is false, all methods are no-ops so output is not polluted in
// piped or CI environments.
type TerminalProgress struct {
	out        io.Writer
	isTerminal bool
}

// NewProgress returns a TerminalProgress that writes to errOut.
// It is active only when isTerminal is true.
func NewProgress(errOut io.Writer, isTerminal bool) *TerminalProgress {
	return &TerminalProgress{out: errOut, isTerminal: isTerminal}
}

// Set writes message to stderr, overwriting the current line.
// No-op when not running in a terminal.
func (p *TerminalProgress) Set(message string) {
	if !p.isTerminal {
		return
	}
	// \r moves to the start of the line; the trailing spaces clear leftover chars.
	fmt.Fprintf(p.out, "\r%-80s", message)
}

// Clear erases the current progress line.
// No-op when not running in a terminal.
func (p *TerminalProgress) Clear() {
	if !p.isTerminal {
		return
	}
	fmt.Fprintf(p.out, "\r%-80s\r", "")
}

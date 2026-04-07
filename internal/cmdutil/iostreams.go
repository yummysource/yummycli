package cmdutil

import "io"

// IOStreams provides the input, output, and error streams for commands.
type IOStreams struct {
	In         io.Reader
	Out        io.Writer
	ErrOut     io.Writer
	IsTerminal bool
}

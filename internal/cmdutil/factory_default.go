package cmdutil

import (
	"os"

	"golang.org/x/term"
)

func NewDefault() *Factory {
	return &Factory{
		IOStreams: &IOStreams{
			In:         os.Stdin,
			Out:        os.Stdout,
			ErrOut:     os.Stderr,
			IsTerminal: term.IsTerminal(int(os.Stdin.Fd())),
		},
	}
}

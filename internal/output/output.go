// Package output provides formatted result writing for commands.
package output

import (
	"encoding/json"
	"io"
)

// Writer writes command results to an output stream.
type Writer struct {
	out io.Writer
}

// New creates a Writer that writes to w.
func New(w io.Writer) *Writer {
	return &Writer{out: w}
}

// JSON encodes v as JSON and writes it to the output stream.
func (w *Writer) JSON(v any) error {
	return json.NewEncoder(w.out).Encode(v)
}

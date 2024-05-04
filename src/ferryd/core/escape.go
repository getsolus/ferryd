package core

import (
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

// EscapedWriter is a writer that escapes non-ASCII sequences written to it using HTML encoding.
type EscapedWriter struct {
	w io.Writer
}

// NewEscapedWriter returns an initialised EscapedWriter.
func NewEscapedWriter(w io.Writer) *EscapedWriter {
	return &EscapedWriter{w: w}
}

// Write implements [io.Writer].
func (e *EscapedWriter) Write(b []byte) (n int, err error) {
	n = len(b)
	escaped := make([]byte, 0, len(b))

	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		escaped = e.appendRune(escaped, r)
		b = b[size:]
	}

	if _, err = e.w.Write(escaped); err != nil {
		return 0, fmt.Errorf("failed to write escaped bytes: %w", err)
	}

	return n, nil
}

func (e *EscapedWriter) appendRune(b []byte, r rune) []byte {
	switch {
	case r == utf8.RuneError:
		return utf8.AppendRune(b, unicode.ReplacementChar)
	case e.mustEscape(r):
		return append(b, []byte(fmt.Sprintf("&#x%X;", uint32(r)))...)
	default:
		return utf8.AppendRune(b, r)
	}
}

func (e *EscapedWriter) mustEscape(r rune) bool {
	return r > unicode.MaxASCII && unicode.In(r, unicode.Symbol, unicode.C, unicode.Space)
}

// Package textstat counts the lines, words and bytes in a stream of text.
package textstat

import (
	"bufio"
	"io"
	"unicode"
	"unicode/utf8"
)

// Stats is the result of counting one input.
type Stats struct {
	Lines int
	Words int
	Bytes int
}

// Count reads r to the end and returns its line, word and byte counts. A line
// is anything terminated by a newline; a word is a maximal run of
// non-space runes, as in wc(1).
func Count(r io.Reader) (Stats, error) {
	var s Stats
	br := bufio.NewReader(r)
	inWord := false
	for {
		ch, size, err := br.ReadRune()
		if err == io.EOF {
			return s, nil
		}
		if err != nil {
			return s, err
		}
		s.Bytes += size
		if ch == '\n' {
			s.Lines++
		}
		switch {
		case unicode.IsSpace(ch):
			inWord = false
		case ch == utf8.RuneError && size == 1:
			// Invalid UTF-8 still counts toward bytes but not words.
		case !inWord:
			inWord = true
			s.Words++
		}
	}
}

// Add returns the element-wise sum of s and o.
func (s Stats) Add(o Stats) Stats {
	return Stats{Lines: s.Lines + o.Lines, Words: s.Words + o.Words, Bytes: s.Bytes + o.Bytes}
}

package timewarp

import (
	"bufio"
	"bytes"
	"io"
)

// Scanner represents a lexical scanner for timerangeQL
type Scanner struct {
	r *reader
}

// NewScanner initializes a new lexical scanner
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: &reader{r: bufio.NewReader(r)}}
}

// Scan returns the next token and the literal value
func (s *Scanner) Scan() (tok Token, pos Pos, lit string) {
	// read the next rune
	ch, pos := s.r.read()

	// If we see a whitespace, consume all contiguous whitespace
	// If we see a letter or digit, consume as an ident or a reserved word
	if isWhitespace(ch) {
		s.r.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) || isDigit(ch) {
		s.r.unread()
		return s.scanIdent()
	}

	// Otherwise read the individual character
	switch ch {
	case eof:
		return EOF, pos, ""
	case '(':
		return LPAREN, pos, ""
	case ')':
		return RPAREN, pos, ""
	}

	return ILLEGAL, pos, string(ch)
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() (tok Token, pos Pos, lit string) {
	// create a buffer and read the current character onto it
	ch, pos := s.r.read()

	var buf bytes.Buffer
	_, _ = buf.WriteRune(ch)

	// read every subsequent whitespace character onto the buffer.
	// Non-whitespace characters and eof will cause the loop to exit.
	for {
		if ch, _ = s.r.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.r.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return WS, pos, buf.String()
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanIdent() (tok Token, pos Pos, lit string) {
	// create a buffer and read the current character into it
	ch, pos := s.r.read()

	var buf bytes.Buffer
	_, _ = buf.WriteRune(ch)

	// read every subsequent ident character onto the buffer.
	// Non-ident characters and eof will cause the loop to exit.
	for {
		if ch, _ = s.r.read(); ch == eof {
			break
		} else if !isIdentChar(ch) {
			s.r.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	lit = buf.String()

	// if the literal matches a keyword then return the keyword
	if tok = Lookup(lit); tok != IDENT {
		return tok, pos, ""
	}

	return
}

// reader represents a buffered rune reader used by the scanner.  It provides a
// fixed length circular buffer that can be unread.
type reader struct {
	r   io.RuneScanner
	i   int // buffer index
	n   int // buffer char count
	pos Pos // last read rune position
	buf [3]struct {
		ch  rune
		pos Pos
	}
}

// ReadRune reads the next rune from te reader
// This is a wrapper function to implement the io.RuneReader interface.
// Note that this function does not return size.
func (r *reader) ReadRune() (ch rune, size int, err error) {
	ch, _ = r.read()
	if ch == eof {
		err = io.EOF
	}
	return
}

// UnreadRune pushes the previously read rune back onto the buffer.  This is a
// wrapper function to implement the io.RuneScanner interface
func (r *reader) UnreadRune() error {
	r.unread()
	return nil
}

// read reads the next rune from the reader.
func (r *reader) read() (ch rune, pos Pos) {
	// If we have unread characters then read them off the buffer first.
	if r.n > 0 {
		r.n--
		return r.curr()
	}

	// Read the next rune from the underlying reader.
	// Any error (including io.EOF) should return EOF
	ch, _, err := r.r.ReadRune()
	if err != nil {
		ch = eof
	} else if ch == '\r' {
		if ch, _, err := r.r.ReadRune(); err != nil {
			// nop
		} else if ch != '\n' {
			_ = r.r.UnreadRune()
		}
		ch = '\n'
	}

	// Save the character and the position to the buffer.
	r.i = (r.i + 1) % len(r.buf)
	buf := &r.buf[r.i]
	buf.ch, buf.pos = ch, r.pos

	// Update position, only counting EOF once
	if ch == '\n' {
		r.pos.Line++
		r.pos.Char = 0
	} else if ch != eof {
		r.pos.Char++
	}
	return r.curr()
}

// unread pushes the previously read rune back onto the buffer.
func (r *reader) unread() {
	r.n++
}

// curr returns the last read character and position.
func (r *reader) curr() (ch rune, pos Pos) {
	i := (r.i - r.n + len(r.buf)) % len(r.buf)
	buf := &r.buf[i]
	return buf.ch, buf.pos
}

// eof represents an end of file
const eof = rune(0)

// isWhitespace returns true if the rune is a space, tab, or newline.
func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }

// isLetter returns true if the rune is a letter.
func isLetter(ch rune) bool { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') }

// isDigit returns true if the rune is a digit.
func isDigit(ch rune) bool { return (ch >= '0' && ch <= '9') }

// isIdentChar returns true if the rune is an ident character.
func isIdentChar(ch rune) bool { return isLetter(ch) || isDigit(ch) }

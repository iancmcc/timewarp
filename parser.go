package timewarp

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

const timefmt = "1504"

// Parser represents a wrapper for scanner to add a buffer.
// It provides a fixed-length circular buffer that can be unread.
type Parser struct {
	s   *Scanner
	i   int // buffer index
	n   int // buffer size
	buf [3]struct {
		tok Token
		pos Pos
		lit string
	}
}

// NewParser instantiates a parser
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// Parse returns a filter for the provided statement
func (p *Parser) Parse() (f Filter, err error) {
	f, err = p.parseExpr()
	if err != nil {
		return nil, err
	}

	tok, pos, lit := p.scanIgnoreWhitespace()
	if tok != EOF {
		return nil, newParseError(tokstr(tok, lit), []string{"EOF"}, pos)
	}
	return
}

// ParseFilter returns a filter for each individual statement.
func (p *Parser) ParseFilter() (f Filter, err error) {
	// inspect the first token
	tok, pos, lit := p.scanIgnoreWhitespace()
	switch tok {
	case LPAREN:
		f, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
		tok, pos, lit := p.scanIgnoreWhitespace()
		if tok != RPAREN {
			return nil, newParseError(tokstr(tok, lit), []string{")"}, pos)
		}
		return
	case NOT:
		f, err = p.ParseFilter()
		if err != nil {
			return nil, err
		}
		f = f.Not()
		return
	case YEAR:
		return p.parseYearFilter()
	case MONTH:
		return p.parseMonthFilter(0)
	case WEEK:
		return p.parseWeekFilter(0)
	case DAY:
		return p.parseDayFilter(0)
	case TIME:
		return p.parseTimeFilter()
	default:
		return nil, newParseError(tokstr(tok, lit), []string{"(", "NOT", "YEAR", "MONTH", "WEEK", "DAY", "TIME"}, pos)
	}
}

// parseExpr returns the resulting filter from joining multiple filters.
func (p *Parser) parseExpr() (f Filter, err error) {
	// read the first statement
	f, err = p.ParseFilter()
	if err != nil {
		return nil, err
	}

	for {
		tok, pos, lit := p.scanIgnoreWhitespace()
		switch tok {
		case EOF:
			return f, nil
		case RPAREN:
			p.unscan()
			return f, nil
		case AND:
			filter, err := p.ParseFilter()
			if err != nil {
				return nil, err
			}
			f = f.Union(filter)
		case IN:
			filter, err := p.ParseFilter()
			if err != nil {
				return nil, err
			}
			f = f.Intersect(filter)
		case OF:
			tok, pos, lit = p.scanIgnoreWhitespace()

			var v = 1
			if tok == IDENT {
				var err error
				v, err = strconv.Atoi(lit)
				if err != nil {
					return nil, &ParseError{
						Message: "unable to parse number",
						Pos:     pos,
					}
				} else if v == 0 {
					return nil, &ParseError{
						Message: "ordinal cannot be zero",
						Pos:     pos,
					}
				}
			} else {
				p.unscan()
			}
			filter, err := p.parseOrdinal(v)
			if err != nil {
				return nil, err
			}
			f = f.Ordinal(v, filter)
		default:
			return nil, newParseError(tokstr(tok, lit), []string{"AND", "IN", "OF"}, pos)
		}
	}
}

// parseOrdinal handles sub-filter values under token "OF"
func (p *Parser) parseOrdinal(v int) (f Filter, err error) {
	// inspect the first token
	tok, pos, lit := p.scanIgnoreWhitespace()
	switch tok {
	case MONTH:
		return p.parseMonthFilter(v)
	case WEEK:
		return p.parseWeekFilter(v)
	case DAY:
		return p.parseDayFilter(v)
	default:
		return nil, newParseError(tokstr(tok, lit), []string{"MONTH", "WEEK", "DAY"}, pos)
	}
}

// parseYearFilter returns a filter for a given year
func (p *Parser) parseYearFilter() (f Filter, err error) {
	tok, pos, lit := p.scanIgnoreWhitespace()
	switch tok {
	case IDENT:
		v, err := strconv.Atoi(lit)
		if err != nil {
			return nil, &ParseError{
				Message: "unable to parse year",
				Pos:     pos,
			}
		} else if v <= 0 {
			return nil, &ParseError{
				Message: "year must be greater than 0",
				Pos:     pos,
			}
		}

		return Year(v).Filter(), nil
	default:
		return nil, &ParseError{
			Message: "missing year",
			Pos:     pos,
		}
	}
}

// parseMonthFilter returns a filter for a given month
func (p *Parser) parseMonthFilter(v int) (f Filter, err error) {
	tok, _, _ := p.scanIgnoreWhitespace()

	var m time.Month
	if tok.isMonthOfYear() {
		m = getMonthOfYear(tok)
	} else {
		p.unscan()
	}

	if v != 0 {
		return TheMonth(m).Filter(), nil
	}
	return Month(m).Filter(), nil
}

// parseWeekFilter returns a filter for the given week
func (p *Parser) parseWeekFilter(v int) (f Filter, err error) {
	tok, _, _ := p.scanIgnoreWhitespace()

	var w time.Weekday
	if tok.isDayOfWeek() {
		w = getDayOfWeek(tok)
	} else {
		p.unscan()
		w = -1
	}

	if v > 0 {
		return TheWeek(w, 7, -v+1, 2*v-1).Filter(), nil
	}
	return Week(w, 7).Filter(), nil
}

// parseDayFilter returns a filter for the given day
func (p *Parser) parseDayFilter(v int) (f Filter, err error) {
	tok, pos, lit := p.scanIgnoreWhitespace()
	if tok == IDENT {
		d, err := strconv.Atoi(lit)
		if err != nil {
			return nil, &ParseError{
				Message: "could not parse days",
				Pos:     pos,
			}
		}

		tok, pos, lit = p.scanIgnoreWhitespace()
		if tok == IDENT {
			n, err := strconv.Atoi(lit)
			if err != nil {
				return nil, &ParseError{
					Message: "could not parse number of consecutive days",
					Pos:     pos,
				}
			}

			if v != 0 {
				return TheDays(d, n).Filter(), nil
			}

			return Days(d, n).Filter(), nil
		}

		p.unscan()
		if v != 0 {
			return TheDays(0, d).Filter(), nil
		}

		return Days(0, d).Filter(), nil

	} else if tok.isDayOfWeek() {
		if v != 0 {
			return nil, &ParseError{
				Message: "can not parse weekdays with ordinal, use WEEK instead",
				Pos:     pos,
			}
		}

		var (
			day   = getDayOfWeek(tok)
			delta = 1
		)

		tok, pos, lit = p.scanIgnoreWhitespace()
		if tok.isDayOfWeek() {
			delta = getWeekdayDelta(getDayOfWeek(tok), day)
		} else {
			p.unscan()
		}
		return Week(day, delta).Filter(), nil
	} else {
		p.unscan()
		if v != 0 {
			return TheDays(-v+1, 2*v-1).Filter(), nil
		}
		return Days(0, 1).Filter(), nil
	}
}

// parseTimeFilter returns a filter for the given time
func (p *Parser) parseTimeFilter() (f Filter, err error) {
	tok, pos, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, newParseError(tokstr(tok, lit), []string{"IDENT"}, pos)
	}

	if _, err := time.Parse(timefmt, lit); err != nil {
		return nil, &ParseError{
			Message: "invalid time format",
			Pos:     pos,
		}
	}
	t1 := lit

	tok, pos, lit = p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, newParseError(tokstr(tok, lit), []string{"IDENT"}, pos)
	}

	if _, err := time.Parse(timefmt, lit); err != nil {
		return nil, &ParseError{
			Message: "invalid time format",
			Pos:     pos,
		}
	}
	t2 := lit

	return Times(timefmt, t1, t2).Filter(), nil
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, pos Pos, lit string) {
	tok, pos, lit = p.scan()
	if tok == WS {
		tok, pos, lit = p.scan()
	}
	return
}

// scan returns the next token from the underlying scanner.
// If the token has been unscanned, then read that instead.
func (p *Parser) scan() (tok Token, pos Pos, lit string) {
	// If we have unread tokens on the buffer then read them off the buffer
	// first.
	if p.n > 0 {
		p.n--
		return p.curr()
	}

	// Read the next token from the scanner and write to the buffer.
	tok, pos, lit = p.s.Scan()
	p.i = (p.i + 1) % len(p.buf)
	buf := &p.buf[p.i]
	buf.tok, buf.pos, buf.lit = tok, pos, lit
	return
}

// unscan pushes the previously pushed token back onto the buffer
func (p *Parser) unscan() {
	p.n++
}

// curr returns the last read token, position, and literal
func (p *Parser) curr() (tok Token, pos Pos, lit string) {
	i := (p.i - p.n + len(p.buf)) % len(p.buf)
	buf := &p.buf[i]
	return buf.tok, buf.pos, buf.lit
}

// ParseError represents an error that occurred during parsing.
type ParseError struct {
	Message  string
	Found    string
	Expected []string
	Pos      Pos
}

// newParseError returns a new instance of ParseError.
func newParseError(found string, expected []string, pos Pos) *ParseError {
	return &ParseError{Found: found, Expected: expected, Pos: pos}
}

// Error returns the string representation of the error.
func (e *ParseError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s at %s", e.Message, e.Pos)
	}
	return fmt.Sprintf("found %s, expected %s at %s", e.Found, strings.Join(e.Expected, ", "), e.Pos)
}

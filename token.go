package timewarp

import (
	"fmt"
	"strings"
	"time"
)

// Token describes lexer terms.
type Token int

const (
	// ILLEGAL Token, EOF, and WS are special timerangeQL tokens
	ILLEGAL Token = iota
	EOF
	WS

	literalBeg
	// IDENT is a timerangeQL literal token
	IDENT // main
	literalEnd

	operatorBeg
	// AND and the following are timerangeQL operators
	AND // AND
	IN  // IN
	OF  // OF
	NOT // NOT
	operatorEnd

	keywordBeg
	// YEAR and the following are timerangeQL keywords
	YEAR  // YEAR
	MONTH // MONTH
	WEEK  // WEEK
	DAY   // DAY
	TIME  // TIME
	keywordEnd

	moyBeg
	// JANUARY and the following are timerangeQL keywords representing months
	// of the year.
	JANUARY   // JANUARY
	FEBRUARY  // FEBRUARY
	MARCH     // MARCH
	APRIL     // APRIL
	MAY       // MAY
	JUNE      // JUNE
	JULY      // JULY
	AUGUST    // AUGUST
	SEPTEMBER // SEPTEMBER
	OCTOBER   // OCTOBER
	NOVEMBER  // NOVEMBER
	DECEMBER  // DECEMBER
	moyEnd

	dowBeg
	// MONDAY and the following are timerangeQL keywords representing days of
	// the week.
	MONDAY    // MONDAY
	TUESDAY   // TUESDAY
	WEDNESDAY // WEDNESDAY
	THURSDAY  // THURSDAY
	FRIDAY    // FRIDAY
	SATURDAY  // SATURDAY
	SUNDAY    // SUNDAY
	dowEnd

	LPAREN // (
	RPAREN // )
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	WS:      "WS",

	IDENT: "IDENT",

	AND: "AND",
	IN:  "IN",
	OF:  "OF",
	NOT: "NOT",

	YEAR:  "YEAR",
	MONTH: "MONTH",
	WEEK:  "WEEK",
	DAY:   "DAY",
	TIME:  "TIME",

	JANUARY:   "JANUARY",
	FEBRUARY:  "FEBRUARY",
	MARCH:     "MARCH",
	APRIL:     "APRIL",
	MAY:       "MAY",
	JUNE:      "JUNE",
	JULY:      "JULY",
	AUGUST:    "AUGUST",
	SEPTEMBER: "SEPTEMBER",
	OCTOBER:   "OCTOBER",
	NOVEMBER:  "NOVEMBER",
	DECEMBER:  "DECEMBER",

	MONDAY:    "MONDAY",
	TUESDAY:   "TUESDAY",
	WEDNESDAY: "WEDNESDAY",
	THURSDAY:  "THURSDAY",
	FRIDAY:    "FRIDAY",
	SATURDAY:  "SATURDAY",
	SUNDAY:    "SUNDAY",

	LPAREN: "(",
	RPAREN: ")",
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for tok := operatorBeg + 1; tok < operatorEnd; tok++ {
		keywords[strings.ToLower(tokens[tok])] = tok
	}

	for tok := keywordBeg + 1; tok < keywordEnd; tok++ {
		keywords[strings.ToLower(tokens[tok])] = tok
	}

	for tok := moyBeg + 1; tok < moyEnd; tok++ {
		keywords[strings.ToLower(tokens[tok])] = tok
	}

	for tok := dowBeg + 1; tok < dowEnd; tok++ {
		keywords[strings.ToLower(tokens[tok])] = tok
	}
}

// String returns the string representation of the token
func (tok Token) String() string {
	if tok >= 0 && tok < Token(len(tokens)) {
		return tokens[tok]
	}
	return ""
}

// isMonthOfYear returns true for month of year tokens
func (tok Token) isMonthOfYear() bool {
	return tok > moyBeg && tok < moyEnd
}

// isDayOfWeek returns true for day of week tokens
func (tok Token) isDayOfWeek() bool {
	return tok > dowBeg && tok < dowEnd
}

// tokstr returns the token string or the literal if available.
func tokstr(tok Token, lit string) string {
	if lit != "" {
		return lit
	}
	return tok.String()
}

// Lookup returns the token associated with a given string
func Lookup(ident string) Token {
	if tok, ok := keywords[strings.ToLower(ident)]; ok {
		return tok
	}
	return IDENT
}

// Pos specifies the line and character position of a token.
// The Char and Line are both zero-based indices.
type Pos struct {
	Line int
	Char int
}

// String returns the string representation of the position
func (p Pos) String() string {
	return fmt.Sprintf("%d col %d", p.Line+1, p.Char+1)
}

// getMonthOfYear converts the month token into a time.Month value
func getMonthOfYear(tok Token) time.Month {
	switch tok {
	case JANUARY:
		return time.January
	case FEBRUARY:
		return time.February
	case MARCH:
		return time.March
	case APRIL:
		return time.April
	case MAY:
		return time.May
	case JUNE:
		return time.June
	case JULY:
		return time.July
	case AUGUST:
		return time.August
	case SEPTEMBER:
		return time.September
	case OCTOBER:
		return time.October
	case NOVEMBER:
		return time.November
	case DECEMBER:
		return time.December
	default:
		return -1
	}
}

// getDayOfWeek converts the day token into a time.Weekday value
func getDayOfWeek(tok Token) time.Weekday {
	switch tok {
	case MONDAY:
		return time.Monday
	case TUESDAY:
		return time.Tuesday
	case WEDNESDAY:
		return time.Wednesday
	case THURSDAY:
		return time.Thursday
	case FRIDAY:
		return time.Friday
	case SATURDAY:
		return time.Saturday
	case SUNDAY:
		return time.Sunday
	default:
		return -1
	}
}

package timewarp_test

import (
	"bytes"

	. "github.com/iancmcc/fasterstronger/timerange"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scanner", func() {
	var (
		s   *Scanner
		buf bytes.Buffer
	)

	JustBeforeEach(func() {
		s = NewScanner(&buf)
	})

	DescribeTable("Identify tokens",
		func(str string, tok Token, lit string) {
			_, _ = buf.WriteString(str)
			t, _, l := s.Scan()
			Expect(t).To(Equal(tok))
			Expect(l).To(Equal(lit))
		},
		Entry("EOF", ``, EOF, ``),
		Entry("ILLEGAL", `#`, ILLEGAL, `#`),
		Entry("WS", "\n \r\n", WS, "\n \n"),
		Entry("WS", "\n \r", WS, "\n \n"),
		Entry("WS", "\n \r ", WS, "\n \n "),

		Entry("IDENT <1st>", `1st`, IDENT, `1st`),
		Entry("IDENT <ms>", `ms`, IDENT, `ms`),

		Entry("AND", "and", AND, ""),
		Entry("IN", "in", IN, ""),
		Entry("OF", "of", OF, ""),
		Entry("NOT", "not", NOT, ""),

		Entry("YEAR", "year", YEAR, ""),
		Entry("MONTH", "month", MONTH, ""),
		Entry("WEEK", "week", WEEK, ""),
		Entry("DAY", "day", DAY, ""),
		Entry("TIME", "time", TIME, ""),

		Entry("JANUARY", "january", JANUARY, ""),
		Entry("FEBRUARY", "february", FEBRUARY, ""),
		Entry("MARCH", "march", MARCH, ""),
		Entry("APRIL", "april", APRIL, ""),
		Entry("MAY", "may", MAY, ""),
		Entry("JUNE", "june", JUNE, ""),
		Entry("JULY", "july", JULY, ""),
		Entry("AUGUST", "august", AUGUST, ""),
		Entry("SEPTEMBER", "september", SEPTEMBER, ""),
		Entry("OCTOBER", "october", OCTOBER, ""),
		Entry("NOVEMBER", "november", NOVEMBER, ""),
		Entry("DECEMEBER", "december", DECEMBER, ""),

		Entry("MONDAY", "monday", MONDAY, ""),
		Entry("TUESDAY", "tuesday", TUESDAY, ""),
		Entry("WEDNESDAY", "wednesday", WEDNESDAY, ""),
		Entry("THURSDAY", "thursday", THURSDAY, ""),
		Entry("FRIDAY", "friday", FRIDAY, ""),
		Entry("SATURDAY", "saturday", SATURDAY, ""),
		Entry("SUNDAY", "sunday", SUNDAY, ""),

		Entry("LPAREN", "(", LPAREN, ""),
		Entry("RPAREN", ")", RPAREN, ""),
	)

	Describe("Sentence", func() {
		BeforeEach(func() {
			_, _ = buf.WriteString("DAY TUESDAY AND DAY\n WEDNESDAY AND DAY FRIDAY SUNDAY")
		})

		ExpectScanned := func(tok Token, pos Pos, lit string) {
			t, p, l := s.Scan()
			Expect(tok).To(Equal(t))
			Expect(pos).To(Equal(p))
			Expect(lit).To(Equal(l))
		}

		Specify("tokens scanned", func() {
			ExpectScanned(DAY, Pos{0, 0}, "")
			ExpectScanned(WS, Pos{0, 3}, " ")
			ExpectScanned(TUESDAY, Pos{0, 4}, "")
			ExpectScanned(WS, Pos{0, 11}, " ")
			ExpectScanned(AND, Pos{0, 12}, "")
			ExpectScanned(WS, Pos{0, 15}, " ")
			ExpectScanned(DAY, Pos{0, 16}, "")
			ExpectScanned(WS, Pos{0, 19}, "\n ")
			ExpectScanned(WEDNESDAY, Pos{1, 1}, "")
			ExpectScanned(WS, Pos{1, 10}, " ")
			ExpectScanned(AND, Pos{1, 11}, "")
			ExpectScanned(WS, Pos{1, 14}, " ")
			ExpectScanned(DAY, Pos{1, 15}, "")
			ExpectScanned(WS, Pos{1, 18}, " ")
			ExpectScanned(FRIDAY, Pos{1, 19}, "")
			ExpectScanned(WS, Pos{1, 25}, " ")
			ExpectScanned(SUNDAY, Pos{1, 26}, "")
			ExpectScanned(EOF, Pos{1, 32}, "")
			ExpectScanned(EOF, Pos{1, 32}, "")
			ExpectScanned(EOF, Pos{1, 32}, "")
		})

	})
})

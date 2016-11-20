package timewarp_test

import (
	"bytes"
	"time"

	. "github.com/FasterStronger/timewarp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parser", func() {
	const (
		datefmt = "01-02-06"
		timefmt = "1504"
	)

	var (
		r      TimeRange
		in     string
		out    Filter
		result Filter
		err    error
	)

	BeforeEach(func() {
		r, _ = Parse(datefmt, "01-01-2006", "01-01-2017")
	})

	JustBeforeEach(func() {
		result, err = NewParser(bytes.NewBufferString(in)).Parse()
	})

	AssertFilter := func() {
		Specify("a matching filter", func() {
			Expect(err).To(BeNil())
			Expect(result(r)).To(Equal(out(r)))
		})
	}

	AssertError := func() {
		Specify("an error", func() {
			Expect(err).ToNot(BeNil())
			Expect(result).To(BeNil())
		})
	}

	Context("Not a number", func() {
		BeforeEach(func() {
			in = `DAY A 1 OF MONTH JUNE IN YEAR 2007`
		})
		AssertError()
	})

	Context("June 5th 2007", func() {
		BeforeEach(func() {
			in = `DAY 5 OF MONTH JUNE IN YEAR 2007`
			out = Days(5, 1).Filter().Of(1, TheMonth(time.June)).In(Year(2007))
		})
		AssertFilter()
	})

	Context("Leap days", func() {
		BeforeEach(func() {
			in = `DAY 29 OF MONTH FEBRUARY`
			out = Days(29, 1).Filter().Of(1, TheMonth(time.February))
		})
		AssertFilter()
	})

	Context("The second Tuesday of the month", func() {
		BeforeEach(func() {
			in = `DAY TUESDAY OF 2 MONTH`
			out = Week(time.Tuesday, 1).Filter().Of(2, TheMonth(0))
		})
		AssertFilter()
	})

	Context("Mondays-Wednesdays, and Fridays from 4-6p", func() {
		BeforeEach(func() {
			in = `DAY MONDAY WEDNESDAY AND DAY FRIDAY IN TIME 1600 1800`
			out = Week(time.Monday, 3).Filter().And(Week(time.Friday, 1)).In(Times(timefmt, "1600", "1800"))
		})
		AssertFilter()
	})

	Context("Sundays from 8-10a, Tuesdays from 4-9p", func() {
		BeforeEach(func() {
			in = `(DAY SUNDAY IN TIME 0800 1000) AND (DAY TUESDAY IN TIME 1600 2100)`
			f1 := Week(time.Sunday, 1).Filter().In(Times(timefmt, "0800", "1000"))
			f2 := Week(time.Tuesday, 1).Filter().In(Times(timefmt, "1600", "2100"))
			out = f1.Union(f2)
		})
		AssertFilter()
	})

	Context("Dangling AND", func() {
		BeforeEach(func() {
			in = `DAY SUNDAY AND`
		})
		AssertError()
	})

	Context("Invalid time format 1st arg", func() {
		BeforeEach(func() {
			in = `TIME 9999 1234`
		})
		AssertError()
	})

	Context("Invald time format 2nd arg", func() {
		BeforeEach(func() {
			in = `TIME 1234 9999`
		})
		AssertError()
	})

	Context("Invalid time format missing 1st arg", func() {
		BeforeEach(func() {
			in = `TIME AND`
		})
		AssertError()
	})

	Context("Invalid time format missing second arg", func() {
		BeforeEach(func() {
			in = `TIME 1000`
		})
		AssertError()
	})

	Context("Invalid consecutive days", func() {
		BeforeEach(func() {
			in = `DAY 4 GHOST`
		})
		AssertError()
	})

	Context("Every three days", func() {
		BeforeEach(func() {
			in = `DAY OF 3 DAY`
			out = Days(0, 1).Filter().Of(3, TheDays(-2, 5))
		})
		AssertFilter()
	})

	Context("Invalid day ordinal", func() {
		BeforeEach(func() {
			in = `DAY OF DAY TUESDAY`
		})
		AssertError()
	})

	Context("Three days on, 2 days off", func() {
		BeforeEach(func() {
			in = `DAY 3 OF DAY 5`
			out = Days(0, 3).Filter().Of(1, TheDays(0, 5))
		})
		AssertFilter()
	})

	Context("First Monday from three days from now", func() {
		BeforeEach(func() {
			in = `DAY MONDAY OF DAY 3 7`
			out = Week(time.Monday, 1).Filter().Of(1, TheDays(3, 7))
		})
		AssertFilter()
	})

	Context("The week, starting from Monday", func() {
		BeforeEach(func() {
			in = `WEEK MONDAY`
			out = Week(time.Monday, 7).Filter()
		})
		AssertFilter()
	})

	Context("Tuesday, Wednesday, every 3 weeks", func() {
		BeforeEach(func() {
			in = `DAY TUESDAY WEDNESDAY OF 3 WEEK MONDAY`
			out = Week(time.Tuesday, 2).Filter().Of(3, TheWeek(time.Monday, 7, -2, 5))
		})
		AssertFilter()
	})

	Context("Tuesday, Wednesday every 3 weeks", func() {
		BeforeEach(func() {
			in = `DAY TUESDAY WEDNESDAY OF 3 WEEK`
			out = Week(time.Tuesday, 2).Filter().Of(3, TheWeek(-1, 7, -2, 5))
		})
		AssertFilter()
	})

	Context("Fourth Thursday of the month", func() {
		BeforeEach(func() {
			in = `DAY THURSDAY OF 4 MONTH`
			out = Week(time.Thursday, 1).Filter().Of(4, Month(0))
		})
		AssertFilter()
	})

	Context("June", func() {
		BeforeEach(func() {
			in = `MONTH JUNE`
			out = Month(time.June).Filter()
		})
		AssertFilter()
	})

	Context("Invalid ordinal", func() {
		BeforeEach(func() {
			in = `DAY 6 1 OF YEAR 2016`
		})
		AssertError()
	})

	Context("Missing year", func() {
		BeforeEach(func() {
			in = `YEAR`
		})
		AssertError()
	})

	Context("Can't parse year", func() {
		BeforeEach(func() {
			in = `YEAR hat`
		})
		AssertError()
	})

	Context("Invalid year", func() {
		BeforeEach(func() {
			in = `YEAR 0`
		})
		AssertError()
	})

	Context("Invalid IN", func() {
		BeforeEach(func() {
			in = `DAY FRIDAY IN TUESDAY`
		})
		AssertError()
	})

	Context("Invalid ordinal number", func() {
		BeforeEach(func() {
			in = `DAY TUESDAY OF X MONTH`
		})
		AssertError()
	})

	Context("Invalid ordinal value", func() {
		BeforeEach(func() {
			in = `DAY TUESDAY OF 0 MONTH`
		})
		AssertError()
	})

	Context("Bad operator", func() {
		BeforeEach(func() {
			in = `DAY TUESDAY MONTH MAY`
		})
		AssertError()
	})

	Context("Except the first Tuesday of March", func() {
		BeforeEach(func() {
			in = `NOT (DAY TUESDAY OF MONTH MARCH)`
			out = Week(time.Tuesday, 1).Filter().Of(1, TheMonth(time.March))
		})
		AssertFilter()
	})

	Context("Missing close paren", func() {
		BeforeEach(func() {
			in = `NOT (DAY TUESDAY`
		})
		AssertError()
	})

	Context("Bad expression in left paren", func() {
		BeforeEach(func() {
			in = `NOT (YEAR`
		})
		AssertError()
	})
})

package timewarp_test

import (
	"bytes"
	"time"

	. "github.com/takeinitiative/timewarp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parser", func() {
	const (
		datefmt = "01-02-06"
		timefmt = "1504"
	)

	var (
		r      *TimeRange
		in     string
		out    Filter
		result Filter
		err    error
	)

	BeforeEach(func() {
		r, err = Parse(datefmt, "01-01-18", "01-01-19")
		Expect(err).ShouldNot(HaveOccurred())
	})

	JustBeforeEach(func() {
		result, err = NewParser(bytes.NewBufferString(in)).Parse()
	})

	AssertFilter := func() {
		Specify("a matching filter", func() {
			Expect(err).To(BeNil())
			Expect(result(*r)).To(Equal(out(*r)))
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

	Context("First Monday, Wednesday and Friday of July", func() {
		BeforeEach(func() {
			in = `(day monday friday of month july) in not (day tuesday and day thursday)`
			d1 := Week(time.Tuesday, 1).And(Week(time.Thursday, 1)).Negate()
			out = Week(time.Monday, 5).Of(1, TheMonth(time.July)).Intersect(d1)
		})
		AssertFilter()
	})

	Context("July 2, 4, 6", func() {
		BeforeEach(func() {
			in = `(DAY 2 OF MONTH JULY) and (DAY 4 OF MONTH JULY) and (DAY 6 oF MONTH JULY)`

			d1 := Days(1, 1).Of(1, TheMonth(time.July))
			d2 := Days(3, 1).Of(1, TheMonth(time.July))
			d3 := Days(5, 1).Of(1, TheMonth(time.July))

			out = d1.Union(d2, d3)
		})
		AssertFilter()
	})

	Context("July 7 830a-1130p", func() {
		BeforeEach(func() {
			in = `day 7 of month july in time 0830 1130`
			out = Days(6, 1).Of(1, TheMonth(time.July)).In(Times(timefmt, "0830", "1130"))
		})
		AssertFilter()
	})

	Context("Mondays and Saturdays, Tuesdays and Thursdays after 1p, Wednesdays and Fridays before 2p, but not the first week in july", func() {
		BeforeEach(func() {
			in = `(day monday and day saturday) and ((day tuesday and day thursday) in time 1300 0000) and ((day wednesday and day friday) in time 0000 1400) in not (week monday of month july)`

			d1 := Week(time.Monday, 1).And(Week(time.Saturday, 1))
			d2 := Week(time.Tuesday, 1).And(Week(time.Thursday, 1)).In(Times(timefmt, "1300", "0000"))
			d3 := Week(time.Wednesday, 1).And(Week(time.Friday, 1)).In(Times(timefmt, "0000", "1400"))
			d4 := Week(time.Monday, 7).Of(1, TheMonth(time.July)).Negate()
			out = d1.Union(d2, d3).Intersect(d4)
		})
		AssertFilter()
	})

	Context("Sundays through Tuesdays, Wednesday through Saturdays except Thursdays before 4p, excluding July 10-17", func() {
		BeforeEach(func() {
			in = `day sunday tuesday and ((day wednesday saturday in not day thursday) in time 0000 1600) in not (day 10 17 of month july)`
			d1 := Week(time.Sunday, 3).Filter()
			d2 := Week(time.Wednesday, 4).Filter().Intersect(Week(time.Thursday, 1).Not()).In(Times(timefmt, "0000", "1600"))
			d3 := Days(9, 8).Of(1, TheMonth(time.July)).Negate()
			out = d1.Union(d2).Intersect(d3)
		})
		AssertFilter()
	})

	Context("June 5th 2006", func() {
		BeforeEach(func() {
			in = `DAY 5 OF MONTH JUNE IN YEAR 2006`
			out = Days(4, 1).Of(1, TheMonth(time.June)).In(Year(2006))
		})
		AssertFilter()
	})

	Context("Leap days", func() {
		BeforeEach(func() {
			in = `DAY 29 OF MONTH FEBRUARY`
			out = Days(29, 1).Of(1, TheMonth(time.February))
		})
		AssertFilter()
	})

	Context("The second Tuesday of the month", func() {
		BeforeEach(func() {
			in = `DAY TUESDAY OF 2 MONTH`
			out = Week(time.Tuesday, 1).Of(2, TheMonth(0))
		})
		AssertFilter()
	})

	Context("Mondays-Wednesdays, and Fridays from 4-6p", func() {
		BeforeEach(func() {
			in = `DAY MONDAY WEDNESDAY AND DAY FRIDAY IN TIME 1600 1800`
			out = Week(time.Monday, 3).And(Week(time.Friday, 1)).In(Times(timefmt, "1600", "1800"))
		})
		AssertFilter()
	})

	Context("Sundays from 8-10a, Tuesdays from 4-9p", func() {
		BeforeEach(func() {
			in = `(DAY SUNDAY IN TIME 0800 1000) AND (DAY TUESDAY IN TIME 1600 2100)`
			f1 := Week(time.Sunday, 1).In(Times(timefmt, "0800", "1000"))
			f2 := Week(time.Tuesday, 1).In(Times(timefmt, "1600", "2100"))
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
			out = Days(0, 1).Of(3, TheDays(-2, 5))
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
			in = `DAY 1 3 OF DAY 0 5`
			out = Days(0, 3).Of(1, TheDays(0, 5))
		})
		AssertFilter()
	})

	Context("First Monday from three days from now", func() {
		BeforeEach(func() {
			in = `DAY MONDAY OF DAY 3 7`
			out = Week(time.Monday, 1).Of(1, TheDays(3, 7))
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
			out = Week(time.Tuesday, 2).Of(3, TheWeek(time.Monday, 7, -2, 5))
		})
		AssertFilter()
	})

	Context("Tuesday, Wednesday every 3 weeks", func() {
		BeforeEach(func() {
			in = `DAY TUESDAY WEDNESDAY OF 3 WEEK`
			out = Week(time.Tuesday, 2).Of(3, TheWeek(-1, 7, -2, 5))
		})
		AssertFilter()
	})

	Context("Fourth Thursday of the month", func() {
		BeforeEach(func() {
			in = `DAY THURSDAY OF 4 MONTH`
			out = Week(time.Thursday, 1).Of(4, Month(0))
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
			out = (Week(time.Tuesday, 1).Of(1, TheMonth(time.March))).Negate()
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

package timerange_test

import (
	"time"

	. "github.com/iancmcc/fasterstronger/timerange"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Query", func() {
	const datefmt = "01-02-06"

	var (
		q      Query
		in     TimeRange
		out    TimeRange
		result TimeRange
		ok     bool
	)

	AssertNotInRange := func() {
		It("should not be in range", func() {
			Expect(ok).To(BeFalse())
		})
	}

	AssertInRangeEquals := func() {
		It("should equal the input", func() {
			Expect(ok).To(BeTrue())
			Expect(out).To(Equal(in))
		})
	}

	AssertInRange := func() {
		It("should be in range", func() {
			Expect(ok).To(BeTrue())
			Expect(out).To(Equal(result))
		})
	}

	JustBeforeEach(func() {
		out, ok = q(in)
	})

	Describe("Year", func() {

		BeforeEach(func() {
			in, _ = Parse(datefmt, "05-07-13", "12-29-16")
		})

		Context("The year is earlier than the set range", func() {
			BeforeEach(func() {
				q = Year(2012)
			})
			AssertNotInRange()
		})

		Context("The year is later than the set range", func() {
			BeforeEach(func() {
				q = Year(2017)
			})
			AssertNotInRange()
		})

		Context("The year fills the range", func() {
			BeforeEach(func() {
				in, _ = Parse(datefmt, "05-07-13", "07-12-13")
				q = Year(2013)
			})
			AssertInRangeEquals()
		})

		Context("The year is a left split on the range", func() {
			BeforeEach(func() {
				result, _ = Parse(datefmt, "05-07-13", "01-01-14")
				q = Year(2013)
			})
			AssertInRange()
		})

		Context("The year is a right split on the range", func() {
			BeforeEach(func() {
				result, _ = Parse(datefmt, "01-01-16", "12-29-16")
				q = Year(2016)
			})
			AssertInRange()
		})

		Context("The year is a subset of the range", func() {
			BeforeEach(func() {
				result, _ = Parse(datefmt, "01-01-14", "01-01-15")
				q = Year(2014)
			})
			AssertInRange()
		})
	})

	Describe("Month", func() {

		BeforeEach(func() {
			in, _ = Parse(datefmt, "05-07-13", "07-12-13")
		})

		Context("The month is earlier than the set range", func() {
			BeforeEach(func() {
				q = Month(time.February)
			})
			AssertNotInRange()
		})

		Context("The month is later than the set range", func() {
			BeforeEach(func() {
				q = Month(time.December)
			})
			AssertNotInRange()
		})

		Context("Left split on the month", func() {
			BeforeEach(func() {
				q = Month(time.May)
				result, _ = Parse(datefmt, "05-07-13", "06-01-13")
			})
			AssertInRange()
		})

		Context("Right split on the month", func() {
			BeforeEach(func() {
				q = Month(time.July)
				result, _ = Parse(datefmt, "07-01-13", "07-12-13")
			})
			AssertInRange()
		})

		Context("The month is completely in range", func() {
			BeforeEach(func() {
				q = Month(time.June)
				result, _ = Parse(datefmt, "06-01-13", "07-01-13")
			})
			AssertInRange()
		})

		Context("The month is the range", func() {
			BeforeEach(func() {
				in, _ = Parse(datefmt, "05-07-13", "05-20-13")
				q = Month(time.May)
			})
			AssertInRangeEquals()
		})
	})

	Describe("TheMonth", func() {

		BeforeEach(func() {
			in, _ = Parse(datefmt, "05-07-13", "07-12-13")
		})

		Context("The first month", func() {
			BeforeEach(func() {
				q = TheMonth(0)
				result, _ = Parse(datefmt, "05-01-13", "06-01-13")
			})
			AssertInRange()
		})

		Context("The month is earlier than the set range", func() {
			BeforeEach(func() {
				q = TheMonth(time.February)
			})
			AssertNotInRange()
		})

		Context("The month is later than the set range", func() {
			BeforeEach(func() {
				q = TheMonth(time.December)
			})
			AssertNotInRange()
		})

		Context("Left split on the month", func() {
			BeforeEach(func() {
				q = TheMonth(time.May)
				result, _ = Parse(datefmt, "05-01-13", "06-01-13")
			})
			AssertInRange()
		})

		Context("Right split on the month", func() {
			BeforeEach(func() {
				q = TheMonth(time.July)
				result, _ = Parse(datefmt, "07-01-13", "08-01-13")
			})
			AssertInRange()
		})

		Context("The month is completely in range", func() {
			BeforeEach(func() {
				q = TheMonth(time.June)
				result, _ = Parse(datefmt, "06-01-13", "07-01-13")
			})
			AssertInRange()
		})
	})

	Describe("Week", func() {

		BeforeEach(func() {
			// Monday, Tuesday, Wednesday
			in, _ = Parse(datefmt, "11-07-16", "11-10-16")
		})

		Context("The week is earlier than the set range", func() {
			BeforeEach(func() {
				q = Week(time.Sunday, 1)
			})
			AssertNotInRange()
		})

		Context("The week is later than the set range", func() {
			BeforeEach(func() {
				q = Week(time.Thursday, 2)
			})
			AssertNotInRange()
		})

		Context("Left split on the week", func() {
			BeforeEach(func() {
				q = Week(time.Sunday, 2)
			})
			AssertNotInRange()
		})

		Context("Right split on the week", func() {
			BeforeEach(func() {
				q = Week(time.Wednesday, 7)
				result, _ = Parse(datefmt, "11-09-16", "11-10-16")
			})
			AssertInRange()
		})

		Context("The week is completely in range", func() {
			BeforeEach(func() {
				q = Week(time.Tuesday, 1)
				result, _ = Parse(datefmt, "11-08-16", "11-09-16")
			})
			AssertInRange()
		})

		Context("The week is the range", func() {
			BeforeEach(func() {
				q = Week(time.Monday, 3)
			})
			AssertInRangeEquals()
		})
	})

	Describe("TheWeek", func() {
		BeforeEach(func() {
			// Monday, Tuesday, Wednesday
			in, _ = Parse(datefmt, "11-07-16", "11-10-16")
		})

		Context("The week is earlier than the set range", func() {
			BeforeEach(func() {
				q = TheWeek(time.Monday, 3, -5, 3)
			})
			AssertNotInRange()
		})

		Context("The week is later than the set range", func() {
			BeforeEach(func() {
				q = TheWeek(time.Monday, 3, 5, 10)
			})
			AssertNotInRange()
		})

		Context("Left split on the week", func() {
			BeforeEach(func() {
				q = TheWeek(time.Monday, 2, -1, 2)
				result, _ = Parse(datefmt, "11-05-16", "11-09-16")
			})
			AssertInRange()
		})

		Context("Right split on the week", func() {
			BeforeEach(func() {
				q = TheWeek(time.Tuesday, 2, 0, 2)
				result, _ = Parse(datefmt, "11-08-16", "11-12-16")
			})
			AssertInRange()
		})

		Context("The week is completely in range", func() {
			BeforeEach(func() {
				q = TheWeek(-1, 3, -1, 3)
				result, _ = Parse(datefmt, "11-04-16", "11-13-16")
			})
			AssertInRange()
		})
	})

	Describe("Days", func() {

		BeforeEach(func() {
			in, _ = Parse(datefmt, "11-07-16", "11-10-16")
		})

		Context("The days occur earlier than the set range", func() {
			BeforeEach(func() {
				q = Days(-5, 3)
			})
			AssertNotInRange()
		})

		Context("The days are later than the set range", func() {
			BeforeEach(func() {
				q = Days(5, 10)
			})
			AssertNotInRange()
		})

		Context("Left split on days", func() {
			BeforeEach(func() {
				q = Days(-1, 2)
				result, _ = Parse(datefmt, "11-07-16", "11-08-16")
			})
			AssertInRange()
		})

		Context("Right split on days", func() {
			BeforeEach(func() {
				q = Days(2, 5)
				result, _ = Parse(datefmt, "11-09-16", "11-10-16")
			})
			AssertInRange()
		})

		Context("The days are completely in range", func() {
			BeforeEach(func() {
				q = Days(1, 1)
				result, _ = Parse(datefmt, "11-08-16", "11-09-16")
			})
			AssertInRange()
		})

		Context("The days are the range", func() {
			BeforeEach(func() {
				q = Days(0, 3)
			})
			AssertInRangeEquals()
		})
	})

	Describe("TheDays", func() {

		BeforeEach(func() {
			in, _ = Parse(datefmt, "11-07-16", "11-10-16")
		})

		Context("The days occur earlier than the set range", func() {
			BeforeEach(func() {
				q = TheDays(-5, 3)
			})
			AssertNotInRange()
		})

		Context("The days are later than the set range", func() {
			BeforeEach(func() {
				q = TheDays(5, 10)
			})
			AssertNotInRange()
		})

		Context("Left split on days", func() {
			BeforeEach(func() {
				q = TheDays(-1, 2)
				result, _ = Parse(datefmt, "11-06-16", "11-08-16")
			})
			AssertInRange()
		})

		Context("Right split on days", func() {
			BeforeEach(func() {
				q = TheDays(2, 5)
				result, _ = Parse(datefmt, "11-09-16", "11-14-16")
			})
			AssertInRange()
		})

		Context("The days are completely in range", func() {
			BeforeEach(func() {
				q = TheDays(-2, 7)
				result, _ = Parse(datefmt, "11-05-16", "11-12-16")
			})
			AssertInRange()
		})

		Context("The days are the range", func() {
			BeforeEach(func() {
				q = TheDays(0, 3)
			})
			AssertInRangeEquals()
		})
	})

	Describe("Times", func() {
		const (
			datetimefmt = "01-02-06 3:04PM"
			timefmt     = "3:04PM"
		)

		BeforeEach(func() {
			in, _ = Parse(datetimefmt, "11-12-16 7:15PM", "11-13-16 4:10AM")
		})

		Context("The time is earlier than the set range", func() {
			BeforeEach(func() {
				q = Times(timefmt, "5:00PM", "7:00PM")
			})
			AssertNotInRange()
		})

		Context("The time is later than the set range", func() {
			BeforeEach(func() {
				q = Times(timefmt, "5:00AM", "7:00AM")
			})
			AssertNotInRange()
		})

		Context("The time fills the range", func() {
			BeforeEach(func() {
				q = Times(timefmt, "6:00PM", "6:00AM")
			})
			AssertInRangeEquals()
		})

		Context("The time is a left split on the range", func() {
			BeforeEach(func() {
				q = Times(timefmt, "6:00PM", "8:00PM")
				result, _ = Parse(datetimefmt, "11-12-16 7:15PM", "11-12-16 8:00PM")
			})
			AssertInRange()
		})

		Context("The time is a right split on the range", func() {
			BeforeEach(func() {
				q = Times(timefmt, "4:00AM", "5:00AM")
				result, _ = Parse(datetimefmt, "11-13-16 4:00AM", "11-13-16 4:10AM")
			})
			AssertInRange()
		})

		Context("The time is a subset of the range", func() {
			BeforeEach(func() {
				q = Times(timefmt, "8:00PM", "12:00AM")
				result, _ = Parse(datetimefmt, "11-12-16 8:00PM", "11-13-16 12:00AM")
			})
			AssertInRange()
		})
	})
})

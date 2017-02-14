package timewarp_test

import (
	"time"

	. "github.com/FasterStronger/timewarp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Filter", func() {
	const datefmt = "01-02-06"

	var (
		f      Filter
		in     TimeRange
		out    []TimeRange
		result []TimeRange
	)

	JustBeforeEach(func() {
		out = f(in)
	})

	Context("Filter from Query", func() {

		BeforeEach(func() {
			f = Month(time.June).Filter()
			in, _ = Parse(datefmt, "06-12-13", "06-20-15")

			slot1, _ := Parse(datefmt, "06-12-13", "07-01-13")
			slot2, _ := Parse(datefmt, "06-01-14", "07-01-14")
			slot3, _ := Parse(datefmt, "06-01-15", "06-20-15")
			result = []TimeRange{slot1, slot2, slot3}
		})

		It("should return the filtered results", func() {
			Expect(out).To(Equal(result))
		})

	})

	Context("Negate", func() {

		BeforeEach(func() {
			f = Month(time.June).Not()
			in, _ = Parse(datefmt, "06-12-13", "06-20-15")

			slot1, _ := Parse(datefmt, "07-01-13", "06-01-14")
			slot2, _ := Parse(datefmt, "07-01-14", "06-01-15")
			result = []TimeRange{slot1, slot2}
		})

		It("should return the inverse filtered results", func() {
			Expect(out).To(Equal(result))
		})
	})

	Context("Union", func() {

		BeforeEach(func() {
			mf1 := Month(time.June).Filter()
			mf2 := Month(time.July).Filter()
			mf3 := Month(time.November).Filter()
			f = mf1.Union(mf2, mf3)

			in, _ = Parse(datefmt, "11-04-13", "08-01-14")

			slot1, _ := Parse(datefmt, "11-04-13", "12-01-13")
			slot2, _ := Parse(datefmt, "06-01-14", "07-01-14")
			slot3, _ := Parse(datefmt, "07-01-14", "08-01-14")
			result = []TimeRange{slot1, slot2, slot3}
		})

		It("should return the union of the results", func() {
			Expect(out).To(ConsistOf(result))
		})
	})

	Context("Intersect", func() {

		BeforeEach(func() {
			mf := Month(time.June).Filter()
			yf := Year(2013).Filter()
			f = mf.Intersect(yf)

			in, _ = Parse(datefmt, "03-13-13", "04-10-15")

			slot1, _ := Parse(datefmt, "06-01-13", "07-01-13")
			result = []TimeRange{slot1}
		})

		It("should return an intersection of the results", func() {
			Expect(out).To(Equal(result))
		})
	})

	Context("Ordinal", func() {

		BeforeEach(func() {
			mf := TheMonth(time.November).Filter()
			df := Week(time.Thursday, 1).Filter()
			f = df.Ordinal(4, mf)

			in, _ = Parse(datefmt, "11-11-16", "11-30-16")

			slot1, _ := Parse(datefmt, "11-24-16", "11-25-16")
			result = []TimeRange{slot1}
		})

		It("should return the 4th Thursday of November", func() {
			Expect(out).To(Equal(result))
		})
	})
})

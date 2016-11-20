package timewarp_test

import (
	"time"

	. "github.com/FasterStronger/timewarp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TimeRange", func() {

	Describe("Slot", func() {
		var (
			format string
			start  string
			end    string
			slot   TimeRange
			err    error
		)

		JustBeforeEach(func() {
			slot, err = Parse(format, start, end)
		})

		Context("invalid time range", func() {
			BeforeEach(func() {
				format, start, end = "xyz", "12:00AM", "7:00AM"
			})

			It("should have an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(slot.Duration()).To(Equal(time.Duration(0)))
			})
		})

		Context("valid time range", func() {
			BeforeEach(func() {
				format, start, end = time.Kitchen, "12:00AM", "7:00AM"
			})

			It("should create a time range", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(slot.Duration()).To(Equal(7 * time.Hour))
			})
		})
	})

	Describe("Sorting and Searching", func() {
		var (
			slots []TimeRange
			slot1 TimeRange
			slot2 TimeRange
			slot3 TimeRange
			slot4 TimeRange
			slot5 TimeRange
		)

		BeforeEach(func() {
			slot1, _ = Parse(time.Kitchen, "12:00AM", "5:00AM")
			slot2, _ = Parse(time.Kitchen, "1:00AM", "5:00AM")
			slot3, _ = Parse(time.Kitchen, "12:00AM", "3:00AM")
			slot4, _ = Parse(time.Kitchen, "4:00AM", "7:00AM")
			slot5, _ = Parse(time.Kitchen, "8:00AM", "2:00PM")
			slots = []TimeRange{slot1, slot2, slot3, slot4, slot5}
		})

		Context("sorted", func() {
			JustBeforeEach(func() {
				Sort(slots)
			})

			It("should be sorted", func() {
				Expect(slots).To(Equal([]TimeRange{slot3, slot1, slot2, slot4, slot5}))
			})
		})

		Context("not found", func() {
			var (
				search1 TimeRange
				search2 TimeRange
				search3 TimeRange
			)

			BeforeEach(func() {
				search1, _ = Parse(time.Kitchen, "12:00AM", "2:00PM")
				search2, _ = Parse(time.Kitchen, "2:00AM", "6:00AM")
				search3, _ = Parse(time.Kitchen, "6:00AM", "10:00AM")
			})

			It("should not find a time slot", func() {
				Expect(SearchIndex(slots, search1)).To(Equal(-1))
				Expect(SearchIndex(slots, search2)).To(Equal(-1))
				Expect(SearchIndex(slots, search3)).To(Equal(-1))
			})
		})

		Context("found", func() {
			var (
				search1 TimeRange
				search2 TimeRange
			)

			BeforeEach(func() {
				search1, _ = Parse(time.Kitchen, "12:00AM", "3:00AM")
				search2, _ = Parse(time.Kitchen, "12:00AM", "5:00AM")
				Sort(slots)
			})

			It("should find time slots", func() {
				Expect(slots[SearchIndex(slots, search1)]).To(Equal(slot3))
				Expect(slots[SearchIndex(slots, search2)]).To(Equal(slot1))
			})
		})

	})

	Describe("Merging", func() {
		var slots []TimeRange

		BeforeEach(func() {
			slot1, _ := Parse(time.Kitchen, "2:00PM", "4:00PM")
			slot2, _ := Parse(time.Kitchen, "12:00PM", "5:00PM")
			slot3, _ := Parse(time.Kitchen, "6:00PM", "9:00PM")
			slot4, _ := Parse(time.Kitchen, "9:00PM", "10:00PM")
			slot5, _ := Parse(time.Kitchen, "10:00AM", "3:00PM")
			slots = []TimeRange{slot1, slot2, slot3, slot4, slot5}
		})

		JustBeforeEach(func() {
			Merge(&slots)
		})

		It("should merge overlapping time slots", func() {
			slot1, _ := Parse(time.Kitchen, "10:00AM", "5:00PM")
			slot2, _ := Parse(time.Kitchen, "6:00PM", "10:00PM")

			Expect(slots).To(Equal([]TimeRange{slot1, slot2}))
		})
	})
})

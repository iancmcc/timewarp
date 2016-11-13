package timerange

// Query is a function that finds the first matching slot in a time range.
type Query func(input TimeRange) (output TimeRange, ok bool)

// Filter changes the query into a filter.
func (q Query) Filter() Filter {
	return func(input TimeRange) []TimeRange {
		var result []TimeRange

		for input.Duration() > 0 {
			var output, ok = q(input)
			if !ok {
				break
			}

			result = append(result, output)
			input.Start = output.End
		}

		return result
	}
}

// Filter is a function that returns all matching slots in a time range.
type Filter func(input TimeRange) []TimeRange

// Not returns a filter that returns the inverse results
func (f Filter) Not() Filter {
	return func(input TimeRange) []TimeRange {
		var result []TimeRange

		for _, s := range f(input) {
			if input.Start.Before(s.Start) {
				result = append(result, TimeRange{input.Start, s.Start})
			}
			input.Start = s.End
		}

		if input.Start.Before(input.End) {
			result = append(result, input)
		}

		return result
	}
}

// Union returns a filter that's result comprises of multiple filters
func (f Filter) Union(filters ...Filter) Filter {
	return func(input TimeRange) []TimeRange {
		var result = f(input)

		for _, f := range filters {
			result = append(result, f(input)...)
		}

		Merge(&result)
		return result
	}
}

// And is same as Union, but passes a query instead of a filter
func (f Filter) And(queries ...Query) Filter {
	var filters []Filter
	for _, q := range queries {
		filters = append(filters, q.Filter())
	}
	return f.Union(filters...)
}

// Intersect returns a filter that's result must satisfy all filters
func (f Filter) Intersect(filters ...Filter) Filter {
	return func(input TimeRange) []TimeRange {
		var result = f(input)

		for _, f := range filters {
			var output []TimeRange

			for _, s := range result {
				output = append(output, f(s)...)
			}

			result = output
		}

		return result
	}
}

// In is the same as Intersect but passes a query instead of a filter
func (f Filter) In(queries ...Query) Filter {
	var filters []Filter
	for _, q := range queries {
		filters = append(filters, q.Filter())
	}
	return f.Intersect(filters...)
}

// Ordinal returns a filter of ranges within the ordinal range
func (f Filter) Ordinal(order int, filter Filter) Filter {
	if order == 0 {
		panic("ordinal cannot be zero")
	}

	return func(input TimeRange) (result []TimeRange) {
		for _, v := range filter(input) {
			var r = f(v)

			// find the range that satisfies the ordinal
			var output TimeRange
			if size := len(r); order < 0 {
				if -order > size {
					continue
				}
				output = r[order+size]
			} else {
				if order > size {
					continue
				}
				output = r[order-1]
			}

			// continue if the objective value exists, but out of scope
			if !output.Start.Before(input.End) || !output.End.After(input.Start) {
				continue
			}

			// adjust start and end to meet input criteria
			if output.Start.Before(input.Start) {
				output.Start = input.Start
			}
			if output.End.After(input.End) {
				output.End = input.End
			}
			result = append(result, output)
		}
		return
	}
}

// Of is the same as Ordinal, but passes a query instead of a filter
func (f Filter) Of(order int, q Query) Filter {
	return f.Ordinal(order, q.Filter())
}

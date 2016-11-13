# TimeWarp
## Time ranges from natural language
TimeWarp makes describing  recurring time ranges easy!  It has **no external dependencies** and is useful for articulating dates, times, and recurrances for natural language processing.

## Usage
```go
import (
    "bytes"
    "github.com/FasterStronger/timewarp"
)

func main() {
    // initialize the range to filter from
    in := timewarp.TimeRange{
        Start: time.Now().AddDate(0, -1, 0), 
        End: time.Now().AddDate(0, 1, 0),
    }
    
    // define the recurrance that we want to parse (every Tuesday)
    p := timewarp.NewParser(bytes.NewParser("DAY TUESDAY"))
    
    // get the filter function
    filter, err := p.Parse()
    if err != nil {
        panic(err)
    }
    
    // print all the time ranges
    times := filter(in)
    for _, t := range times {
        fmt.Println(t.Start, t.End)
    }
}
```

## Parser Syntax
TimeWarp uses a basis syntax parser to procedurally generate functions that will find all time ranges that apply to the input timerange.

### Example: The second Tuesday of March from 12-2p
Syntax: `DAY TUESDAY OF 2 MONTH MARCH IN TIME 1200 1400`

### Example: Fridays through Sundays and Wednesdays
Syntax: `DAY FRIDAY SUNDAY AND WEDNESDAY`

### Example: July 15, 2008
Syntax: `DAY 15 1 OF MONTH JULY IN YEAR 2008`

### Example: Weekdays from 5-11a except Tuesdays
Syntax: `DAY MONDAY FRIDAY IN TIME 0500 1100 AND NOT DAY TUESDAY`

### Example: Every other Saturday
Syntax: `DAY SATURDAY OF 2 WEEK SATURDAY`

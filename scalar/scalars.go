package scalar

import (
	"fmt"
	"time"
	"io"
	"strconv"
)

const DateFormat = "2006-01-02"
// const DateTimeFormat = "2006-01-02 15:04:00"

type CustomDate time.Time
type CustomDateTime time.Time

func (t CustomDate) MarshalGQL(w io.Writer) {
	goTime := time.Time(t)
	dateStr := goTime.Format(DateFormat)
	io.WriteString(w, strconv.Quote(dateStr))
}


func (t *CustomDate) UnmarshalGQL(v any) error {
	dateStr, ok := v.(string)
	if !ok {
		return fmt.Errorf("date must be a string")
	}
	parsedTime, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return fmt.Errorf("date must be in YYYY-MM-dd format: &%w", err)
	}
	*t = CustomDate(parsedTime)
	return nil
}



package messaging

import (
	"fmt"
	"strings"
	"time"
)

const (
	timeFormat = time.RFC3339Nano
)

func NewTime(t time.Time) CustomTime {
	return CustomTime{Time: t}
}

func ParseTime(s string) (*time.Time, error) {
	parse, err := time.Parse(timeFormat, s)
	if err != nil {
		return nil, err
	}
	return &parse, nil
}

type CustomTime struct {
	time.Time
}
func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(timeFormat, s)
	return
}

var nilTime = (time.Time{}).UnixNano()

func (ct *CustomTime) MarshalJSON() ([]byte, error) {
	if ct.Time.UnixNano() == nilTime {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(timeFormat))), nil
}

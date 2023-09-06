package marshal

import (
	"strings"
	"time"
)

type CustomTime struct {
	time.Time
}

func (ct CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse("2006-01-02", s)
	return
}

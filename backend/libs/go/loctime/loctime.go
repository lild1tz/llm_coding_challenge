package loctime

import (
	"time"
)

var loc *time.Location

func init() {
	loc, _ = time.LoadLocation("Europe/Moscow")
}

func Transfer(t time.Time) time.Time {
	return t.In(loc)
}

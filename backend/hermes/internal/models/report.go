package models

import (
	"time"

	"github.com/lild1tz/llm_coding_challenge/backend/libs/go/loctime"
)

type Report struct {
	ID            int
	ChatContextID int
	StartedAt     time.Time
	LastUpdatedAt time.Time
	FinishedAt    *time.Time
}

func (r *Report) IsFinished() bool {
	return r.FinishedAt != nil
}

func (r *Report) IsNeedToFinish(finishHour int) bool {
	now := loctime.Transfer(time.Now())
	now = now.Truncate(time.Hour * 24)

	return now.Hour() >= finishHour
}

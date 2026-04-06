package task

import "time"

type Status string
type RecurrenceType string

const (
	StatusNew        Status         = "new"
	StatusInProgress Status         = "in_progress"
	StatusDone       Status         = "done"
	Daily            RecurrenceType = "daily"
	Monthly          RecurrenceType = "monthly"
	Specific         RecurrenceType = "specific"
	EvenDays         RecurrenceType = "even"
	OddDays          RecurrenceType = "odd"
)

type Task struct {
	ID             int64          `json:"id"`
	Title          string         `json:"title"`
	Description    string         `json:"description"`
	Status         Status         `json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	RecurrenceType RecurrenceType `json:"recurrence_type"`
	RecurrenceRule string         `json:"recurrence_rule"`
}

func (s Status) Valid() bool {
	switch s {
	case StatusNew, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

package reminder

import (
	"time"
)

type Reminder struct {
	Message  string
	At       time.Duration
	Sent     bool
	timer    *time.Timer
	notifier func(string)
}

func NewReminder(message string, at time.Duration, notifier func(string)) *Reminder {
	return &Reminder{
		Message:  message,
		At:       at,
		Sent:     false,
		timer:    nil,
		notifier: notifier,
	}
}

func (r *Reminder) Send() {
	if r.Sent {
		return
	}
	r.notifier(r.Message)
	r.Sent = true
}

func (r *Reminder) Start() {
	r.timer = time.AfterFunc(r.At, r.Send)
}

func (r *Reminder) Stop() {
	r.timer.Stop()
}

package reminder

import (
	"fmt"
	"time"
)

type Reminder struct {
	Message string
	At      time.Duration
	Sent    bool
	timer   *time.Timer
}

func NewReminder(message string, at time.Duration) *Reminder {
	return &Reminder{
		Message: message,
		At:      at,
		Sent:    false,
		timer:   nil,
	}
}

func (r *Reminder) Send() {
	if r.Sent {
		return
	}
	fmt.Println("Reminder!", r.Message)
	r.Sent = true
}

func (r *Reminder) Start() {
	r.timer = time.AfterFunc(r.At, r.Send)
}

func (r *Reminder) Stop() {
	r.timer.Stop()
}

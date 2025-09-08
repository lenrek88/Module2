package reminder

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Reminder struct {
	Message  string
	At       time.Duration
	Sent     bool
	timer    *time.Timer
	notifier func(string)
}

var ErrEmptyMessage = errors.New("message is empty")

func NewReminder(message string, at time.Duration, notifier func(string)) (*Reminder, error) {
	if len(strings.TrimSpace(message)) == 0 {
		return nil, fmt.Errorf("can't create reminder: %w", ErrEmptyMessage)
	}

	return &Reminder{
		Message:  message,
		At:       at,
		Sent:     false,
		timer:    nil,
		notifier: notifier,
	}, nil
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

package events

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/araddon/dateparse"
	"github.com/google/uuid"
	"github.com/lenrek88/Module2/reminder"
)

type Event struct {
	ID       string
	Title    string
	StartAt  time.Time
	Priority Priority
	Reminder *reminder.Reminder
}

func getNextID() string {
	return uuid.New().String()
}

func (e *Event) UpdateEvent(title string, date string, priority string) error {

	t, err := dateparse.ParseAny(date)

	if err != nil {
		return errors.New("неверный формат даты")
	}

	if !IsValidTitle(title) {
		return errors.New("некорректное имя задачи")
	}

	p := Priority(priority)
	if err := p.Validate(); err != nil {
		return errors.New("неверный приоритет")
	}

	e.Title = title
	e.StartAt = t
	e.Priority = p
	return nil
}

func (e *Event) AddReminder(message string, at time.Time, notifier func(string)) error {
	startTime := at
	endTime := time.Now()
	duration := startTime.Sub(endTime)

	r, err := reminder.NewReminder(message, duration, notifier)
	if err != nil {
		return fmt.Errorf("can't add reminder to event: %w", err)
	}
	e.Reminder = r
	e.Reminder.Start()
	return nil
}

func (e *Event) RemoveReminder() error {
	if e.Reminder == nil {
		return errors.New("reminder is nil")
	}
	e.Reminder.Stop()
	e.Reminder = nil
	return nil
}

func NewEvent(title string, dateStr string, priority string) (*Event, error) {
	t, err := dateparse.ParseAny(dateStr)

	if err != nil {
		return &Event{}, errors.New("неверный формат даты")
	}

	if !IsValidTitle(title) {
		return &Event{}, errors.New("некорректное имя задачи")
	}

	p := Priority(priority)
	if err := p.Validate(); err != nil {
		return &Event{}, errors.New("неверный приоритет ")
	}

	return &Event{
		ID:       getNextID(),
		Title:    title,
		StartAt:  t,
		Priority: p,
		Reminder: nil,
	}, nil
}

func IsValidTitle(title string) bool {
	pattern := `^[А-Яа-яЁё0-9 ,\.]{3,40}$`
	matched, err := regexp.MatchString(pattern, title)
	if err != nil {
		return false
	}
	return matched
}

package events

import (
	"errors"
	"github.com/araddon/dateparse"
	"github.com/google/uuid"
	"github.com/lenrek88/app/reminder"
	"regexp"
	"time"
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

func (e *Event) AddReminder(message string, at time.Time) {
	e.Reminder = reminder.NewReminder(message, at)
}

func (e *Event) RemoveReminder() {
	e.Reminder = nil
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

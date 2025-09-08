package calendar

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/araddon/dateparse"
	"github.com/lenrek88/app/events"
	"github.com/lenrek88/app/storage"
)

type Calendar struct {
	eventsMap    map[string]*events.Event
	storage      storage.Store
	Notification chan string
}

func NewCalendar(s storage.Store) *Calendar {
	return &Calendar{
		eventsMap:    make(map[string]*events.Event),
		storage:      s,
		Notification: make(chan string),
	}
}

func (c *Calendar) Save() error {
	data, err := json.Marshal(&c.eventsMap)
	if err != nil {
		return err
	}
	err = c.storage.Save(data)
	return err
}

func (c *Calendar) Load() error {
	data, err := c.storage.Load()
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &c.eventsMap)
	return err
}

func (c *Calendar) SetEventReminder(id string, message string, dateStr string) error {
	t, err := dateparse.ParseLocal(dateStr)
	if err != nil {
		return errors.New("неверный формат даты")
	}
	e, exists := c.eventsMap[id]
	if !exists {
		return fmt.Errorf("event with key %q not found", id)
	}
	err = e.AddReminder(message, t, c.Notify)
	if err != nil {
		return fmt.Errorf("set reminder failed: %w", err)
	}
	return nil
}

func (c *Calendar) CancelEventReminder(id string) error {
	e, exists := c.eventsMap[id]
	if !exists {
		return fmt.Errorf("event with key %q not found", id)
	}
	e.RemoveReminder()
	return nil
}

func (c *Calendar) Notify(msg string) {
	c.Notification <- msg
}

func (c *Calendar) AddEvent(title string, date string, priority string) (*events.Event, error) {
	e, err := events.NewEvent(title, date, priority)
	if err != nil {
		return &events.Event{}, err
	}
	c.eventsMap[e.ID] = e
	return e, nil
}

func (c *Calendar) DeleteEvent(key string) {
	delete(c.eventsMap, key)

}

func (c *Calendar) EditEvent(id string, title string, date string, priority string) error {
	e, exists := c.eventsMap[id]
	if !exists {
		return fmt.Errorf("event with key %q not found", id)
	}
	err := e.UpdateEvent(title, date, priority)
	return err
}

func (c *Calendar) GetEvents() map[string]*events.Event {
	return c.eventsMap
}

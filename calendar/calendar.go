package calendar

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/araddon/dateparse"
	"github.com/lenrek88/Module2/events"
	"github.com/lenrek88/Module2/storage"
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

func (c *Calendar) SetEventReminder(title string, message string, dateStr string) error {
	calculatedID, err := c.GetID(title)
	if err != nil {
		return fmt.Errorf("event %q not found", title)
	}
	t, err := dateparse.ParseLocal(dateStr)
	if err != nil {
		return errors.New("неверный формат даты")
	}
	e, exists := c.eventsMap[calculatedID]
	if !exists {
		return fmt.Errorf("event with key %q not found", calculatedID)
	}
	err = e.AddReminder(message, t, c.Notify)
	if err != nil {
		return fmt.Errorf("set reminder failed: %w", err)
	}
	return nil
}

func (c *Calendar) CancelEventReminder(title string) error {
	calculatedID, err := c.GetID(title)
	if err != nil {
		return fmt.Errorf("event %q not found", title)
	}
	e, exists := c.eventsMap[calculatedID]
	if !exists {
		return fmt.Errorf("event with key %q not found", calculatedID)
	}
	err = e.RemoveReminder()
	if err != nil {
		return fmt.Errorf("remove reminder failed: %w", err)
	}
	return nil
}

func (c *Calendar) Notify(msg string) {
	c.Notification <- msg
}

func (c *Calendar) AddEvent(title string, date string, priority string) (*events.Event, error) {
	for _, event := range c.eventsMap {
		if event.Title == title {
			return &events.Event{}, errors.New("the title is repeated")
		}
	}
	e, err := events.NewEvent(title, date, priority)
	if err != nil {
		return &events.Event{}, err
	}
	c.eventsMap[e.ID] = e
	return e, nil
}

func (c *Calendar) DeleteEvent(title string) error {
	calculatedID, err := c.GetID(title)
	if err != nil {
		return fmt.Errorf("event %q not found", title)
	}
	delete(c.eventsMap, calculatedID)
	return nil

}

func (c *Calendar) GetID(title string) (string, error) {
	for _, event := range c.eventsMap {
		if event.Title == title {
			return event.ID, nil
		}
	}
	return "", errors.New("the title not found")
}

func (c *Calendar) EditEvent(title string, newTitle string, date string, priority string) error {
	calculatedID, err := c.GetID(title)
	if err != nil {
		return fmt.Errorf("event %q not found", title)
	}
	e, exists := c.eventsMap[calculatedID]

	if !exists {
		return fmt.Errorf("event with key %q not found", calculatedID)
	}
	err = e.UpdateEvent(newTitle, date, priority)
	return err
}

func (c *Calendar) GetEvents() map[string]*events.Event {
	return c.eventsMap
}

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/google/shlex"
	"github.com/lenrek88/app/calendar"
)

type Cmd struct {
	calendar *calendar.Calendar
}

func NewCmd(c *calendar.Calendar) *Cmd {
	return &Cmd{
		calendar: c,
	}
}

func (c *Cmd) Run() {
	p := prompt.New(
		c.executor,
		c.completer,
		prompt.OptionPrefix("> "),
	)
	go func() {
		for msg := range c.calendar.Notification {
			fmt.Println(msg)
		}
	}()
	p.Run()

}

func (c *Cmd) executor(input string) {
	parts, err := shlex.Split(input)
	if err != nil {
		fmt.Println("ошибка разделения строки", err)
		return
	}
	cmd := strings.ToLower(parts[0])

	switch cmd {
	case "add":
		if len(parts) < 4 {
			fmt.Println("Формат: add \"название события\" \"дата и время\" \"приоритет\"")
			return
		}

		title := parts[1]
		date := parts[2]
		priority := parts[3]

		e, err := c.calendar.AddEvent(title, date, priority)
		if err != nil {
			fmt.Println("Ошибка добавления:", err)
		} else {
			fmt.Println("Событие:", e.Title, "добавлено")
		}

	case "list":
		events := c.calendar.GetEvents()
		for _, e := range events {
			fmt.Println(e.ID, e.Title, "-", e.StartAt.Format("2006-01-02 15:04"), e.Priority)
		}

	case "remove":
		if len(parts) < 2 {
			fmt.Println("Формат: remove \"ключ\"")
			return
		}
		c.calendar.DeleteEvent(parts[1])

	case "update":
		if len(parts) < 5 {
			fmt.Println("Формат: update \"ключ события\" \"название события\" \"дата и время\" \"приоритет\"")
			return
		}
		err := c.calendar.EditEvent(parts[1], parts[2], parts[3], parts[4])
		if err != nil {
			fmt.Println("Ошибка редактирования:", err)
			return
		}

	case "reminder":
		if len(parts) < 4 {
			fmt.Println("Формат: reminder \"ключ события\" \"сообщение напоминания\" \"дата и время\"")
			return
		}
		err := c.calendar.SetEventReminder(parts[1], parts[2], parts[3])
		if err != nil {
			fmt.Println("Ошибка создания напоминания:", err)
		}
	case "help":
		fmt.Println("add - Добавить событие \n" +
			"list - Показать все события \n" +
			"remove - Удалить событие \n" +
			"update - Редактировать событие \n" +
			"reminder - Добавить напоминание \n" +
			"help - Показать справку \n" +
			"exit - Выход из программы")

	case "exit":
		err := c.calendar.Save()
		close(c.calendar.Notification)
		if err != nil {
			fmt.Println("Ошибка сохранения данных", err)
		}
		os.Exit(0)

	default:
		fmt.Println("Неизвестная команда:")
		fmt.Println("Введите 'help' для списка команд")
	}

}

func (c *Cmd) completer(d prompt.Document) []prompt.Suggest {
	suggestions := []prompt.Suggest{
		{Text: "add", Description: "Добавить событие"},
		{Text: "list", Description: "Показать все события"},
		{Text: "remove", Description: "Удалить событие"},
		{Text: "update", Description: "Редактировать событие"},
		{Text: "reminder", Description: "Добавить напоминание"},
		{Text: "help", Description: "Показать справку"},
		{Text: "exit", Description: "Выйти из программы"},
	}
	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}

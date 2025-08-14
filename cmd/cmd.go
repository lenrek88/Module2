package cmd

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/google/shlex"
	"github.com/lenrek88/app/calendar"
	"os"
	"strings"
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
		c.calendar.ShowEvents()

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

	case "help":
		fmt.Println("add - Добавить событие \n list - Показать все события \n remove - Удалить событие \n update - Редактировать событие \n help - Показать справку \n exit - Выход из программы")

	case "exit":
		err := c.calendar.Save()
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
		{Text: "help", Description: "Показать справку"},
		{Text: "exit", Description: "Выйти из программы"},
	}
	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}

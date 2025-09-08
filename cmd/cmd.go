package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/google/shlex"
	"github.com/lenrek88/app/calendar"
	"github.com/lenrek88/app/logger"
	"github.com/lenrek88/app/reminder"
)

type Cmd struct {
	calendar *calendar.Calendar
	logger   *Logger
}

func NewCmd(c *calendar.Calendar) *Cmd {
	return &Cmd{
		calendar: c,
		logger:   &Logger{},
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
			c.PrintWithLog("Reminder!" + msg)
		}
	}()
	p.Run()

}

func (c *Cmd) executor(input string) {
	parts, err := shlex.Split(input)
	if err != nil {
		c.PrintWithLog("error executing command")
		return
	}
	cmd := strings.ToLower(parts[0])

	switch cmd {
	case "add":
		c.logger.AddLogMessage("add command called")

		if len(parts) < 4 {
			c.PrintWithLog("Формат: add \"название события\" \"дата и время\" \"приоритет\"")
			return
		}

		title := parts[1]
		date := parts[2]
		priority := parts[3]

		e, err := c.calendar.AddEvent(title, date, priority)

		if err != nil {
			c.PrintWithLog("failed to add calendar events")
			logger.Error("failed to add calendar events")
		} else {
			c.PrintWithLog("Событие:" + e.Title + "добавлено")
			logger.Info("\"Событие:\", e.Title, \"добавлено\"")
		}

	case "list":
		c.logger.AddLogMessage("list command called")

		events := c.calendar.GetEvents()
		for _, e := range events {
			fmt.Println(e.ID, e.Title, "-", e.StartAt.Format("2006-01-02 15:04"), e.Priority)
		}

	case "remove":
		c.logger.AddLogMessage("remove command called")
		if len(parts) < 2 {
			c.PrintWithLog("Формат: remove \"ключ\"")
			return
		}
		c.calendar.DeleteEvent(parts[1])

	case "update":
		c.logger.AddLogMessage("update command called")
		if len(parts) < 5 {
			c.PrintWithLog("Формат: update \"ключ события\" \"название события\" \"дата и время\" \"приоритет\"")
			return
		}
		err := c.calendar.EditEvent(parts[1], parts[2], parts[3], parts[4])
		if err != nil {
			c.PrintWithLog("Update Error:" + err.Error())
			return
		}

	case "reminder":
		c.logger.AddLogMessage("reminder command called")

		if len(parts) < 4 {
			c.PrintWithLog("Формат: reminder \"ключ события\" \"сообщение напоминания\" \"дата и время\"")
			return
		}
		err := c.calendar.SetEventReminder(parts[1], parts[2], parts[3])
		if errors.Is(err, reminder.ErrEmptyMessage) {
			c.PrintWithLog("Can't set reminder with empty message")
		} else {
			c.PrintWithLog("Ошибка создания напоминания:" + err.Error())
		}
	case "help":
		c.logger.AddLogMessage("help command called")
		fmt.Println("add - Добавить событие \n" +
			"list - Показать все события \n" +
			"remove - Удалить событие \n" +
			"update - Редактировать событие \n" +
			"reminder - Добавить напоминание \n" +
			"help - Показать справку \n" +
			"logger - Показать весь лог \n" +
			"history - Показать историю ввода-вывода \n " +
			"exit - Выход из программы")
	case "logger":
		c.logger.AddLogMessage("logger command called")
		file, err := os.Open("app.log")
		if err != nil {
			log.Println("Ошибка открытия файла:", err)
			return
		}
		defer file.Close()

		_, err = io.Copy(os.Stdout, file)
		if err != nil {
			log.Println("Ошибка чтения файла:", err)
		}
	case "history":
		c.logger.AddLogMessage("ioHistory command called")
		c.ShowLog()
	case "exit":
		c.logger.AddLogMessage("exit command called")
		err := c.calendar.Save()
		close(c.calendar.Notification)
		if err != nil {
			log.Println("Ошибка сохранения данных", err)
		}
		os.Exit(0)

	default:
		log.Println("Неизвестная команда:")
		log.Println("Введите 'help' для списка команд")
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
		{Text: "logger", Description: "Показать весь лог"},
		{Text: "history", Description: "Показать историю ввода/вывода"},
		{Text: "exit", Description: "Выйти из программы"},
	}
	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}
